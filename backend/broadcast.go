package main

import (
	"backend/common"
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
	"html"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	iso8601 = "2006-01-02T15:04:05-0700"
)

//go:embed newsletter.html
var newsletterFile string

func main() {
	lambda.Start(handleBroadcast)
}

// handleBroadcast creates a campaign and schedules to be sent.
func handleBroadcast(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx = context.WithValue(ctx, common.ConfigKey, common.NewConfig())

	status, msg := processRequest(ctx, req.Body)

	return common.MakeResponse(status, msg), nil
}

// processRequest schedules an email campaign about an hour form now.
// For Ghost not to retry the request because of short timeout, update webhook timeout. See https://forum.ghost.org/t/webhook-getting-triggered-multiple-times/16503/3
func processRequest(ctx context.Context, body string) (int, string) {
	var reqData lambdaReqBody
	err := json.Unmarshal([]byte(body), &reqData)
	if err != nil {
		return http.StatusBadRequest, fmt.Sprintf("request body parse error: %v", err)
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

	postData := reqData.Post.Current
	if postData.Status != "published" || math.Abs(postData.PublishedAt.Sub(postData.UpdatedAt).Minutes()) > 30 {
		return http.StatusBadRequest, fmt.Sprintf(
			"this post is too old to be rescheduled. It was created at %v and updated at %v ",
			postData.PublishedAt,
			postData.UpdatedAt,
		)
	}

	post := reqData.toPost()

	content, err := parseEmailTemplate(post)
	if err != nil {
		log.Println("parse template error: ", err)
		return http.StatusInternalServerError, "something went wrong"
	}

	campaignId, err := createCampaign(ctx, post)
	if err != nil {
		return http.StatusBadGateway, err.Error()
	}

	err = setCampaignContent(ctx, campaignId, content)
	if err != nil {
		return http.StatusBadGateway, err.Error()
	}

	err = scheduleCampaign(ctx, campaignId)
	if err != nil {
		return http.StatusBadGateway, err.Error()
	}

	return http.StatusOK, "Successfully scheduled"
}

func scheduleCampaign(ctx context.Context, campaignId string) error {
	// Use a time an hour into the future, so I can manually cancel the campaign if I don't want it to be sent.
	f := time.Now().Add(time.Hour)
	// Round it to quarterly hour. See https://mailchimp.com/developer/marketing/api/campaigns/schedule-campaign
	minutes := f.Minute()
	remainder := minutes % 15
	var when time.Time
	if (remainder * 2) < 15 { // If the remainder is less than 7.5 (halfway between 0 and 15), round down; otherwise, round up
		when = time.Date(f.Year(), f.Month(), f.Day(), f.Hour(), minutes-remainder, 0, 0, f.Location())
	} else {
		when = time.Date(f.Year(), f.Month(), f.Day(), f.Hour(), minutes+(15-remainder), 0, 0, f.Location())
	}

	reqBody := map[string]any{"schedule_time": when.Format(iso8601)}

	success := common.MakeMailChimpRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("campaigns/%v/actions/schedule", campaignId),
		reqBody,
		nil,
	)
	if !success {
		return errors.New("failed to schedule newsletter")
	}
	return nil
}

func setCampaignContent(ctx context.Context, campaignId string, content string) error {
	reqBody := map[string]any{"html": content}

	success := common.MakeMailChimpRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("campaigns/%v/content", campaignId),
		reqBody,
		nil,
	)
	if !success {
		return errors.New("failed to set newsletter content")
	}
	return nil
}

func createCampaign(ctx context.Context, post Post) (string, error) {
	reqBody, err := post.toRequestBody(common.ConfigFromContext(&ctx))
	if err != nil {
		return "", errors.New("failed to create campaign")
	}

	var campaign struct {
		ID string `json:"id"`
	}

	success := common.MakeMailChimpRequest(
		ctx,
		http.MethodPost,
		"campaigns",
		reqBody,
		&campaign,
	)
	if !success {
		return "", errors.New("failed to create campaign")
	}
	return campaign.ID, nil
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

func (p Post) toRequestBody(config *common.Config) (map[string]any, error) {
	var segment int
	switch p.Tag {
	case common.Programming:
		segment, _ = strconv.Atoi(config.ProgrammingSegment)
	case common.Photography:
		segment, _ = strconv.Atoi(config.PhotographySegment)
	}

	if segment == 0 {
		return nil, errors.New("won't send campaigns for non-(programming or photography) articles")
	}

	return map[string]any{
		"type": "regular",
		"settings": map[string]string{
			"title":        fmt.Sprintf("New Publication: %v", p.Title),
			"subject_line": p.Title,
			"preview_text": p.Excerpt,
			"from_name":    p.Author,
			"reply_to":     config.EmailSender,
		},
		"recipients": map[string]any{
			"list_id": config.NewsletterListId,
			"segment_opts": map[string]int{
				"saved_segment_id": segment,
			},
		},
	}, nil
}

func (e lambdaReqBody) toPost() Post {
	p := e.Post.Current

	featureImage := p.FeatureImage

	// Check if this post is a reference to an external article and retrieve the feature image.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(p.HTML))
	if err == nil {
		bookmark := doc.Find("figure.kg-bookmark-card")

		// A post which is just a reference to an external article
		// will contain nothing but the bookmark card and the reading time.
		if bookmark.Length() > 0 && bookmark.Children().Length() != 2 {

			// Only set the feature image if this post didn't have one
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
		Tag:                 p.PrimaryTag.Slug,
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
	Tag                 string
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

			PrimaryAuthor struct {
				Name string `json:"name" validate:"required"`
			} `json:"primary_author" validate:"required"`

			PrimaryTag struct {
				Name string `json:"name" validate:"required"`
				Slug string `json:"slug" validate:"required"`
			} `json:"primary_tag" validate:"required"`
		} `json:"current" validate:"required"`
	} `json:"post" validate:"required"`
}
