package fxgomysqlserver_test

import (
	basesql "database/sql"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestFxGoMySQLServerModuleInfo(t *testing.T) {
	t.Parallel()

	t.Run("test info name", func(t *testing.T) {
		t.Parallel()

		info := fxgomysqlserver.NewFxGoMySQLServerModuleInfo(nil, nil)

		assert.Equal(t, fxgomysqlserver.ModuleName, info.Name())
	})

	t.Run("test info data", func(t *testing.T) {
		t.Parallel()

		serverConfig := config.NewGoMySQLServerConfig(
			config.WithTransport(config.MemoryTransport),
		)

		srv, err := server.NewDefaultGoMySQLServerFactory().Create(
			server.WithConfig(serverConfig),
		)
		assert.NoError(t, err)

		go func() {
			startErr := srv.Start()
			assert.NoError(t, startErr)
		}()

		dsn, err := serverConfig.DSN()
		assert.NoError(t, err)

		db, err := basesql.Open("mysql", dsn)
		assert.NoError(t, err)

		err = db.Ping()
		assert.NoError(t, err)

		info := fxgomysqlserver.NewFxGoMySQLServerModuleInfo(serverConfig, srv)

		assert.Equal(t, "user:password@memory(bufconn)/db", info.Data()["dsn"])

		sessions, ok := info.Data()["sessions"].(map[uint32]interface{})
		if !ok {
			t.Error("info sessions is not a map")
		}

		assert.Len(t, sessions, 1)

		err = db.Close()
		assert.NoError(t, err)

		err = srv.Close()
		assert.NoError(t, err)
	})
}
