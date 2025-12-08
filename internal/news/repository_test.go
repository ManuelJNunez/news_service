package news

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close() //nolint:errcheck
	})

	repo := NewRepository(db)

	assert.NotNil(t, repo)
	assert.IsType(t, &postgresRepository{}, repo)
}

func TestGetByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close() //nolint:errcheck
	})

	repo := NewRepository(db)

	expectedArticle := &Article{
		Title:    "Test Article",
		Body:     "Test Body",
		Datetime: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"title", "body", "datetime"}).
		AddRow(expectedArticle.Title, expectedArticle.Body, expectedArticle.Datetime)

	mock.ExpectQuery("SELECT title, body, datetime FROM news WHERE id=(.+)").
		WillReturnRows(rows)

	article, err := repo.GetByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.NotNil(t, article)
	assert.Equal(t, expectedArticle.Title, article.Title)
	assert.Equal(t, expectedArticle.Body, article.Body)
	assert.WithinDuration(t, expectedArticle.Datetime, article.Datetime, time.Second)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close() //nolint:errcheck
	})

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT title, body, datetime FROM news WHERE id=(.+)").
		WillReturnError(sql.ErrNoRows)

	article, err := repo.GetByID(context.Background(), "999")

	assert.Error(t, err)
	assert.Nil(t, article)
	assert.Equal(t, ErrNewsNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close() //nolint:errcheck
	})

	repo := NewRepository(db)

	expectedError := errors.New("database connection error")

	mock.ExpectQuery("SELECT title, body, datetime FROM news WHERE id=(.+)").
		WillReturnError(expectedError)

	article, err := repo.GetByID(context.Background(), "1")

	assert.Error(t, err)
	assert.Nil(t, article)
	assert.Equal(t, expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
