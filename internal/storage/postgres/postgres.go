package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"url-service/internal/storage"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connectionString string) (*Postgres, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening pgType connection: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL UNIQUE,
		alias TEXT NOT NULL UNIQUE
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating table urls: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS urls_alias_idx ON urls(alias);`)
	if err != nil {
		return nil, fmt.Errorf("error creating index: %v", err)
	}

	return &Postgres{db: db}, nil
}

func (pg *Postgres) SaveURL(url, alias string) error {
	_, err := pg.db.Exec(`INSERT INTO urls(url, alias) VALUES($1, $2)`, url, alias)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "urls_url_key":
				return storage.ErrURLExists
			case "urls_alias_key":
				return storage.ErrAliasExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (pg *Postgres) GetURL(alias string) (string, error) {
	var url string
	err := pg.db.QueryRow(`SELECT url FROM urls WHERE alias = $1`, alias).Scan(&url)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return "", storage.ErrURLNotFound
		}
		return "", err
	}
	return url, nil
}

var _ storage.Storage = (*Postgres)(nil)
