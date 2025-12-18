package news

import (
	"context"
	"log/slog"
)

type Service interface {
	GetByID(ctx context.Context, id uint64) (*Article, error)
}

type service struct {
	repo Repository
}

// Constructor
func NewService(repo Repository) Service {
	slog.Info("news service initialized")
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, id uint64) (*Article, error) {
	slog.Debug("service: fetching article", slog.Uint64("id", id))
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		slog.Error("service: failed to fetch article", slog.Uint64("id", id), slog.Any("error", err))
		return nil, err
	}
	slog.Info("service: article fetched successfully", slog.Uint64("id", id))
	return article, nil
}
