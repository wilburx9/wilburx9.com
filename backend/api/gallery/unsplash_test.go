package gallery

import (
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
	var images = u.FetchImage([]models.Image{}, 0)

	if len(images) != expectedResults {
		t.Errorf("Recursive fetching of images failed. Expected 2 but got %d", len(images))
	}

	first := images[0]
	if first.Url != "https://unsplash.com/photos/blah_blah" {
		t.Error("Failed to parse image url")
	}

	if first.Caption != "ABC" {
		t.Error("Failed to parse image caption")
	}

	if first.UploadedAt.Year() == time.Now().Year() {
		t.Error("Failed to parse image creation date")
	}

	if first.Src != "https://images.unsplash.com/photo-56789-098yhj?crop=entropy&cs=srgb&fm=jpg&ixid=OIFGHJIUGGH=rb-1.2.1&q=85" {
		t.Error("Failed to parse image src")
		t.Errorf("Image %q", first.Src)
	}

	user, ok := first.Meta["user"].(models.User)
	if !ok {
		t.Errorf("Failed to parse image user. Got %T but wanted user", first.Meta["user"])
	}

	if user.Username != "aafgotiigg" || user.Name != "Larry Emeka" {
		t.Error("Failed to parse image user")
	}
}
