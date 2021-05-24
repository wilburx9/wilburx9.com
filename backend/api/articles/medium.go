package articles

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	mediumKey = "Medium"
)

// Medium encapsulates the fetching and caching of medium articles
type Medium struct {
	Name string // should be Medium username (e.g "@Wilburx9") or publication (e.g flutter-community)
	internal.Fetcher
}

// FetchAndCache fetches and caches all Medium Articles
func (m Medium) FetchAndCache() {
	articles := m.FetchArticles()
	buf, _ := json.Marshal(articles)
	m.CacheData(getCacheKey(mediumKey), buf)
}

// GetCached returns cached Medium articles
func (m Medium) GetCached() ([]byte, error) {
	return m.GetCachedData(getCacheKey(mediumKey))
}

// FetchArticles fetches articles via HTTP
func (m Medium) FetchArticles() []Article {
	url := fmt.Sprintf("https://medium.com/feed/%s", m.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't init http request")
		return nil
	}

	res, err := m.HttpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't send request")
		return nil
	}
	defer res.Body.Close()

	var rss rss
	err = xml.NewDecoder(res.Body).Decode(&rss)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warning("Couldn't Unmarshall data")
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
		thumbnail, excerpt := getMediumThumbAndExcerpt(e.Encoded)
		articles[i] = Article{
			Title:     e.Title,
			Url:       e.Link,
			Thumbnail: thumbnail,
			PostedAt:  internal.StringToTime(time.RFC1123, e.PubDate),
			UpdatedAt: internal.StringToTime(time.RFC3339, e.Updated),
			Excerpt:   excerpt,
		}
	}
	return articles
}

func getMediumThumbAndExcerpt(content string) (thumbnail string, excerpt string) {

	// Get text with p tags
	var collectText func(*html.Node, *bytes.Buffer)
	collectText = func(n *html.Node, buf *bytes.Buffer) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			collectText(c, buf)
		}
	}

	// Get img src
	collectImgSrc := func(n *html.Node) string {
		for _, a := range n.Attr {
			if a.Key == "src" {
				return a.Val
			}
		}
		return ""
	}

	// Craw through the HTML nodes
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if excerpt == "" && node.Type == html.ElementNode && node.DataAtom == atom.P {
			buffer := &bytes.Buffer{}
			collectText(node.FirstChild, buffer)
			cleaned := strings.TrimSpace(internal.GetFirstNCodePoints(buffer.String(), 200))
			if utf8.RuneCountInString(cleaned) >= 80 {
				excerpt = fmt.Sprintf("%v.", strings.Split(cleaned, ".")[0])
			}
			return
		}
		if thumbnail == "" && node.Type == html.ElementNode && node.DataAtom == atom.Img {
			thumbnail = collectImgSrc(node)
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if thumbnail != "" && excerpt != "" {
				// Gotten both thumbnail and except. No need to continue the loop
				return
			}
			crawler(child)
		}
	}

	node, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "string": content}).Warning("Cannot parse content")
		return
	}
	crawler(node)
	return
}
