package sqlite

import (
	"database/sql"
	"fmt"
	str "url-shortener/internal/lib/string"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"
	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url (
		    id INTEGER PRIMARY KEY,
		    alias VARCHAR (10) NOT NULL UNIQUE,
		    url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS index_alias ON url(alias)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Create(urlValue string) (int64, error) {
	const op = "storage.sqlite.Create"

	stmt, err := s.db.Prepare("INSERT INTO url (alias, url) VALUES (?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	// TODO Check that alias is not exists in db
	alias := str.Generate(6)

	res, err := stmt.Exec(alias, urlValue)

	if err != nil {
		return 0, fmt.Errorf("%s: exec statement: %w", op, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("%s: can't get last insert id: %w", op, err)
	}

	return id, nil
}
