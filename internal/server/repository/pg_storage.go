package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/zelas91/metric-collector/internal/server/types"
)

var (
	gauge   int
	counter int
)

type DBStorage struct {
	db  *sql.DB
	ctx context.Context
}

func NewDBStorage(ctx context.Context, dbURL string) *DBStorage {
	db := newPostgresDB(dbURL)
	migration(db)
	initType(db)
	return &DBStorage{ctx: ctx, db: db}
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

func (d *DBStorage) AddMetric(ctx context.Context, metric Metric) *Metric {
	isMetric, err := d.isMetric(ctx, metric.ID)
	if err != nil {
		log.Errorf("add metric err: %v", err)
	}

	if isMetric {
		if err := d.updateMetric(ctx, metric); err != nil {
			log.Errorf("add metric err:%v", err)
			return nil
		}
	} else {
		if err := d.addMetric(ctx, metric); err != nil {
			log.Errorf("add metric err:%v", err)
			return nil
		}
	}
	return &metric
}

func (d *DBStorage) updateMetric(ctx context.Context, metric Metric) error {
	switch metric.MType {
	case types.GaugeType:
		_, err := d.db.ExecContext(ctx, "update metrics set value=$1 where name=$2", *metric.Value, metric.ID)
		if err != nil {
			return fmt.Errorf("update metric err:%w", err)
		}
	case types.CounterType:
		_, err := d.db.ExecContext(ctx, "update metrics set delta=$1 where name=$2", *metric.Delta, metric.ID)
		if err != nil {
			return fmt.Errorf("update metric err:%w", err)
		}

	}
	return nil
}

func (d *DBStorage) addMetric(ctx context.Context, metric Metric) error {
	switch metric.MType {
	case types.GaugeType:
		_, err := d.db.ExecContext(ctx, "insert into metrics (name , type , value) values ($1,$2,$3)",
			metric.ID, d.convertType(metric.MType), *metric.Value)
		if err != nil {
			return fmt.Errorf("add metric err:%w", err)
		}
	case types.CounterType:
		_, err := d.db.ExecContext(ctx, "insert into metrics (name , type , delta) values ($1,$2,$3)",
			metric.ID, d.convertType(metric.MType), *metric.Delta)
		if err != nil {
			return fmt.Errorf("add metric err:%w", err)
		}

	}
	return nil
}

func (d *DBStorage) GetMetric(ctx context.Context, name string) (*Metric, error) {
	row := d.db.QueryRowContext(ctx, "select name, type, delta, value from metrics where name=$1", name)
	var metric Metric
	if err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
		log.Infof("get metrics row err:%v", err)
	}
	return &metric, nil
}

func (d *DBStorage) GetMetrics(ctx context.Context) []Metric {
	return nil
}

func (d *DBStorage) Ping() error {
	return d.db.Ping()
}
func (d *DBStorage) isMetric(ctx context.Context, name string) (bool, error) {
	var isMetric bool
	if err := d.db.QueryRowContext(ctx, "select exists (select 1 from metrics where name=$1)", name).
		Scan(&isMetric); err != nil {
		return false, fmt.Errorf("is metric err :%w", err)
	}
	return isMetric, nil
}

func (d *DBStorage) convertType(t string) int {
	switch t {
	case types.GaugeType:
		return gauge
	case types.CounterType:
		return counter
	}
	return 0
}
func initType(db *sql.DB) {
	rows, err := db.Query("select id, name from metric_type")
	if err != nil {
		log.Errorf("init type err:%v", err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Errorf("init type rows close err:%v", err)
		}
	}()
	var id int
	var name string
	for rows.Next() {
		rows.Scan(&id, &name)
		switch name {
		case types.GaugeType:
			gauge = id
		case types.CounterType:
			counter = id
		}
		if rows.Err() != nil {
			log.Errorf("init type rows err:%v", err)
			return
		}
	}
}
