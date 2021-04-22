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
	name string // should be medium username (e.g "@Wilburx9") or publication
}

func (m medium) fetchArticles(client common.HttpClient) []Article {
	url := fmt.Sprintf("https://medium.com/feed/%s", m.name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		common.LogError(err)
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		common.LogError(err)
		return nil
	}
	defer res.Body.Close()

	var rss rss
	err = xml.NewDecoder(res.Body).Decode(&rss)
	if err != nil {
		common.LogError(err)
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
	for i, e := range r.Channel.Item {
		articles[i] = Article{
			Title:     e.Title,
			Url:       e.Link,
			Thumbnail: getThumbnail(e.Encoded),
			PostedAt:  common.StringToTime(time.RFC1123, e.PubDate),
			UpdatedAt: common.StringToTime(time.RFC3339, e.Updated),
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
