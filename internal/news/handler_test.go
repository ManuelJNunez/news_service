package news

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
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

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["error"] != "id parameter is required" {
		t.Fatalf("unexpected error message: %v", resp)
	}
}

func TestHandlerGetNewsNotFound(t *testing.T) {
	router := setupRouter(&stubService{err: ErrNewsNotFound})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/news?id=1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "news not found" {
		t.Fatalf("unexpected error message: %v", resp)
	}
}

func TestHandlerGetNewsInternalError(t *testing.T) {
	router := setupRouter(&stubService{err: errors.New("db down")})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/news?id=1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "internal error" {
		t.Fatalf("unexpected error message: %v", resp)
	}
}

func TestHandlerGetNewsSuccess(t *testing.T) {
	article := &Article{Title: "fake_title", Body: "fake_body", Datetime: time.Now()}
	router := setupRouter(&stubService{article: article})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/news?id=99", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "fake_title|fake_body") {
		t.Fatalf("expected rendered template to contain article data, got %q", body)
	}
}
