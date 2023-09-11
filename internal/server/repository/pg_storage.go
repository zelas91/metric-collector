package repository

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	db  *sql.DB
	ctx context.Context
	mem *MemStorage
}

func NewDBStorage(ctx context.Context, dbURL string) *DBStorage {
	db := newPostgresDB(dbURL)
	return &DBStorage{ctx: ctx, db: db, mem: NewMemStorage()}
}

func newPostgresDB(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("db open err : %v", err)
	}
	return db
}

func (d *DBStorage) Shutdown() error {
	return d.db.Close()
}

func (d *DBStorage) AddMetric(metric Metric) *Metric {
	return d.mem.AddMetric(metric)
}

func (d *DBStorage) GetMetric(name string) (*Metric, error) {
	return d.mem.GetMetric(name)
}

func (d *DBStorage) GetMetrics() []Metric {
	return d.mem.GetMetrics()
}

func (d *DBStorage) Ping() error {
	return d.db.Ping()
}
