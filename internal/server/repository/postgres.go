package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewPostgresDB(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("db open err : %v", err)
	}
	return db
}
