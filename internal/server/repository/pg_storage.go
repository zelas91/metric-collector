package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	db  *sql.DB
	ctx context.Context
}

func NewDBStorage(ctx context.Context, dbURL string) *DBStorage {
	db := newPostgresDB(dbURL)
	if err := migration(db); err != nil {
		log.Fatal(err)
	}
	return &DBStorage{ctx: ctx, db: db}
}

func migration(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration error WithInstance, err:%v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./schema",
		"metrics", driver)
	if err != nil {
		return fmt.Errorf("migration NewWithDatabaseInstance, err:%v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration UP err : %v", err)
	}
	return nil
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
	isMetric, err := d.isMetricExists(ctx, metric.ID)
	if err != nil {
		log.Errorf("add metric err: %v", err)
		return nil
	}

	if isMetric {
		_, err = d.db.ExecContext(ctx, "update metrics set value=$1, delta=delta+$2 where name=$3",
			metric.Value, metric.Delta, metric.ID)
		if err != nil {
			log.Errorf("update metric err: %v", err)
			return nil
		}
	} else {
		_, err = d.db.ExecContext(ctx, "INSERT INTO metrics (name, type, value, delta) SELECT $1, id, $2 , $3 FROM metric_type WHERE name =$4",
			metric.ID, metric.Value, metric.Delta, metric.MType)
		if err != nil {
			log.Errorf("add metric err: %v", err)
			return nil
		}
	}
	return &metric
}

func (d *DBStorage) GetMetric(ctx context.Context, name string) (*Metric, error) {
	row := d.db.QueryRowContext(ctx, "select name, (select name from metric_type where id=type)as type, delta, value from metrics where name=$1", name)
	var metric Metric
	if err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
		log.Errorf("get metrics row err:%v", err)
		return nil, err
	}
	return &metric, nil
}

func (d *DBStorage) GetMetrics(ctx context.Context) []Metric {
	var metrics []Metric
	rows, err := d.db.QueryContext(ctx, "select name , (select name from metric_type where id=type)  as type, value ,delta from metrics")
	if err != nil {
		log.Errorf("get metrics query err: %v", err)
		return nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Errorf("rows close err :%v", err)
			return
		}
	}()
	for rows.Next() {
		var metric Metric
		if err = rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta); err != nil {
			log.Errorf("rows scan err: %v", err)
			return nil
		}
		metrics = append(metrics, metric)
	}
	if err = rows.Err(); err != nil {
		log.Errorf("rows err: %v", err)
		return nil
	}
	return metrics
}

func (d *DBStorage) Ping() error {
	return d.db.Ping()
}
func (d *DBStorage) isMetricExists(ctx context.Context, name string) (bool, error) {
	var isMetric bool
	if err := d.db.QueryRowContext(ctx, "select exists (select 1 from metrics where name=$1)", name).
		Scan(&isMetric); err != nil {
		return false, fmt.Errorf("is metric err :%w", err)
	}
	return isMetric, nil
}

func (d *DBStorage) AddMetrics(ctx context.Context, metrics []Metric) error {

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("new transactional err: %w ", err)
	}

	add, err := tx.PrepareContext(ctx, "INSERT INTO metrics (name, type, value, delta) SELECT $1, id, $2 , $3 FROM metric_type WHERE name =$4")
	if err != nil {
		return fmt.Errorf("add metrics, add prepare err: %w", err)
	}

	update, err := tx.PrepareContext(ctx, "update metrics set value=$1, delta=delta+$2 where name=$3")
	if err != nil {
		return fmt.Errorf("add metrics, update prepare err: %w", err)
	}

	existsMetric, err := tx.PrepareContext(ctx, "select exists (select 1 from metrics where name=$1)")
	if err != nil {
		return fmt.Errorf("add metrics, exists prepare err: %w", err)
	}

	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(sql.ErrTxDone, err) {
			log.Errorf("Rollback err: %v", err)
			return
		}
	}()

	var isMetrics bool

	for _, metric := range metrics {
		if err = existsMetric.QueryRowContext(ctx, metric.ID).Scan(&isMetrics); err != nil {
			return fmt.Errorf("query exists err: %w", err)
		}

		if isMetrics {
			if _, err = update.ExecContext(ctx, metric.Value, metric.Delta, metric.ID); err != nil {
				return fmt.Errorf("query update err: %w", err)
			}
		} else {
			if _, err = add.ExecContext(ctx, metric.ID, metric.Value, metric.Delta, metric.MType); err != nil {
				return fmt.Errorf("query update err: %w", err)
			}
		}
	}
	return tx.Commit()
}
