package articles

import (
	"github.com/stretchr/testify/assert"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"testing"
)

func TestMediumFetchArticles(t *testing.T) {
	var m = Medium{Name: "testUser", Fetch: internal.Fetch{
		Db:         nil,
		HttpClient: &internal.HttpClientMock{ResponseFilePath: "./testdata/medium_response.xml"},
	}}
	var result = m.fetchArticles().Articles

	first := result[0]
	second := result[1]

	assert.Equal(t, first.Title, "Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	assert.Equal(t, first.Thumbnail, "https://cdn-images-1.medium.com/max/960/1*bbkcrsggiQLxNDRAgHiSBQ.png")
	assert.Equal(t, first.Url, "https://medium.com/lorem/lorem---q")
	assert.NotEmpty(t, first.Excerpt)
	assert.Empty(t, second.Thumbnail)
}
