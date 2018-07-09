// Copyright 2016 songliwei
//
// HelloTalk.inc

package common

import "github.com/gansidui/gotcp"

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string { return string(err) }

// Conn represents a connection to a server.
type Conn interface {
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value if the connection is broken. The returned
	// value is either the first non-nil value returned from the underlying
	// network connection or a protocol parsing error. Applications should
	// close broken connections.
	Err() error

	// Do sends a command to the server and returns the received reply.
	Do(req gotcp.Packet) (rsp gotcp.Packet, err error)
	// Send writes the command to the client's output buffer.
	Send(p gotcp.Packet) error
	// Flush flushes the output buffer to the server.
	Flush() error

	// Receive receives a single reply from the server
	Receive() (p gotcp.Packet, err error)
}
