package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Conn *sql.DB
}

func New(dsn string) (*DB, error) {
	conn, err := sql.Open("mysql", dsn)
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{Conn: conn}, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
