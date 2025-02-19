package storage

import "errors"

type Storage interface {
	SaveURL(url string, alias string) error
	GetURL(alias string) (string, error)
}

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url already exists")
	ErrAliasExists = errors.New("alias already exists")
)
