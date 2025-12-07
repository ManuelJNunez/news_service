package news

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

// ErrNotFound is used when it is not possible to find the requested News.
var ErrNewsNotFound = errors.New("news not found")

type Repository interface {
	GetByID(ctx context.Context, id string) (*Article, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (s *postgresRepository) GetByID(ctx context.Context, id string) (*Article, error) {
	slog.Debug("fetching news", slog.String("id", id))
	query := fmt.Sprintf("SELECT title, body, datetime FROM news WHERE id=%s;", id)

	// Get a single row from the database (the first one) and copy the fetched data into the Article struct
	var n Article
	err := s.db.QueryRowContext(ctx, query).Scan(
		&n.Title,
		&n.Body,
		&n.Datetime,
	)

	if errors.Is(err, sql.ErrNoRows) {
		slog.Warn("news not found", slog.String("id", id))
		return nil, ErrNewsNotFound
	}
	if err != nil {
		slog.Error("error fetching news", slog.String("id", id), slog.Any("error", err))
		return nil, err
	}

	slog.Info("successfully fetched news", slog.String("id", id))
	return &n, nil
}
