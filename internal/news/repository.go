package news

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

// ErrNotFound is used when it is not possible to find the requested News.
var ErrArticleNotFound = errors.New("article not found")

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

	const query = "SELECT title, body, datetime FROM news WHERE id=$1;"

	// Get a single row from the database (the first one) and copy the fetched data into the Article struct
	var article Article
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&article.Title,
		&article.Body,
		&article.Datetime,
	)

	// Check error returned by the query
	if errors.Is(err, sql.ErrNoRows) {
		slog.Warn("article not found", slog.String("id", id))
		return nil, ErrArticleNotFound
	}
	if err != nil {
		slog.Error("error fetching article", slog.String("id", id), slog.Any("error", err))
		return nil, err
	}

	slog.Info("successfully fetched article", slog.String("id", id))
	return &article, nil
}
