package news

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stubRepository struct {
	article *Article
	err     error
	called  bool
	lastID  uint64
}

func (s *stubRepository) GetByID(_ context.Context, id uint64) (*Article, error) {
	s.called = true
	s.lastID = id
	return s.article, s.err
}

func TestServiceGetByIDSuccess(t *testing.T) {
	article := &Article{Title: "fake_title", Body: "fake_body", Datetime: time.Now()}
	repo := &stubRepository{article: article}

	svc := NewService(repo)
	got, err := svc.GetByID(context.Background(), 123)

	assert.NoError(t, err)
	assert.True(t, repo.called)
	assert.Equal(t, uint64(123), repo.lastID)
	assert.Equal(t, article, got)
}

func TestServiceGetByIDError(t *testing.T) {
	repo := &stubRepository{err: errors.New("failed to fetch article")}
	svc := NewService(repo)

	_, err := svc.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.True(t, repo.called)
}
