package gallery

import (
	"github.com/stretchr/testify/assert"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestUnsplashFetchImages(t *testing.T) {
	const expectedResults = 3
	var header = http.Header{}
	header.Add("X-Total", strconv.Itoa(expectedResults))

	var u = Unsplash{Username: "x", AccessKey: "xa", Fetch: internal.Fetch{
		HttpClient: &internal.HttpClientMock{ResponseFilePath: "./testdata/unsplash_response.json", Header: header},
	}}
	var images = u.FetchImages()

	assert.Equal(t, len(images), expectedResults)

	first := images[0].(models.Image)
	assert.Equal(t, first.Url, "https://images.unsplash.com/photo-56789-098yhj?crop=entropy&cs=srgb&fm=jpg&ixid=OIFGHJIUGGH=rb-1.2.1&q=85")
	assert.Equal(t, first.Page, "https://unsplash.com/photos/blah_blah")
	assert.Equal(t, first.Caption, "ABC")
	assert.NotEqual(t, first.UploadedOn.Year(), time.Now().Year())

	user, ok := first.Meta["user"].(models.User)
	if assert.True(t, ok) {
		assert.Equal(t, user.Username, "aafgotiigg")
		assert.Equal(t, user.Name, "Larry Emeka")
	}
}
