package fxclickhouse_test

import (
	"testing"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ankorstore/yokai-contrib/fxclickhouse"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxClickHouseModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxclickhouse.FxClickhouseModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestFxClickHouseConn(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var conf *config.Config
	var conn driver.Conn

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxclickhouse.FxClickhouseModule,
		fx.Populate(&conf, &conn),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	assert.Equal(t, conf.GetInt("modules.clickhouse.maxOpenConns"), conn.Stats().MaxOpenConns)
	assert.Equal(t, conf.GetInt("modules.clickhouse.maxIdleConns"), conn.Stats().MaxIdleConns)
	assert.NoError(t, app.Err())
	assert.NotNil(t, conn)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}
