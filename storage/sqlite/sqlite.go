package sqlite

import (
	"PracticeBot/storage"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// NewStorage makes new SQLite storage.
func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping storage: %w", err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	query := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, query, p.URL, p.UserName)
	if err != nil {
		return fmt.Errorf("failed to save page: %w", err)
	}

	return nil
}

// PickRandom picks random page from storage with reported userName.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	query := `SELECT url FROM pages WHERE user_name = ? ORDER BY random() LIMIT 1`

	result := s.db.QueryRowContext(ctx, query, userName)

	var u string
	err := result.Scan(&u)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to pick random page: %w", err)
	}

	return &storage.Page{URL: u, UserName: userName}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	query := `DELETE FROM pages WHERE url = ? and user_name = ?`
	_, err := s.db.ExecContext(ctx, query, p.URL, p.UserName)
	if err != nil {
		return fmt.Errorf("failed to remove page: %w", err)
	}

	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	query := `SELECT url FROM pages WHERE url = ? and user_name = ?`
	result := s.db.QueryRowContext(ctx, query, p.URL, p.UserName)
	var u string
	err := result.Scan(&u)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check page existence: %w", err)
	}

	return true, nil
}

func (s *Storage) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
