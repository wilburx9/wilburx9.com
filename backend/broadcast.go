package main

import (
	. "backend/common"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/mailerlite/mailerlite-go"
	"github.com/samber/lo"
	"html"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//go:embed newsletter.html
var newsletterFile string

func main() {
	lambda.Start(handleBroadcast)
}

// handleBroadcast creates a campaign and schedules to be sent.
func handleBroadcast(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	origin := req.Headers["origin"]

	if !InitSuccess() {
		return GenerateResponse(origin, http.StatusInternalServerError, "Something went wrong"), nil
	}

	status, msg := processBroadcastRequest(ctx, req.Body)
	return GenerateResponse(origin, status, msg), nil
}

// processBroadcastRequest schedules an email campaign about an hour form now.
// For Ghost not to retry the request because of short timeout, update webhook timeout. See https://forum.ghost.org/t/webhook-getting-triggered-multiple-times/16503/3
func processBroadcastRequest(ctx context.Context, body string) (int, string) {
	var reqData lambdaReqBody
	err := json.Unmarshal([]byte(body), &reqData)
	if err != nil {
		return http.StatusBadRequest, "invalid request body"
	}

	// Validate the request body
	validate := validator.New()
	err = validate.Struct(reqData)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	// 1. To prevent sending notification for drafts, confirm this post is published.
	// 2. To prevent sending emails for unpublished and republished posts, ensure
	//	  that the diff b/w publication and updated dates is not more than 30 minutes.
	if !reqData.canBroadcast() {
		return http.StatusBadRequest, fmt.Sprintf(
			"this post is too old to be rescheduled. It was created at %v and updated at %v ",
			reqData.Post.Current.PublishedAt,
			reqData.Post.Current.UpdatedAt,
		)
	}

	post := reqData.toPost()

	content, err := parseEmailTemplate(post)
	if err != nil {
		log.Println("parse template error: ", err)
		return http.StatusInternalServerError, "something went wrong"
	}

	campaignId, err := createCampaign(ctx, post, content)
	if err != nil {
		log.Println("campaign creation error: ", err)
		return http.StatusBadGateway, "something went wrong while creating campaign"
	}

	err = scheduleCampaign(ctx, campaignId)
	if err != nil {
		log.Println("campaign scheduling error: ", err)
		return http.StatusBadGateway, "something went wrong while scheduling campaign"
	}

	return http.StatusOK, "Successfully scheduled"
}

func scheduleCampaign(ctx context.Context, campaignId string) error {
	timezone := AppConfig.TimeZone
	when := time.Now().Add(time.Hour)

	list, _, err := MailClient.Timezone.List(ctx)
	if err != nil {
		return err
	}

	var timezoneId int
	for _, tz := range list.Data {
		if tz.Name == timezone {
			timezoneId, err = strconv.Atoi(tz.Id)
			if err != nil {
				return err
			}
			break
		}
	}
	if timezoneId == 0 {
		return fmt.Errorf("no valid timezone id found for %v. OS timezone: %v", timezone, os.Getenv("TZ"))
	}

	schedule := &mailerlite.ScheduleCampaign{
		Delivery: mailerlite.CampaignScheduleTypeScheduled,
		Schedule: &mailerlite.Schedule{
			Date:       when.Format("2006-01-02"),
			Hours:      fmt.Sprintf("%2d", when.Hour()),
			Minutes:    fmt.Sprintf("%2d", when.Minute()),
			TimezoneID: timezoneId,
		},
	}

	_, _, err = MailClient.Campaign.Schedule(ctx, campaignId, schedule)
	return err
}

func createCampaign(ctx context.Context, post Post, content string) (string, error) {
	allSegments, _, err := MailClient.Segment.List(ctx, &mailerlite.ListSegmentOptions{})
	if err != nil {
		return "", err
	}

	primarySegment := fmt.Sprintf("%v: %v", Blog, post.PrimaryTag)
	segment, ok := lo.Find(allSegments.Data, func(seg mailerlite.Segment) bool {
		return strings.EqualFold(seg.Name, primarySegment)
	})

	if !ok {
		return "", errors.New("won't send campaigns for non-(software or photography) articles")
	}

	sender := AppConfig.EmailSender
	emails := &[]mailerlite.Emails{
		{
			Subject:  post.Title,
			FromName: post.Author,
			From:     sender,
			Content:  content,
		},
	}
	campaign := &mailerlite.CreateCampaign{
		Name:     post.Title,
		Type:     mailerlite.CampaignTypeRegular,
		Emails:   *emails,
		Segments: []string{segment.ID},
	}
	c, _, err := MailClient.Campaign.Create(ctx, campaign)
	if err != nil {
		return "", err
	}
	return c.Data.ID, nil
}

func parseEmailTemplate(post Post) (string, error) {
	var emailContent bytes.Buffer
	t, err := template.New("newsletter").Parse(newsletterFile)
	if err != nil {
		return "", err
	}

	err = t.Execute(&emailContent, post)
	if err != nil {
		return "", err
	}

	return emailContent.String(), nil
}

func (l lambdaReqBody) canBroadcast() bool {
	// 1. To prevent sending notification for drafts, confirm this post is published.
	// 2. To prevent sending emails for a posts that are unpublished and  then republished, ensure
	//	  that the diff b/w publication and updated dates is not more than 30 minutes.

	postData := l.Post.Current
	return postData.Status == "published" && math.Abs(postData.PublishedAt.Sub(postData.UpdatedAt).Minutes()) <= 30
}

func (l lambdaReqBody) toPost() Post {
	p := l.Post.Current

	featureImage := p.FeatureImage

	// Check if this post is a reference to an external article and retrieve the feature image.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(p.HTML))
	if err == nil {

		if slices.ContainsFunc(p.Tags, func(item slug) bool {
			return strings.EqualFold(item.Name, "#external")
		}) {
			if featureImage == "" {
				img := doc.Find("div.kg-bookmark-thumbnail img")
				if img.Length() > 0 {
					featureImage, _ = img.Attr("src")
				}
			}
		}
	}

	// Remove html tags from the caption
	cleanedCaption := regexp.MustCompile("<[^>]+>").ReplaceAllString(p.FeatureImageCaption, "")
	featureImageCaption := html.UnescapeString(cleanedCaption)

	return Post{
		Author:              p.PrimaryAuthor.Name,
		Title:               p.Title,
		PubDate:             p.PublishedAt.Format("02 Jan 2006"),
		FeatureImage:        featureImage,
		FeatureImageCaption: featureImageCaption,
		Excerpt:             p.Excerpt,
		URL:                 p.URL,
		PrimaryTag:          p.PrimaryTag,
		Tags:                p.Tags,
	}
}

type Post struct {
	Author              string
	Title               string
	PubDate             string
	FeatureImage        string
	FeatureImageCaption string
	Excerpt             string
	URL                 string
	PrimaryTag          slug
	Tags                []slug
}

type lambdaReqBody struct {
	Post struct {
		Current struct {
			Excerpt             string    `json:"excerpt" validate:"required"`
			FeatureImage        string    `json:"feature_image" validate:"http_url"`
			FeatureImageCaption string    `json:"feature_image_caption" validate:"required"`
			ID                  string    `json:"id" validate:"required"`
			PublishedAt         time.Time `json:"published_at" validate:"required"`
			ReadingTime         int64     `json:"reading_time" validate:"required"`
			Status              string    `json:"status" validate:"required"`
			Title               string    `json:"title" validate:"required"`
			UpdatedAt           time.Time `json:"updated_at" validate:"required"`
			URL                 string    `json:"url" validate:"http_url"`
			Visibility          string    `json:"visibility" validate:"required"`
			HTML                string    `json:"html" validate:"required"`
			PrimaryTag          slug      `json:"primary_tag" validate:"required"`
			Tags                []slug    `json:"tags" validate:"required"`
			PrimaryAuthor       struct {
				Name string `json:"name" validate:"required"`
			} `json:"primary_author" validate:"required"`
		} `json:"current" validate:"required"`
	} `json:"post" validate:"required"`
}

type slug struct {
	Slug string `json:"slug" validate:"required"`
	Name string `json:"name" validate:"required"`
}
