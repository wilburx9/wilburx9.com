package articles

import (
	"github.com/wilburt/wilburx9.dev/backend/common"
	"testing"
)

func TestMediumFetchArticles(t *testing.T) {
	var m = medium{name: "testUser"}
	clientMock := common.HttpClientMock{ResponseFilePath: "./testdata/medium_response.xml"}
	var articles = m.fetchArticles(&clientMock)

	first := articles[0]
	second := articles[1]

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
