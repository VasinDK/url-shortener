package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"mod_shortener/internal/lib/api/user"
	"mod_shortener/internal/lib/crypto"
	"mod_shortener/internal/lib/logger/sl"
	"mod_shortener/internal/storage"
	"time"

	"log/slog"

	"github.com/mattn/go-sqlite3"
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

	// stmt, err = db.

	// DROP TABLE IF EXISTS user;
	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS user(
			id INTEGER PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			name TEXT,
			surname TEXT ,
			email TEXT UNIQUE,
			phone TEXT,
			pass TEXT NOT NULL,
			refresh_token TEXT UNIQUE,
			expires_at DATETIME
		);

		CREATE INDEX idx_login ON user(login);
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

func (s *Storage) AddUser(user *user.User, log *slog.Logger) (int64, error) {
	const op = "Sqlite.Storage.AddUser"

	stmt, err := s.db.Prepare(`
		INSERT INTO user(login, name, surname, email, phone, pass, refresh_token) 
		VALUES(?,?,?,?,?,?,?)
	`)

	if err != nil {
		log.Error(op+".Prepare", sl.Err(err))
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	user.Pass, _ = crypto.HashPass(user.Pass)
	user.Refresh_token, _ = crypto.HashPass(user.Refresh_token)

	res, err := stmt.Exec(
		user.Login,
		user.Name,
		user.Surname,
		user.Email,
		user.Phone,
		user.Pass,
		user.Refresh_token,
	)
	if err != nil {
		log.Error(op+".Exec", sl.Err(err))
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetPass(login string, log *slog.Logger) (string, string, string, error) {
	const op = "storage.sqlite.GetPass"

	stmt, err := s.db.Prepare(`
		SELECT id, pass, refresh_token 
		FROM user 
		WHERE login = ?
	`)

	if err != nil {
		log.Error(op+" Prepare", sl.Err(err))
	}

	var pass string
	var token string
	var id string

	err = stmt.QueryRow(login).Scan(&id, &pass, &token)

	if err != nil {
		log.Error(op+" Exec", sl.Err(err))
	}

	if errors.Is(err, sql.ErrNoRows) {
		return "", "", "", sql.ErrNoRows
	}

	return pass, token, id, nil
}

func (s *Storage) GetUser(login string, log *slog.Logger) (*user.User, error) {
	const op = "storage.sqlite.GetUser"

	stmt, err := s.db.Prepare(`
		SELECT id, login, name, surname, email, phone
		FROM user 
		WHERE login = ?
	`)

	if err != nil {
		log.Error(op+" Prepare", sl.Err(err))
	}

	userCurrent := &user.User{}

	err = stmt.QueryRow(login).Scan(
		&userCurrent.Id,
		&userCurrent.Login,
		&userCurrent.Name,
		&userCurrent.Surname,
		&userCurrent.Email,
		&userCurrent.Phone,
	)

	if err != nil {
		log.Error(op+" Exec", sl.Err(err))
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &user.User{}, sql.ErrNoRows
	}

	return userCurrent, nil
}

func (s *Storage) UpdateRefreshToken(RefToken string, ExpAT time.Time, log *slog.Logger, idUser string) (string, error) {
	const op = "storage.sqlite.UpdateRefreshToken"

	stmt, err := s.db.Prepare(`
		UPDATE user 
		SET refresh_token = ?, expires_at = ?
		WHERE id = ?
	`)

	if err != nil {
		log.Error(op, sl.Err(err))
		return "", err
	}

	res, err := stmt.Exec(RefToken, ExpAT, idUser)

	if err != nil {
		log.Error(op+".Exec", sl.Err(err))
		return "", err
	}

	return fmt.Sprint(res.LastInsertId()), nil
}
