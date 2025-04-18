package fxclickhouse_test

import (
	"testing"
	"time"

	clickhousesdk "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ankorstore/yokai-contrib/fxclickhouse"
	"github.com/stretchr/testify/assert"
)

func TestWithAddr(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithAddr([]string{"node1", "node2"})(&opt)

	assert.Equal(t, []string{"node1", "node2"}, opt.Addr)
}

func TestWithAuth(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithAuth(clickhousesdk.Auth{
		Database: "db_1",
		Username: "usr_1",
		Password: "pwd_1",
	})(&opt)

	assert.Equal(t, clickhousesdk.Auth{
		Database: "db_1",
		Username: "usr_1",
		Password: "pwd_1",
	}, opt.Auth)
}

func TestWithDebug(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithDebug(true)(&opt)

	assert.Equal(t, true, opt.Debug)
}

func TestWithSettings(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithSettings(clickhousesdk.Settings{
		"foo": "bar",
	})(&opt)

	assert.Equal(t, clickhousesdk.Settings{
		"foo": "bar",
	}, opt.Settings)
}

func TestWithCompression(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithCompression(&clickhousesdk.Compression{
		Method: clickhousesdk.CompressionGZIP,
	})(&opt)

	assert.Equal(t, &clickhousesdk.Compression{
		Method: clickhousesdk.CompressionGZIP,
	}, opt.Compression)
}

func TestWithDialTimeout(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithDialTimeout(time.Duration(5) * time.Second)(&opt)

	assert.Equal(t, time.Duration(5)*time.Second, opt.DialTimeout)
}

func TestWithMaxOpenConns(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithMaxOpenConns(99)(&opt)

	assert.Equal(t, 99, opt.MaxOpenConns)
}

func TestWithMaxIdleConns(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithMaxIdleConns(99)(&opt)

	assert.Equal(t, 99, opt.MaxIdleConns)
}

func TestWithConnMaxLifetime(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithConnMaxLifetime(time.Duration(5) * time.Minute)(&opt)

	assert.Equal(t, time.Duration(5)*time.Minute, opt.ConnMaxLifetime)
}

func TestWithConnOpenStrategy(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithConnOpenStrategy(clickhousesdk.ConnOpenRoundRobin)(&opt)

	assert.Equal(t, clickhousesdk.ConnOpenRoundRobin, opt.ConnOpenStrategy)
}

func TestWithBlockBufferSize(t *testing.T) {
	t.Parallel()

	opt := fxclickhouse.DefaultClickhouseOptions()
	fxclickhouse.WithBlockBufferSize(1)(&opt)

	assert.Equal(t, uint8(1), opt.BlockBufferSize)
}
