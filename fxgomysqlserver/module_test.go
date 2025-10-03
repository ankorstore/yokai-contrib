package fxgomysqlserver_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/testdata/transport"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/trace/tracetest"
	sqle "github.com/dolthub/go-mysql-server/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxGoMySQLServerModule(t *testing.T) {
	var server *sqle.Server

	serverPort := transport.FindUnusedTestTCPPort(t)

	t.Run("test with tcp server", func(t *testing.T) {
		t.Setenv("APP_CONFIG_PATH", "testdata/config")
		t.Setenv("SERVER_TRANSPORT", "tcp")
		t.Setenv("SERVER_PORT", fmt.Sprintf("%d", serverPort))

		app := fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxtrace.FxTraceModule,
			fxgomysqlserver.FxGoMySQLServerModule,
			fx.Populate(&server),
		)

		app.RequireStart()
		assert.NoError(t, app.Err())

		db, err := sql.Open("mysql", fmt.Sprintf("user:password@tcp(localhost:%d)/db", serverPort))
		assert.NoError(t, err)

		err = db.Ping()
		assert.NoError(t, err)

		_, err = db.Exec("CREATE TABLE tcp (id int)")
		assert.NoError(t, err)

		row := db.QueryRow("SHOW TABLES")
		assert.NoError(t, row.Err())

		var tableName string
		err = row.Scan(&tableName)
		assert.NoError(t, err)

		if tableName != "tcp" {
			t.Errorf("expected to find tcp, but got %s", tableName)
		}

		err = db.Close()
		assert.NoError(t, err)

		app.RequireStop()
		assert.NoError(t, app.Err())
	})

	t.Run("test with memory server", func(t *testing.T) {
		t.Setenv("APP_CONFIG_PATH", "testdata/config")
		t.Setenv("SERVER_TRANSPORT", "memory")

		app := fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxtrace.FxTraceModule,
			fxgomysqlserver.FxGoMySQLServerModule,
			fx.Populate(&server),
		)

		app.RequireStart()
		assert.NoError(t, app.Err())

		db, err := sql.Open("mysql", "user:password@memory(bufconn)/db")
		assert.NoError(t, err)

		err = db.Ping()
		assert.NoError(t, err)

		_, err = db.Exec("CREATE TABLE memory (id int)")
		assert.NoError(t, err)

		row := db.QueryRow("SHOW TABLES")
		assert.NoError(t, row.Err())

		var tableName string
		err = row.Scan(&tableName)
		assert.NoError(t, err)

		if tableName != "memory" {
			t.Errorf("expected to find memory, but got %s", tableName)
		}

		err = db.Close()
		assert.NoError(t, err)

		app.RequireStop()
		assert.NoError(t, app.Err())
	})

	t.Run("test with observability", func(t *testing.T) {
		t.Setenv("APP_CONFIG_PATH", "testdata/config")
		t.Setenv("SERVER_TRANSPORT", "memory")
		t.Setenv("SERVER_TRACE_ENABLED", "true")

		var traceExporter tracetest.TestTraceExporter

		app := fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxtrace.FxTraceModule,
			fxgomysqlserver.FxGoMySQLServerModule,
			fx.Populate(&server, &traceExporter),
		)

		app.RequireStart()
		assert.NoError(t, app.Err())

		db, err := sql.Open("mysql", "user:password@memory(bufconn)/db")
		assert.NoError(t, err)

		err = db.Ping()
		assert.NoError(t, err)

		_, err = db.Exec("CREATE TABLE memory (id int)")
		assert.NoError(t, err)

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"parse",
			attribute.String("query", "CREATE TABLE memory (id int)"),
		)

		row := db.QueryRow("SHOW TABLES")
		assert.NoError(t, row.Err())

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"parse",
			attribute.String("query", "SHOW TABLES"),
		)

		var tableName string
		err = row.Scan(&tableName)
		assert.NoError(t, err)

		if tableName != "memory" {
			t.Errorf("expected to find memory, but got %s", tableName)
		}

		err = db.Close()
		assert.NoError(t, err)

		app.RequireStop()
		assert.NoError(t, app.Err())
	})

	t.Run("test with invalid transport", func(t *testing.T) {
		t.Setenv("APP_CONFIG_PATH", "testdata/config")
		t.Setenv("SERVER_TRANSPORT", "invalid")

		err := fx.New(
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxtrace.FxTraceModule,
			fxgomysqlserver.FxGoMySQLServerModule,
			fx.Populate(&server),
		).Start(context.Background())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown transport: invalid")
	})
}
