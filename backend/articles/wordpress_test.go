package articles

import (
	"github.com/wilburt/wilburx9.dev/backend/common"
	"testing"
)

func TestWordPressFetchArticles(t *testing.T) {
	var w = Wordpress{"https://example.com/wp-json/wp/v2/posts"}
	clientMock := common.HttpClientMock{ResponseFilePath: "./testdata/wordpress_response.json"}
	var articles = w.fetchArticles(&clientMock)
	if len(articles) != 2 {
		t.Error()
	}
	if articles[0].Title != "Lorem ipsum is placeholder text commonly used" {
		t.Error()
	}
}
