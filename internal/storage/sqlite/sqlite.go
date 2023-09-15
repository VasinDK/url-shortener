package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"mod_shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"
)

type Storage struct {
	db *sql.DB
}

type Url struct {
	url   string
	alias string
}

func New(storagePath string) (*Storage, error) {
	op := "Storage.Sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX idx_alias ON url (alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "Storage.sqlite.saveUrl"
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)

	if err != nil {

		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s %w", op, err)
	}

	id, err := res.LastInsertId()

	if err != err {
		return 0, fmt.Errorf("%s %w", "last id", err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "Storage.SQLite.GetURL"

	var resURL Url

	stmt, err := s.db.Prepare("SELECT url, alias FROM url WHERE alias = ?")

	if err != nil {
		return "", fmt.Errorf("%s %w", op, err)
	}

	err = stmt.QueryRow(alias).Scan(&resURL.url, &resURL.alias)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s %w", op, err)
	}

	slog.Info("resURL", "структура", resURL)

	return resURL.url, nil
}

func (s *Storage) DeleteURL(id int64) (int64, error) {
	const op = "Sqlite.Storage.DeleteURL"

	stmt, err := s.db.Prepare("Delete FROM url WHERE id = ?")

	res, err := stmt.Exec(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrElemNotFount
		}

		return 0, fmt.Errorf("%s %w", op, err)
	}

	if rows, err := res.RowsAffected(); rows > 0 && err == nil {
		return rows, nil
	}

	return 0, fmt.Errorf("%s %w", op, err)
}
