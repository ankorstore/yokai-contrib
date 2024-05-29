package config

import "strings"

// Transport is an enum for the supported database connection protocols.
type Transport string

const (
	UnknownTransport Transport = "unknown"
	TCPTransport     Transport = "tcp"
	SocketTransport  Transport = "socket"
	MemoryTransport  Transport = "memory"
)

// String returns a string representation of the [Transport].
func (d Transport) String() string {
	return string(d)
}

// FetchTransport returns a [Transport] for a given name.
func FetchTransport(name string) Transport {
	//nolint:exhaustive
	switch d := Transport(strings.ToLower(name)); d {
	case TCPTransport,
		SocketTransport,
		MemoryTransport:
		return d
	default:
		return UnknownTransport
	}
}
