package articles

import (
	"encoding/xml"
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"regexp"
	"time"
)

type medium struct {
	name string
}

func (m medium) fetchArticles(client common.HttpClient) []Article {
	url := fmt.Sprintf("https://medium.com/feed/%s", m.name)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		common.Logger.Errorf("An error while creating http request for %s :: \"%v\"", url, err)
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		common.Logger.Errorf("An error occurred while sending request for %s :: \"%v\"", url, err)
		return nil
	}

	defer res.Body.Close()

	var rss rss
	err = xml.NewDecoder(res.Body).Decode(&rss)
	if err != nil {
		common.Logger.Errorf("An error occurred while decoging response for %s :: \"%v\"", url, err)
		return nil
	}
	return rss.toArticles()
}

type rss struct {
	Channel struct {
		Item []struct {
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			PubDate string `xml:"pubDate"`
			Updated string `xml:"updated"`
			Encoded string `xml:"encoded"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (r rss) toArticles() []Article {
	var articles = make([]Article, len(r.Channel.Item))
	for i, s := range r.Channel.Item {
		articles[i] = Article{
			Title:     s.Title,
			Url:       s.Link,
			Thumbnail: getThumbnail(s.Encoded),
			PostedAt:  stringToTime(time.RFC1123, s.PubDate),
			UpdatedAt: stringToTime(time.RFC3339, s.Updated),
		}
	}
	return articles
}

func getThumbnail(body string) string {
	var imgReg = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	subMatch := imgReg.FindStringSubmatch(body)
	if subMatch == nil {
		return ""
	}
	return subMatch[1]
}

func stringToTime(layout string, timeStr string) time.Time {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		common.Logger.Infof("Could not parse time: \"%s\" with layout:\"%s\" :: \"%v\"", timeStr, layout, err)
		return time.Now()
	}
	return t
}
