package config

import "strings"

// Transport is an enum for the supported database connection protocols.
type Transport string

const (
	UnknownTransport Transport = "unknown"
	TCPTransport     Transport = "tcp"
	MemoryTransport  Transport = "memory"
)

// String returns a string representation of the [Transport].
func (t Transport) String() string {
	return string(t)
}

// FetchTransport returns a [Transport] for a given name.
func FetchTransport(name string) Transport {
	//nolint:exhaustive
	switch t := Transport(strings.ToLower(name)); t {
	case TCPTransport,
		MemoryTransport:
		return t
	default:
		return UnknownTransport
	}
}
