package articles

import (
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"testing"
)

func TestWordPressFetchArticles(t *testing.T) {
	var w = Wordpress{URL: "https://example.com/wp-json/wp/v2/posts", Fetch: internal.Fetch{
		Db:         nil,
		HttpClient: &internal.HttpClientMock{ResponseFilePath: "../testdata/wordpress_response.json"},
	}}
	var results = w.fetchArticles()
	if len(results) != 2 {
		t.Error()
	}
	if results[0].Title != "Lorem ipsum is placeholder text commonly used" {
		t.Error()
	}
}
