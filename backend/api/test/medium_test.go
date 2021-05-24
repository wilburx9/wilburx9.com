package test

import (
	"github.com/wilburt/wilburx9.dev/backend/api/articles"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"testing"
)

func TestMediumFetchArticles(t *testing.T) {
	var m = articles.Medium{Name: "testUser", Fetch: internal.Fetch{
		Db:         nil,
		HttpClient: &internal.HttpClientMock{ResponseFilePath: "./testdata/medium_response.xml"},
	}}
	var result = m.FetchArticles()

	first := result[0]
	second := result[1]

	if first.Title != "Lorem ipsum dolor sit amet, consectetur adipiscing elit" {
		t.Error()
	}
	if first.Thumbnail != "https://cdn-images-1.medium.com/max/960/1*bbkcrsggiQLxNDRAgHiSBQ.png" {
		t.Error()
	}

	if first.Url != "https://medium.com/lorem/lorem---q" {
		t.Error()
	}
	if second.Thumbnail != "" {
		t.Error()
	}
}
