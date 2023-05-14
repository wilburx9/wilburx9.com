package main

import (
	"backend/common"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
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
	http.HandleFunc("/broadcast", handleBroadcast)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBroadcast(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling webhook ============ ")
	// Print the request headers
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Close the request body
	defer r.Body.Close()

	// Parse the request body as JSON
	var requestBody interface{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the request body to JSON
	requestBodyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	campaignId, err := handlePublishEvent(string(requestBodyJSON))
	if err != nil {
		fmt.Println(fmt.Sprintf("%v :: %v", campaignId, err))
		// Respond with an error if the campaign has not been created. This is so Ghost can retry
		if campaignId == "" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Request received successfully!")
}

// handlePublishEvent schedules an email campaign about an hour form now.
// For this to work, increase Ghost's webhook timeout. See https://forum.ghost.org/t/webhook-getting-triggered-multiple-times/16503/3
func handlePublishEvent(body string) (string, error) {
	fmt.Println("Handling publication event")
	var event publishEvent
	err := json.Unmarshal([]byte(body), &event)
	if err != nil {
		return "", fmt.Errorf("request body parse error: %w", err)
	}

	// Confirm that event.Status == "published" before publishing
	// Also check that the difference publication and updated dates is not more than 5 mins,
	//this is so we don't send emails for posts that was unpublished and republished
	post, err := event.toPost()
	if err != nil {
		return "", fmt.Errorf("request body to post mapping error: %w", err)
	}

	fmt.Println(fmt.Sprintf("Post data %+v\n", post))

	content, err := parseEmailTemplate(post, "newsletter.html")
	if err != nil {
		return "", fmt.Errorf("parse template error: %w", err)
	}

	campaignId, err := createCampaign(post)
	if err != nil {
		return campaignId, fmt.Errorf("campaign creation error: %w", err)
	}

	err = setCampaignContent(campaignId, content)
	if err != nil {
		return campaignId, fmt.Errorf("campaign update error: %w", err)
	}

	err = scheduleCampaign(campaignId)
	if err != nil {
		return campaignId, fmt.Errorf("campaign schedule error: %w", err)
	}

	return campaignId, nil
}

func scheduleCampaign(campaignId string) error {
	fmt.Println("Scheduling", "::", campaignId)
	dc := os.Getenv("MAILCHIMP_DC")
	token := os.Getenv("MAILCHIMP_TOKEN")
	u := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/campaigns/%v/actions/schedule", dc, campaignId)

	when := time.Now().Add(time.Hour)
	reqData := map[string]interface{}{"schedule_time": when.Format(iso8601)}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf(string(body))
}

func setCampaignContent(campaignId string, content string) error {
	fmt.Println(fmt.Sprintf("Setting content :: %v", campaignId))
	dc := os.Getenv("MAILCHIMP_DC")
	token := os.Getenv("MAILCHIMP_TOKEN")
	u := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/campaigns/%v/content", dc, campaignId)

	reqData := map[string]interface{}{"html": content}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, u, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf(string(body))
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

func createCampaign(post Post) (string, error) {
	fmt.Println("Creating campaign ")
	dc := os.Getenv("MAILCHIMP_DC")
	token := os.Getenv("MAILCHIMP_TOKEN")
	u := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/campaigns", dc)

	reqBody, err := post.toRequestBody()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf(string(body))
	}

	var campaign struct {
		ID string `json:"id"`
	}

	err = json.NewDecoder(res.Body).Decode(&campaign)
	if err != nil {
		return "", err
	}

	return campaign.ID, nil
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

func (e publishEvent) toPost() (Post, error) {
	p := e.Post.Current

	pubDate, err := time.Parse(time.RFC3339, p.PublishedAt)
	if err != nil {
		return Post{}, err
	}

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
		PubDate:             pubDate.Format("02 Jan 2006"),
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

type publishEvent struct {
	Post struct {
		Current struct {
			Excerpt             string `json:"excerpt"`
			FeatureImage        string `json:"feature_image"`
			FeatureImageCaption string `json:"feature_image_caption"`
			ID                  string `json:"id"`
			PublishedAt         string `json:"published_at"`
			ReadingTime         int64  `json:"reading_time"`
			Status              string `json:"status"`
			Title               string `json:"title"`
			UpdatedAt           string `json:"updated_at"`
			URL                 string `json:"url"`
			Visibility          string `json:"visibility"`
			HTML                string `json:"html"`

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
