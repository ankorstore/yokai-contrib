package fxclickhouse

import (
	"context"
	"errors"

	"github.com/ClickHouse/clickhouse-go/v2"
	clickhousesdk "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "clickhouse"

// FxClickhouseModule is the [Fx] orm module.
//
// [Fx]: https://github.com/uber-go/fx
var FxClickhouseModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewDefaultClickhouseFactory,
		NewFxClickhouse,
	),
)

// FxClickhouseParam allows injection of the required dependencies in [NewFxClickhouse].
type FxClickhouseParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Factory   ClickhouseFactory
	Config    *config.Config
}

// NewFxClickhouse returns a [gorm.DB].
func NewFxClickhouse(p FxClickhouseParam) (driver.Conn, error) {
	addrs := p.Config.GetStringSlice("modules.clickhouse.address")
	auth := clickhousesdk.Auth{
		Database: p.Config.GetString("modules.clickhouse.database"),
		Username: p.Config.GetString("modules.clickhouse.username"),
		Password: p.Config.GetString("modules.clickhouse.password"),
	}
	options := []ClickhouseOptions{
		WithAddr(addrs),
		WithAuth(auth),
	}

	if p.Config.IsSet("modules.clickhouse.debug") {
		options = append(options, WithDebug(p.Config.GetBool("modules.clickhouse.debug")))
	}

	if p.Config.IsSet("modules.clickhouse.maxOpenConns") {
		options = append(options, WithMaxOpenConns(p.Config.GetInt("modules.clickhouse.maxOpenConns")))
	}

	if p.Config.IsSet("modules.clickhouse.maxIdleConns") {
		options = append(options, WithMaxIdleConns(p.Config.GetInt("modules.clickhouse.maxidleConns")))
	}

	if p.Config.IsSet("modules.clickhouse.connOpenStrategy") {
		connOpenStrategy, err := FetchConnOpenStrategy(p.Config.GetString("modules.clickhouse.connOpenStrategy"))
		if err != nil {
			return nil, err
		}
		options = append(options, WithConnOpenStrategy(connOpenStrategy))
	}

	if p.Config.IsSet("modules.clickhouse.blockBufferSize") {
		options = append(options, WithBlockBufferSize(uint8(p.Config.GetUint("modules.clickhouse.blockBufferSize"))))
	}

	if p.Config.IsSet("modules.clickhouse.dialTimeout") {
		options = append(options, WithDialTimeout(p.Config.GetDuration("modules.clickhouse.dialTimeout")))
	}

	if p.Config.IsSet("modules.clickhouse.settings") {
		settings, ok := (p.Config.Get("modules.clickhouse.settings")).(clickhouse.Settings)
		if !ok {
			return nil, errors.New("Unknown settings")
		}

		options = append(options, WithSettings(settings))
	}

	conn, err := p.Factory.Create(options...)
	if err != nil {
		return nil, err
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if !p.Config.IsTestEnv() {
				conn.Close()
				return nil
			}

			return nil
		},
	})

	return conn, nil
}
