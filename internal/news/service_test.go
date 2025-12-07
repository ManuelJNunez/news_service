package news

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubRepository struct {
	article *Article
	err     error
	called  bool
	lastID  string
}

func (s *stubRepository) GetByID(_ context.Context, id string) (*Article, error) {
	s.called = true
	s.lastID = id
	return s.article, s.err
}

func TestServiceGetByIDSuccess(t *testing.T) {
	article := &Article{Title: "fake_title", Body: "fake_body", Datetime: time.Now()}
	repo := &stubRepository{article: article}

	svc := NewService(repo)
	got, err := svc.GetByID(context.Background(), "123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.called || repo.lastID != "123" {
		t.Fatalf("repository was not called correctly: called=%v id=%s", repo.called, repo.lastID)
	}
	if got != article {
		t.Fatalf("expected article pointer to match, got %+v", got)
	}
}

func TestServiceGetByIDError(t *testing.T) {
	repo := &stubRepository{err: errors.New("failed to fetch article")}
	svc := NewService(repo)

	_, err := svc.GetByID(context.Background(), "abc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !repo.called {
		t.Fatalf("repository should have been called")
	}
}
