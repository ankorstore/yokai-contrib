package fxclickhouse

import (
	"errors"
	"strings"
	"time"

	clickhousesdk "github.com/ClickHouse/clickhouse-go/v2"
)

// Options are options for the [ClickhouseFactory] implementations.
type Options struct {
	Addr             []string
	Auth             clickhousesdk.Auth
	Settings         clickhousesdk.Settings
	Debug            bool
	Compression      *clickhousesdk.Compression
	DialTimeout      time.Duration
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
	ConnOpenStrategy clickhousesdk.ConnOpenStrategy
	BlockBufferSize  uint8
}

// DefaultClickhouseOptions are the default options used in the [DefaultClickhouseFactory].
func DefaultClickhouseOptions() Options {
	return Options{
		Addr:             []string{},
		Auth:             clickhousesdk.Auth{},
		Debug:            false,
		Settings:         clickhousesdk.Settings{},
		Compression:      nil,
		DialTimeout:      time.Duration(1) * time.Second,
		MaxOpenConns:     10,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Duration(1) * time.Hour,
		ConnOpenStrategy: clickhousesdk.ConnOpenInOrder,
		BlockBufferSize:  2,
	}
}

// ClickhouseOption are functional options for the [ClickhouseFactory] implementations.
type ClickhouseOptions func(o *Options)

// WithAddr is used to specify the database DSN to use.
func WithAddr(a []string) ClickhouseOptions {
	return func(o *Options) {
		o.Addr = a
	}
}

// WithAuth is used to specify the [clickhouse.Auth] to use.
func WithAuth(a clickhousesdk.Auth) ClickhouseOptions {
	return func(o *Options) {
		o.Auth = a
	}
}

// WithDebug is used to specify the [clickhouse.Auth] to use.
func WithDebug(d bool) ClickhouseOptions {
	return func(o *Options) {
		o.Debug = d
	}
}

// WithSettings is used to specify the [clickhouse.Settings] to use.
func WithSettings(s clickhousesdk.Settings) ClickhouseOptions {
	return func(o *Options) {
		o.Settings = s
	}
}

// WithCompression is used to specify the [clickhouse.Compression] to use.
func WithCompression(c *clickhousesdk.Compression) ClickhouseOptions {
	return func(o *Options) {
		o.Compression = c
	}
}

// WithDialTimeout is used to specify the database [DialTimeout] to use.
func WithDialTimeout(d time.Duration) ClickhouseOptions {
	return func(o *Options) {
		o.DialTimeout = d
	}
}

// WithMaxOpenConns is used to specify the database [MaxOpenConns] to use.
func WithMaxOpenConns(m int) ClickhouseOptions {
	return func(o *Options) {
		o.MaxOpenConns = m
	}
}

// WithMaxIdleConns is used to specify the database [MaxIdleConns] to use.
func WithMaxIdleConns(m int) ClickhouseOptions {
	return func(o *Options) {
		o.MaxIdleConns = m
	}
}

// WithConnMaxLifetime is used to specify the database [ConnMaxLifetime] to use.
func WithConnMaxLifetime(t time.Duration) ClickhouseOptions {
	return func(o *Options) {
		o.ConnMaxLifetime = t
	}
}

// WithConnOpenStrategy is used to specify [clickhouse.ConnOpenStrategy].
func WithConnOpenStrategy(s clickhousesdk.ConnOpenStrategy) ClickhouseOptions {
	return func(o *Options) {
		o.ConnOpenStrategy = s
	}
}

// WithBlockBufferSize is used to specify the database [BlockBufferSize] to use.
func WithBlockBufferSize(b uint8) ClickhouseOptions {
	return func(o *Options) {
		o.BlockBufferSize = b
	}
}

func FetchConnOpenStrategy(str string) (clickhousesdk.ConnOpenStrategy, error) {
	switch strings.ToLower(str) {
	case "inorder":
		return clickhousesdk.ConnOpenInOrder, nil
	case "random":
		return clickhousesdk.ConnOpenRandom, nil
	case "roundrobin":
		return clickhousesdk.ConnOpenRoundRobin, nil
	default:
		return clickhousesdk.ConnOpenInOrder, errors.New("Unknown ConnOpenStrategy")
	}
}
