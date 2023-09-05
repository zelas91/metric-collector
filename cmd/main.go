package main

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite"
	"time"
)

func main() {
	// при использовании пакета go-sqlite3 имя драйвера — sqlite3
	db, err := sql.Open("sqlite", "video.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}
	// ...
}
