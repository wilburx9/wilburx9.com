package models

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	minAcceptableExcerpt = 80
	minExcerpt           = 200
)

// Rss is a container for Medium Rss feed data
type Rss struct {
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

// ToResult creates ArticleResult by mapping Rss to a slice of Article
func (r Rss) ToResult() ArticleResult {
	var articles = make([]Article, len(r.Channel.Item))
	for i, e := range r.Channel.Item {
		thumbnail, excerpt := getMediumThumbAndExcerpt(e.Encoded)
		articles[i] = Article{
			Title:     e.Title,
			Url:       e.Link,
			Thumbnail: thumbnail,
			PostedAt:  internal.StringToTime(time.RFC1123, e.PubDate),
			UpdatedAt: internal.StringToTime(time.RFC3339, e.Updated),
			Excerpt:   fmt.Sprintf("%v..", excerpt),
		}
	}
	return ArticleResult{
		Result:   internal.Result{UpdatedAt: time.Now()},
		Articles: articles,
	}
}

func getMediumThumbAndExcerpt(content string) (thumbnail string, excerpt string) {
	// Walk through the HTML nodes
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		if excerpt == "" && node.Type == html.ElementNode && node.DataAtom == atom.P {
			buffer := &bytes.Buffer{}
			collectText(node.FirstChild, buffer)
			excerpt = getCleanedText(buffer.String())
			return
		}
		if thumbnail == "" && node.Type == html.ElementNode && node.DataAtom == atom.Img {
			thumbnail = collectImgSrc(node)
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if thumbnail != "" && excerpt != "" {
				// Gotten both thumbnail and excerpt. No need to continue the loop
				return
			}
			walker(child)
		}
	}

	node, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "string": content}).Warning("Cannot parse content")
		return
	}
	walker(node)
	return
}

// Concatenate the sentences until the concatenated string is approx. minExcerpt characters
func getCleanedText(s string) string {
	s = strings.TrimSpace(s)
	buffer := &bytes.Buffer{}
	if utf8.RuneCountInString(s) >= minAcceptableExcerpt {
		splits := strings.Split(s, ".") // Separate it into sentences
		for i := range splits {
			str := splits[i]
			lenStr := len(str)
			if lenStr > 0 {
				// End it with a dot if it doesn't already end with a dot
				if "." != str[lenStr-1:] {
					str = fmt.Sprintf("%v.", str)
				}
				buffer.WriteString(str)
				if utf8.RuneCountInString(buffer.String()) >= minExcerpt {
					return buffer.String()
				}
			}
		}
	}
	return buffer.String()
}

// Get texts in node
func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

// Get img src
func collectImgSrc(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "src" {
			return a.Val
		}
	}
	return ""
}
