package gallery

import (
	"errors"
	"github.com/sirupsen/logrus"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wilburt/wilburx9.dev/backend/api/gallery/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestUnsplashFetchImages(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	assert := testify.New(t)
	file, _ := os.Open("./testdata/unsplash_response.json")
	httpClient := new(internal.MockHttpClient)

	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	var u = Unsplash{HttpClient: httpClient}
	images, err := u.fetchImages()
	assert.Nil(err)
	assert.Equal(1, len(images))
	first := images[0].(models.Image)
	assert.Equal(first.Url, "https://images.unsplash.com/photo-56789-098yhj?crop=entropy&cs=srgb&fm=jpg&ixid=OIFGHJIUGGH=rb-1.2.1&q=85")
	assert.Equal(first.Page, "https://unsplash.com/photos/blah_blah")
	assert.Equal(first.Caption, "ABC")
	user, ok := first.Meta["user"].(models.User)
	assert.True(ok)
	assert.Equal(user.Username, "aafgotiigg")
	assert.Equal(user.Name, "Larry Emeka")

	httpClient = new(internal.MockHttpClient)
	httpClient.On("Do", mock.Anything).Return(nil, errors.New("something went wrong")).Once()
	u = Unsplash{HttpClient: httpClient}
	images, err = u.fetchImages()
	assert.Nil(images)
	assert.NotNil(err)

	httpClient = new(internal.MockHttpClient)
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: io.NopCloser(strings.NewReader("Lorem"))}, nil).Once()
	u = Unsplash{HttpClient: httpClient}
	images, err = u.fetchImages()
	assert.Nil(images)
	assert.NotNil(err)
}

func TestUnsplashCache(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	assert := testify.New(t)
	httpClient := new(internal.MockHttpClient)
	db := new(database.MockDb)

	file, _ := os.Open("./testdata/unsplash_response.json")
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	db.On("Write", mock.Anything, mock.Anything).Return(nil).Once()
	var m = Unsplash{Db: db, HttpClient: httpClient}
	size, err := m.Cache(nil)
	assert.Nil(err)
	assert.Equal(1, size)

	file, _ = os.Open("./testdata/unsplash_response.json")
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	db.On("Write", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
	m = Unsplash{Db: db, HttpClient: httpClient}
	size, err = m.Cache(nil)
	assert.NotNil(err)
	assert.Equal(1, size)

	httpClient.On("Do", mock.Anything).Return(nil, errors.New("test")).Once()
	m = Unsplash{Db: db, HttpClient: httpClient}
	size, err = m.Cache(nil)
	assert.NotNil(err)
	assert.Equal(0, size)
}
