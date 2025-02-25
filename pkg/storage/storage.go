package storage

import (
	"errors"
	"strings"

	"github.com/lib/pq"           // PostgreSQL драйвер
	"github.com/mattn/go-sqlite3" // SQLite драйвер
	"github.com/go-sql-driver/mysql" // MySQL драйвер
)

func IsDuplicatePhoneError(err error) bool {
	var sqliteErr sqlite3.Error
	var pqErr *pq.Error
	var mySQLErr *mysql.MySQLError

	switch {
	case errors.As(err, &pqErr) && pqErr.Code == "23505": // PostgreSQL
		return true
	case errors.As(err, &sqliteErr) && strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed"): // SQLite
		return true
	case errors.As(err, &mySQLErr) && mySQLErr.Number == 1062: // MySQL
		return true
	default:
		return false
	}
}