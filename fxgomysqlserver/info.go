package fxgomysqlserver

import (
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

// FxGoMySQLServerModuleInfo is a module info collector for gomysqlserver.
type FxGoMySQLServerModuleInfo struct {
	config *config.GoMySQLServerConfig
	server *server.Server
}

// NewFxGoMySQLServerModuleInfo returns a new [FxGoMySQLServerModuleInfo].
func NewFxGoMySQLServerModuleInfo(config *config.GoMySQLServerConfig, server *server.Server) *FxGoMySQLServerModuleInfo {
	return &FxGoMySQLServerModuleInfo{
		config: config,
		server: server,
	}
}

// Name return the name of the module info.
func (i *FxGoMySQLServerModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *FxGoMySQLServerModuleInfo) Data() map[string]interface{} {
	data := make(map[string]interface{})

	dsn, err := i.config.DSN()
	if err == nil {
		data["dsn"] = dsn
	}

	sessionVars := make(map[uint32]interface{})
	err = i.server.SessionManager().Iter(func(session sql.Session) (stop bool, err error) {
		sessionVars[session.ID()] = session.GetAllSessionVariables()

		return false, nil
	})
	if err == nil {
		data["sessions"] = sessionVars
	}

	return data
}
