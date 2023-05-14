package main

import (
	"backend/common"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/net/html"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	iso8601 = "2006-01-02T15:04:05-0700"
)

func main() {
	lambda.Start(handleBroadcast)
}

// handleBroadcast creates a campaign and schedules to be sent.
func handleBroadcast(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	campaignId, err := processRequest(ctx, req.Body)

	// The broadcast was successfully created and scheduled
	if err == nil {
		return common.MakeResponse(
			http.StatusOK,
			fmt.Sprintf("Scheduled broadcast %v", campaignId),
		), nil
	}

	// The broadcast was created but not scheduled successfully
	if campaignId != "" {
		return common.MakeResponse(
			http.StatusCreated,
			fmt.Sprintf("Scheduled broadcast %v", campaignId),
		), nil
	}

	// The campaign was neither created nor scheduled. Return non 2XX so Ghost will retry.
	return common.MakeResponse(
		http.StatusInternalServerError,
		err.Error(),
	), nil
}

// processRequest schedules an email campaign about an hour form now.
// For Ghost not to retry the request because of short timeout, update webhook timeout. See https://forum.ghost.org/t/webhook-getting-triggered-multiple-times/16503/3
func processRequest(ctx context.Context, body string) (string, error) {
	var reqData lambdaReqBody
	err := json.Unmarshal([]byte(body), &reqData)
	if err != nil {
		return "", fmt.Errorf("request body parse error: %w", err)
	}

	// 1. To prevent sending notification for drafts, confirm this post is published.
	// 2. To prevent sending emails for unpublished and republished posts, ensure
	//	  that the diff b/w publication and updated dates is not more than 30 minutes.

	postData := reqData.Post.Current
	if postData.Status != "published" || math.Abs(postData.PublishedAt.Sub(postData.UpdatedAt).Minutes()) > 30 {
		log.Println(
			"this post is too old to be rescheduled. It was created at ",
			postData.PublishedAt,
			" and updated at ",
			postData.UpdatedAt,
		)
		return "", nil
	}

	post, err := reqData.toPost()
	if err != nil {
		return "", fmt.Errorf("request body to post mapping error: %w", err)
	}

	content, err := parseEmailTemplate(post, "newsletter.html")
	if err != nil {
		return "", fmt.Errorf("parse template error: %w", err)
	}

	campaignId, err := createCampaign(ctx, post)
	if err != nil {
		return campaignId, fmt.Errorf("campaign creation error: %w", err)
	}

	err = setCampaignContent(ctx, campaignId, content)
	if err != nil {
		return campaignId, fmt.Errorf("campaign update error: %w", err)
	}

	err = scheduleCampaign(ctx, campaignId)
	if err != nil {
		return campaignId, fmt.Errorf("campaign schedule error: %w", err)
	}

	return campaignId, nil
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

	reqBody := map[string]interface{}{"schedule_time": when.Format(iso8601)}

	err := common.MakeMailChimpRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("campaigns/%v/actions/schedule", campaignId),
		reqBody,
		nil,
	)
	return err
}

func setCampaignContent(ctx context.Context, campaignId string, content string) error {
	reqBody := map[string]interface{}{"html": content}

	err := common.MakeMailChimpRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("campaigns/%v/content", campaignId),
		reqBody,
		nil,
	)
	return err
}

func createCampaign(ctx context.Context, post Post) (string, error) {
	reqBody, err := post.toRequestBody()
	if err != nil {
		return "", err
	}

	var campaign struct {
		ID string `json:"id"`
	}

	err = common.MakeMailChimpRequest(
		ctx,
		"campaigns",
		http.MethodPost,
		reqBody,
		&campaign,
	)
	if err != nil {
		return "", err
	}
	return campaign.ID, nil
}

func parseEmailTemplate(post Post, templateFile string) (string, error) {
	fileBytes, err := os.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	var emailContent bytes.Buffer
	t, err := template.New("newsletter").Parse(string(fileBytes))
	if err != nil {
		return "", err
	}

	err = t.Execute(&emailContent, post)
	if err != nil {
		return "", err
	}

	return emailContent.String(), nil
}

func (p Post) toRequestBody() ([]byte, error) {
	listId := os.Getenv("MAILCHIMP_LIST_ID")
	replyTo := os.Getenv("MAILCHIMP_REPLY_TO")
	var segment int
	switch p.Tag {
	case common.Programming:
		segment, _ = strconv.Atoi(os.Getenv("MAILCHIMP_PROGRAMMING_SEGMENT"))
	case common.Photography:
		segment, _ = strconv.Atoi(os.Getenv("MAILCHIMP_PHOTOGRAPHY_SEGMENT"))
	}

	if segment == 0 {
		return nil, errors.New("won't send campaigns for non-(programming or photography) articles")
	}

	data := map[string]interface{}{
		"type": "regular",
		"settings": map[string]string{
			"title":        fmt.Sprintf("New Publication: %v", p.Title),
			"subject_line": p.Title,
			"preview_text": p.Excerpt,
			"from_name":    p.Author,
			"reply_to":     replyTo,
		},
		"recipients": map[string]interface{}{
			"list_id": listId,
			"segment_opts": map[string]int{
				"saved_segment_id": segment,
			},
		},
	}

	return json.Marshal(data)
}

func (e lambdaReqBody) toPost() (Post, error) {
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
	}, nil
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
			Excerpt             string    `json:"excerpt"`
			FeatureImage        string    `json:"feature_image"`
			FeatureImageCaption string    `json:"feature_image_caption"`
			ID                  string    `json:"id"`
			PublishedAt         time.Time `json:"published_at"`
			ReadingTime         int64     `json:"reading_time"`
			Status              string    `json:"status"`
			Title               string    `json:"title"`
			UpdatedAt           time.Time `json:"updated_at"`
			URL                 string    `json:"url"`
			Visibility          string    `json:"visibility"`
			HTML                string    `json:"html"`

			PrimaryAuthor struct {
				Name string `json:"name"`
			} `json:"primary_author"`

			PrimaryTag struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"primary_tag"`
		} `json:"current"`
	} `json:"post"`
}
