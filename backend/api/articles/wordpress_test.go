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
	"time"
)

func TestWordPressFetchArticles(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	assert := testify.New(t)
	file, _ := os.Open("./testdata/wordpress_response.json")
	httpClient := new(internal.MockHttpClient)

	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	w := WordPress{HttpClient: httpClient}
	articles, err := w.fetchArticles()
	assert.Nil(err)
	assert.Equal(len(articles), 2)
	first := articles[0].(models.Article)
	second := articles[1].(models.Article)
	assert.Equal(first.Title, "Lorem ipsum is placeholder text commonly used")
	assert.Equal(second.UpdatedOn, time.Date(2018, time.June, 13, 15, 0, 55, 0, time.UTC))
	assert.Equal(second.Url, "https://www.lorem.com/54321/lorem-dolor-blatty-blah")
	assert.Equal(first.Thumbnail, "https://lorem.com/dolor/s.jpg")

	httpClient = new(internal.MockHttpClient)
	httpClient.On("Do", mock.Anything).Return(nil, errors.New("something went wrong")).Once()
	w = WordPress{HttpClient: httpClient}
	articles, err = w.fetchArticles()
	assert.Nil(articles)
	assert.NotNil(err)

	httpClient = new(internal.MockHttpClient)
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: io.NopCloser(strings.NewReader("Lorem"))}, nil).Once()
	w = WordPress{HttpClient: httpClient}
	articles, err = w.fetchArticles()
	assert.Nil(articles)
	assert.NotNil(err)
}

func TestWordPressCache(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	assert := testify.New(t)
	httpClient := new(internal.MockHttpClient)
	db := new(database.MockDb)

	file, _ := os.Open("./testdata/wordpress_response.json")
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	db.On("Write", mock.Anything, mock.Anything).Return(nil).Once()
	var m = WordPress{Db: db, HttpClient: httpClient}
	size, err := m.Cache(nil)
	assert.Nil(err)
	assert.Equal(2, size)

	file, _ = os.Open("./testdata/wordpress_response.json")
	httpClient.On("Do", mock.Anything).Return(&http.Response{Body: file}, nil).Once()
	db.On("Write", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
	m = WordPress{Db: db, HttpClient: httpClient}
	size, err = m.Cache(nil)
	assert.NotNil(err)
	assert.Equal(2, size)

	httpClient.On("Do", mock.Anything).Return(nil, errors.New("test")).Once()
	m = WordPress{Db: db, HttpClient: httpClient}
	size, err = m.Cache(nil)
	assert.NotNil(err)
	assert.Equal(0, size)
}
