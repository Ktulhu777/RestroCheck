package storage

import (
	"database/sql"
	"errors"
)

type Storage struct {
	DB *sql.DB
}

var (
	ErrOpenDBConnection = errors.New("failed to open database connection")
	ErrPingDB           = errors.New("failed to ping database")
)


func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) Close() error {
	if err := s.DB.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Ping() error {
	if err := s.DB.Ping(); err != nil {
		return ErrPingDB
	}
	return nil
}