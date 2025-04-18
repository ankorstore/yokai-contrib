package fxclickhouse

import (
	clickhousesdk "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// ClickhouseFactory is the interface for [gorm.DB] factories.
type ClickhouseFactory interface {
	Create(options ...ClickhouseOptions) (driver.Conn, error)
}

// DefaultClickhouseFactory is the default [ClickhouseFactory] implementation.
type DefaultClickhouseFactory struct{}

// NewDefaultClickhouseFactory returns a [DefaultClickhouseFactory], implementing [ClickhouseFactory].
func NewDefaultClickhouseFactory() ClickhouseFactory {
	return &DefaultClickhouseFactory{}
}

func (f *DefaultClickhouseFactory) Create(options ...ClickhouseOptions) (driver.Conn, error) {
	appliedOpts := DefaultClickhouseOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	conn, err := f.createDatabase(appliedOpts)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (f *DefaultClickhouseFactory) createDatabase(options Options) (driver.Conn, error) {
	conn, err := clickhousesdk.Open(&clickhousesdk.Options{
		Addr:             options.Addr,
		Auth:             options.Auth,
		Debug:            options.Debug,
		Settings:         options.Settings,
		Compression:      options.Compression,
		DialTimeout:      options.DialTimeout,
		MaxOpenConns:     options.MaxOpenConns,
		MaxIdleConns:     options.MaxIdleConns,
		ConnMaxLifetime:  options.ConnMaxLifetime,
		ConnOpenStrategy: options.ConnOpenStrategy,
		BlockBufferSize:  options.BlockBufferSize,
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
