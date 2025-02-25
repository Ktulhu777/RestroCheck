package postgresql

import (
	"database/sql"
	"fmt"

	"restrocheck/internal/config"
	"restrocheck/internal/storage"

	_ "github.com/lib/pq"
)

func BuildStringConnectDB(cfg *config.Config) string {
	var dbURL string = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	return dbURL
}


func NewStorage(cfg *config.Config) (*storage.Storage, error) {
	const fn = "internal.storage.postgresql.NewStorage"

	var dbURL string = BuildStringConnectDB(cfg)

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return nil, fmt.Errorf("%s: %s: %v", fn, storage.ErrOpenDBConnection, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %s: %v", fn, storage.ErrPingDB, err)
	}

	return storage.NewStorage(db), nil
}