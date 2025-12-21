package user

import (
	"context"
	"log/slog"
)

type Service interface {
	FindOne(ctx context.Context, cred map[string]any) (*UserOutput, error)
	Create(ctx context.Context, input LoginInput) (*UserOutput, error)
}

type service struct {
	repo Repository
}

// Constructor
func NewService(repo Repository) Service {
	slog.Info("news service initialized")
	return &service{repo: repo}
}

func (s *service) FindOne(ctx context.Context, cred map[string]any) (*UserOutput, error) {
	user, err := s.repo.FindOne(ctx, cred)
	if err != nil {
		slog.Error("service: failed to fetch user", slog.Any("error", err))
		return nil, err
	}
	slog.Info("service: user fetched successfully")
	return user, nil
}

func (s *service) Create(ctx context.Context, input LoginInput) (*UserOutput, error) {
	slog.Debug("service: creating user", slog.String("username", input.Username))

	// Register user on DB
	user, err := s.repo.Create(ctx, input)
	if err != nil {
		slog.Error("service: failed to create user", slog.String("username", input.Username), slog.Any("error", err))
		return nil, err
	}
	slog.Info("service: user created successfully", slog.String("username", input.Username))
	return user, nil
}
