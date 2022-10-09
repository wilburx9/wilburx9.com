package articles

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wilburt/wilburx9.dev/backend/api/internal/database"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	assert := testify.New(t)
	db := new(database.MockDb)
	updatedAt := database.UpdatedAt{T: time.Now()}

	db.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(nil, updatedAt, errors.New("test")).Once()
	router := setupRouter(db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/articles", nil)
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusInternalServerError, w.Code)

	var articles []map[string]interface{}
	db.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(articles, updatedAt, nil).Once()
	router = setupRouter(db)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/articles", nil)
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusInternalServerError, w.Code)

	articles = []map[string]interface{}{{"Lorem": "Dolor X"}}
	db.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(articles, updatedAt, nil).Once()
	router = setupRouter(db)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/articles", nil)
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)
	assert.Contains(w.Body.String(), "Dolor X")

}

func setupRouter(db database.ReadWrite) *gin.Engine {
	r := gin.New()
	r.GET("/articles", func(c *gin.Context) { Handler(c, db) })
	return r
}
