package healthcheck_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxclickhouse"
	"github.com/ankorstore/yokai-contrib/fxclickhouse/healthcheck"
	cmock "github.com/srikanthccv/ClickHouse-go-mock"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	t.Parallel()

	conn, err := cmock.NewClickHouseNative(nil)
	assert.NoError(t, err)

	probe := healthcheck.NewClickhouseProbe(conn)

	assert.Equal(t, "clickhouse", probe.Name())
}

func TestSetName(t *testing.T) {
	t.Parallel()

	conn, err := cmock.NewClickHouseNative(nil)
	assert.NoError(t, err)

	probe := healthcheck.NewClickhouseProbe(conn)
	probe.SetName("custom")

	assert.Equal(t, "custom", probe.Name())
}

func TestCheckSuccess(t *testing.T) {
	t.Parallel()

	conn, err := cmock.NewClickHouseNative(nil)
	assert.NoError(t, err)

	probe := healthcheck.NewClickhouseProbe(conn)

	result := probe.Check(context.Background())
	assert.True(t, result.Success)
	assert.Equal(t, "database ping success", result.Message)
}

func TestCheckFailure(t *testing.T) {
	t.Parallel()

	conn, err := fxclickhouse.NewDefaultClickhouseFactory().Create()
	assert.NoError(t, err)

	probe := healthcheck.NewClickhouseProbe(conn)

	result := probe.Check(context.Background())
	assert.False(t, result.Success)
	assert.Contains(t, result.Message, "database ping error:")
}
