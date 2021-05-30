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

// ToArticles maps this Rss to a slice of Article
func (r Rss) ToArticles() []Article {
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
	// Walk through the HTML nodes
	node, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "string": content}).Warning("Cannot parse content")
		return
	}
	return walkNode(node)
}

func walkNode(node *html.Node) (thumbnail string, excerpt string) {
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
			// Gotten both thumbnail and excerpt. No need to continue the loop
			return
		}
		walkNode(child)
	}
	return
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
