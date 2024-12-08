package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
