package config

import (
	"google.golang.org/grpc/test/bufconn"
)

// ConfigOptions are options for the [GoMySQLServerFactory] implementations.
type ConfigOptions struct {
	Transport Transport
	Listener  *bufconn.Listener
	User      string
	Password  string
	Socket    string
	Host      string
	Port      int
	Database  string
}

// DefaultGoMySQLServerConfigOptions are the default options used in the [DefaultGoMySQLServerFactory].
func DefaultGoMySQLServerConfigOptions() ConfigOptions {
	return ConfigOptions{
		Transport: TCPTransport,
		Listener:  bufconn.Listen(DefaultMemoryBufferSize),
		User:      DefaultUser,
		Password:  DefaultPassword,
		Socket:    DefaultSocket,
		Host:      DefaultHost,
		Port:      DefaultPort,
		Database:  DefaultDatabase,
	}
}

// GoMySQLServerConfigOption are functional options for the [GoMySQLServerFactory] implementations.
type GoMySQLServerConfigOption func(o *ConfigOptions)

// WithTransport is used to specify the [Transport].
func WithTransport(transport Transport) GoMySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.Transport = transport
	}
}

// WithUser is used to specify the user.
func WithUser(user string) GoMySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.User = user
	}
}

// WithPassword is used to specify the password.
func WithPassword(password string) GoMySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.Password = password
	}
}

// WithSocket is used to specify the socket.
func WithSocket(socket string) GoMySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.Socket = socket
	}
}

// WithHost is used to specify the host.
func WithHost(host string) GoMySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.Host = host
	}
}

// WithPort is used to specify the port.
func WithPort(port int) GoMySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.Port = port
	}
}

// WithDatabase is used to specify the database.
func WithDatabase(database string) GoMySQLServerConfigOption {
	return func(o *ConfigOptions) {
		o.Database = database
	}
}
