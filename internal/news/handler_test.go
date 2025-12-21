package news

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type stubService struct {
	article *Article
	err     error
}

func (s *stubService) GetByID(_ context.Context, _ string) (*Article, error) {
	return s.article, s.err
}

func setupRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Create a fake HTML template
	r.SetHTMLTemplate(template.Must(template.New("article.html").Parse("{{.Title}}|{{.Body}}")))
	RegisterRoutes(r.Group(""), NewHandler(svc))
	return r
}

func TestHandlerGetNewsMissingID(t *testing.T) {
	router := setupRouter(&stubService{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/news", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "id parameter is required", resp["error"])
}

func TestHandlerGetNewsNotFound(t *testing.T) {
	router := setupRouter(&stubService{err: ErrArticleNotFound})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/news?id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "article not found", resp["error"])
}

func TestHandlerGetNewsInternalError(t *testing.T) {
	router := setupRouter(&stubService{err: errors.New("db down")})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/news?id=1", nil)
	router.ServeHTTP(w, req)

	// Internal errors should return 404 to avoid leaking information
	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "article not found", resp["error"])
}

func TestHandlerGetNewsSuccess(t *testing.T) {
	article := &Article{Title: "fake_title", Body: "fake_body", Datetime: time.Now()}
	router := setupRouter(&stubService{article: article})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/news?id=99", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, "fake_title|fake_body")
}
