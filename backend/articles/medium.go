package articles

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"regexp"
	"time"
)

const (
	mediumKey = "Medium"
)

// Medium encapsulates the fetching and caching of medium articles
type Medium struct {
	Name string // should be Medium username (e.g "@Wilburx9") or publication (e.g flutter-community)
	common.Fetcher
}

// FetchAndCache fetches and caches all Medium Articles
func (m Medium) FetchAndCache() {
	articles := m.fetchArticles()
	buf, _ := json.Marshal(articles)
	m.CacheData(getCacheKey(mediumKey), buf)
}

// GetCached returns cached Medium articles
func (m Medium) GetCached() ([]byte, error) {
	return m.GetCachedData(getCacheKey(mediumKey))
}

func (m Medium) fetchArticles() []Article {
	url := fmt.Sprintf("https://Medium.com/feed/%s", m.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		common.LogError(err)
		return nil
	}

	res, err := m.HttpClient.Do(req)
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
	// Yes bobince, I am parsing HTML with regex. Bite me! (https://stackoverflow.com/a/1732454/6181476)
	var imgReg = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	subMatch := imgReg.FindStringSubmatch(body)
	if subMatch == nil {
		return ""
	}
	return subMatch[1]
}
