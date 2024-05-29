package config

import (
	"fmt"

	"google.golang.org/grpc/test/bufconn"
)

const (
	DefaultUser             = "user"
	DefaultPassword         = "password"
	DefaultSocket           = "/tmp/mysql.sock"
	DefaultHost             = "localhost"
	DefaultPort             = 3306
	DefaultDatabase         = "db"
	DefaultMemoryNetwork    = "memory"
	DefaultMemoryAddress    = "bufconn"
	DefaultMemoryBufferSize = 1024 * 1024
)

// GoMySQLServerConfig is the server configuration.
type GoMySQLServerConfig struct {
	transport Transport
	listener  *bufconn.Listener
	user      string
	password  string
	socket    string
	host      string
	port      int
	database  string
}

// NewGoMySQLServerConfig returns a new [GoMySQLServerConfig].
func NewGoMySQLServerConfig(options ...GoMySQLServerConfigOption) *GoMySQLServerConfig {
	// resolve options
	configOptions := DefaultGoMySQLServerConfigOptions()
	for _, opt := range options {
		opt(&configOptions)
	}

	// create config
	return &GoMySQLServerConfig{
		transport: configOptions.Transport,
		listener:  configOptions.Listener,
		user:      configOptions.User,
		password:  configOptions.Password,
		socket:    configOptions.Socket,
		host:      configOptions.Host,
		port:      configOptions.Port,
		database:  configOptions.Database,
	}
}

// Transport returns the configuration transport.
func (c *GoMySQLServerConfig) Transport() Transport {
	return c.transport
}

// Listener returns the configuration listener.
func (c *GoMySQLServerConfig) Listener() *bufconn.Listener {
	return c.listener
}

// User returns the configuration user.
func (c *GoMySQLServerConfig) User() string {
	return c.user
}

// Password returns the configuration password.
func (c *GoMySQLServerConfig) Password() string {
	return c.password
}

// Socket returns the configuration socket.
func (c *GoMySQLServerConfig) Socket() string {
	return c.socket
}

// Host returns the configuration host.
func (c *GoMySQLServerConfig) Host() string {
	return c.host
}

// Port returns the configuration port.
func (c *GoMySQLServerConfig) Port() int {
	return c.port
}

// Database returns the configuration database.
func (c *GoMySQLServerConfig) Database() string {
	return c.database
}

// DSN returns the dsn, depending on the configuration transport.
func (c *GoMySQLServerConfig) DSN() (string, error) {
	//nolint:exhaustive
	switch c.Transport() {
	case TCPTransport:
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			c.User(),
			c.Password(),
			c.Host(),
			c.Port(),
			c.Database(),
		), nil
	case SocketTransport:
		return fmt.Sprintf(
			"%s:%s@unix(%s)/%s",
			c.User(),
			c.Password(),
			c.Socket(),
			c.Database(),
		), nil
	case MemoryTransport:
		return fmt.Sprintf(
			"%s:%s@%s(%s)/%s",
			c.User(),
			c.Password(),
			DefaultMemoryNetwork,
			DefaultMemoryAddress,
			c.Database(),
		), nil
	default:
		return "", fmt.Errorf("unknown transport: %s", c.Transport().String())
	}
}
