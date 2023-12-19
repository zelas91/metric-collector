// Package repository to work with data.

package repository

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
)

var (
	log  = logger.New()
	once sync.Once
)

// FileStorage to work with data in file.
type FileStorage struct {
	file *os.File
	mem  *MemStorage
	cfg  *config.Config
	ctx  context.Context
}

// NewFileStorage make FileStorage struct.
func NewFileStorage(ctx context.Context, cfg *config.Config) *FileStorage {
	storage := &FileStorage{ctx: ctx, cfg: cfg}

	path := *cfg.FilePath

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("open file err: %v", err)
		return nil
	}
	storage.file = file
	var mem *MemStorage
	if memStorage := storage.getMetricsFile(); memStorage == nil {
		mem = NewMemStorage()
	} else {
		mem = memStorage
	}
	storage.mem = mem
	return storage
}

// Shutdown close file.
func (f *FileStorage) Shutdown() error {
	return f.file.Close()
}

func (f *FileStorage) AddMetric(ctx context.Context, metric Metric) *Metric {
	resultMetrics := f.mem.AddMetric(ctx, metric)
	f.saveMetric()
	return resultMetrics
}

func (f *FileStorage) GetMetric(ctx context.Context, name string) (*Metric, error) {
	metric, err := f.mem.GetMetric(ctx, name)
	return metric, err
}

func (f *FileStorage) GetMetrics(ctx context.Context) []Metric {
	metrics := f.mem.GetMetrics(ctx)
	return metrics
}

func (f *FileStorage) save() error {
	metrics := f.mem.GetMetrics(context.Background())
	data, err := json.MarshalIndent(metrics, "", " ")
	if err != nil {
		return err
	}

	if err = f.file.Truncate(0); err != nil {
		return err
	}
	if _, err = f.file.Seek(0, 0); err != nil {
		return err
	}
	if _, err = f.file.Write(data); err != nil {
		return err
	}
	return nil
}
func (f *FileStorage) asyncSave() {
	once.Do(func() {
		go func() {
			tickerStoreInterval := time.NewTicker(time.Duration(*f.cfg.StoreInterval) * time.Second)
			for {
				select {
				case <-tickerStoreInterval.C:
					if err := f.save(); err != nil {
						log.Errorf("save error %v", err)
					}
				case <-f.ctx.Done():
					return
				}

			}
		}()
	})
}

func (f *FileStorage) syncSave() {
	if err := f.save(); err != nil {
		log.Errorf("save error %v", err)
		return
	}
}
func (f *FileStorage) saveMetric() {
	if f.cfg.StoreInterval == nil {
		return
	}

	if *f.cfg.StoreInterval == 0 {
		f.syncSave()
		return
	}
	f.asyncSave()
}

func (f *FileStorage) getMetricsFile() *MemStorage {
	info, err := f.file.Stat()
	if err != nil {
		return nil
	}
	data := make([]byte, info.Size())
	if _, err = f.file.Read(data); err != nil {
		log.Errorf("read file err: %v", err)
		return nil
	}
	var metrics []Metric
	if err = json.Unmarshal(data, &metrics); err != nil {
		log.Errorf("read metrics db err: %v", err)
		return nil
	}

	mem := make(map[string]Metric, len(metrics))
	for _, metric := range metrics {
		mem[metric.ID] = metric
	}
	return &MemStorage{mem: mem}
}

func (f *FileStorage) AddMetrics(ctx context.Context, metrics []Metric) error {
	for _, metric := range metrics {
		_ = f.AddMetric(ctx, metric)
	}
	return nil
}
