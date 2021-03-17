package articles

import (
	"github.com/wilburt/wilburx9.dev/backend/common"
	"net/http"
	"os"
	"testing"
)

func TestFetchArticles(t *testing.T) {
	common.SetUpLogger(false)
	var m = medium{name: "testUser"}
	var articles = m.fetchArticles(&HttpClientMock{})

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

type HttpClientMock struct{}

func (cm *HttpClientMock) Do(_ *http.Request) (*http.Response, error) {
	file, err := os.Open("./testdata/medium_response.xml")
	if err != nil {
		return nil, err
	}

	return &http.Response{Body: file}, nil
}
