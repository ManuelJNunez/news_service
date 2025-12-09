package news

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
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

	// Block dangerous SQL keywords
	if containsSQLKeywords(id) {
		slog.Warn("blocked SQLi attempt", slog.String("id", id))
		return nil, ErrNewsNotFound
	}

	query := fmt.Sprintf("SELECT title, body, datetime FROM news WHERE id=%s;", id)

	// Get a single row from the database (the first one) and copy the fetched data into the Article struct
	var article Article
	err := s.db.QueryRowContext(ctx, query).Scan(
		&article.Title,
		&article.Body,
		&article.Datetime,
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
	return &article, nil
}

func containsSQLKeywords(input string) bool {
	dangerousKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "ORDER", "GROUP", ";", "/*", "*/", "--"}
	inputUpper := strings.ToUpper(input)

	for _, keyword := range dangerousKeywords {
		if strings.Contains(inputUpper, keyword) {
			return true
		}
	}
	return false
}
