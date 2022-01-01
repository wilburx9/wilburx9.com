package articles

import (
	"github.com/stretchr/testify/assert"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/update"
	"testing"
)

func TestWordPressFetchArticles(t *testing.T) {
	var w = WordPress{URL: "https://example.com/wp-json/wp/v2/posts", BaseCache: update.BaseCache{
		Db:         nil,
		HttpClient: &internal.HttpClientMock{ResponseFilePath: "./testdata/wordpress_response.json"},
	}}
	var results = w.fetchArticles().Articles

	assert.Equal(t, len(results), 2)
	assert.Equal(t, results[0].Title, "Lorem ipsum is placeholder text commonly used")
}
