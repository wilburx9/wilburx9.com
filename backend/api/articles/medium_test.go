package articles

import (
	"errors"
	"github.com/sirupsen/logrus"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wilburt/wilburx9.dev/backend/api/articles/internal/models"
	"github.com/wilburt/wilburx9.dev/backend/api/internal"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestMediumFetchArticles(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	assert := testify.New(t)
	file, _ := os.Open("./testdata/medium_response.xml")
	httpClient := new(internal.MockHttpClient)

	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	var m = Medium{HttpClient: httpClient}
	articles, err := m.fetchArticles()
	assert.Nil(err)
	assert.Equal(2, len(articles))
	first := articles[0].(models.Article)
	second := articles[1].(models.Article)
	assert.Equal(first.Title, "Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	assert.Equal(first.Thumbnail, "https://cdn-images-1.medium.com/max/960/1*bbkcrsggiQLxNDRAgHiSBQ.png")
	assert.Equal(first.Url, "https://medium.com/lorem/lorem---q")
	assert.Equal(first.ID, "f6857d90bdcadcc4a34f8ea0712dd2ef")
	assert.NotEmpty(first.Excerpt)
	assert.Empty(second.Thumbnail)

	httpClient = new(internal.MockHttpClient)
	httpClient.On("Do", mock.Anything).Return(nil, errors.New("something went wrong")).Once()
	m = Medium{HttpClient: httpClient}
	articles, err = m.fetchArticles()
	assert.Nil(articles)
	assert.NotNil(err)

	httpClient = new(internal.MockHttpClient)
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: io.NopCloser(strings.NewReader("Lorem"))}, nil).Once()
	m = Medium{HttpClient: httpClient}
	articles, err = m.fetchArticles()
	assert.Nil(articles)
	assert.NotNil(err)
}

func TestMediumCache(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	assert := testify.New(t)
	httpClient := new(internal.MockHttpClient)
	db := new(database.MockDb)

	file, _ := os.Open("./testdata/medium_response.xml")
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	db.On("Write", mock.Anything, mock.Anything).Return(nil).Once()
	var m = Medium{Db: db, HttpClient: httpClient}
	size, err := m.Cache()
	assert.Nil(err)
	assert.Equal(2, size)

	file, _ = os.Open("./testdata/medium_response.xml")
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	db.On("Write", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
	m = Medium{Db: db, HttpClient: httpClient}
	size, err = m.Cache()
	assert.NotNil(err)
	assert.Equal(2, size)

	httpClient.On("Do", mock.Anything).Return(nil, errors.New("test")).Once()
	m = Medium{Db: db, HttpClient: httpClient}
	size, err = m.Cache()
	assert.NotNil(err)
	assert.Equal(0, size)
}
