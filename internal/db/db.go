package db

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	maxIdleConns = 5
	maxOpenConns = 5
	maxConnLife  = 5 * time.Minute
)

var db *sql.DB

// GetDB initializes the connection to the database on first call, and
// returns the existing connection after that
func GetDB(dbDSN string) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	dbConnection, err := sql.Open("mysql", dbDSN)
	db = dbConnection

	if err != nil {
		return nil, fmt.Errorf("db.connect: %v", err)
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(maxConnLife)

	return db, nil
}
