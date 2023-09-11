package repository

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	db  *sql.DB
	ctx context.Context
	mem *MemStorage
}

func NewDBStorage(ctx context.Context, dbURL string) *DBStorage {
	db := newPostgresDB(dbURL)
	migration(db)
	return &DBStorage{ctx: ctx, db: db, mem: NewMemStorage()}
}

func migration(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Errorf("migration error WithInstance, err:%v", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./schema",
		"metrics", driver)
	if err != nil {
		log.Errorf("migration NewWithDatabaseInstance, err:%v", err)
		return
	}
	if err := m.Up(); err != nil {
		log.Errorf("migration UP err : %v", err)
	}
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
