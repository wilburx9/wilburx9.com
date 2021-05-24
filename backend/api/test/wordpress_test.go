package test

import (
	"github.com/wilburt/wilburx9.dev/backend/api/articles"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"testing"
)

func TestWordPressFetchArticles(t *testing.T) {
	var w = articles.Wordpress{URL: "https://example.com/wp-json/wp/v2/posts", Fetch: internal.Fetch{
		Db:         nil,
		HttpClient: &internal.HttpClientMock{ResponseFilePath: "../testdata/wordpress_response.json"},
	}}
	var results = w.FetchArticles()
	if len(results) != 2 {
		t.Error()
	}
	if results[0].Title != "Lorem ipsum is placeholder text commonly used" {
		t.Error()
	}
}
