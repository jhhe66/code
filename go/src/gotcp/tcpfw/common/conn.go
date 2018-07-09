// Copyright 2016 songliwei
//
// HelloTalk.inc

package common

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gansidui/gotcp"
)

// Error type
var (
	errConnClosing      = errors.New("use of closed network connection")
	ErrIllegalIpOrPort  = errors.New("ip or port illegal")
	errReachMaxTryCount = errors.New("reach max retry count")
)

// conn is the low-level implementation of Conn
type conn struct {

	// Shared
	mu       sync.Mutex
	pending  int
	err      error
	conn     *net.TCPConn
	protocol gotcp.Protocol // customize packet protocol

	// Read
	readTimeout time.Duration
	//	br          *bufio.Reader

	// Write
	writeTimeout time.Duration
	bw           *bufio.Writer
}

// DialOption specifies an option for dialing a Redis server.
type DialOption struct {
	f func(*dialOptions)
}

type dialOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	dial         func(ip, port string) (*net.TCPConn, error)
}

// DialReadTimeout specifies the timeout for reading a single command reply.
func DialReadTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.readTimeout = d
	}}
}

// DialWriteTimeout specifies the timeout for writing a single command.
func DialWriteTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.writeTimeout = d
	}}
}

func dialTCPConnect(ip, port string) (conn *net.TCPConn, err error) {
	if ip == "" || port == "" {
		err = ErrIllegalIpOrPort
		return nil, err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+port)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	err = conn.SetKeepAlive(true)
	if err != nil {
		return conn, err
	}

	err = conn.SetKeepAlivePeriod(time.Minute)
	if err != nil {
		return conn, err
	}
	return conn, err
}

// Dial connects to the api server at the given ip and
// port using the specified options.
func Dial(ip, port string, proto gotcp.Protocol, options ...DialOption) (Conn, error) {
	do := dialOptions{
		dial: dialTCPConnect,
	}
	for _, option := range options {
		option.f(&do)
	}

	netConn, err := do.dial(ip, port)
	if err != nil {
		return nil, err
	}

	c := &conn{
		conn: netConn,
		bw:   bufio.NewWriter(netConn),
		//br:           bufio.NewReader(netConn),
		protocol:     proto,
		readTimeout:  do.readTimeout,
		writeTimeout: do.writeTimeout,
	}

	return c, nil
}

func (c *conn) Close() error {
	c.mu.Lock()
	err := c.err
	if c.err == nil {
		c.err = errors.New("conn: closed")
		err = c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) fatal(err error) error {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		// Close connection to force errors on subsequent calls and to unblock
		// other reader or writer.
		c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *conn) writeBytes(p []byte) error {
	_, err := c.bw.Write(p)
	return err
}

type protocolError string

func (pe protocolError) Error() string {
	return fmt.Sprintf("tcp_api: %s (possible server error or unsupported concurrent read by application)", string(pe))
}

func (c *conn) Send(req gotcp.Packet) error {
	c.mu.Lock()
	c.pending += 1
	c.mu.Unlock()
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	buf := req.Serialize()
	if err := c.writeBytes(buf); err != nil {
		return c.fatal(err)
	}
	return nil
}

func (c *conn) Flush() error {
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if err := c.bw.Flush(); err != nil {
		return c.fatal(err)
	}
	return nil
}

func (c *conn) Receive() (p gotcp.Packet, err error) {
	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	// The pending field is decremented after the reply is read to handle the
	// case where Receive is called before Send.
	c.mu.Lock()
	if c.pending > 0 {
		c.pending -= 1
	}
	c.mu.Unlock()
	p, err = c.protocol.ReadPacket(c.conn)
	if err != nil { // 读取报文失败 关闭连接
		return nil, c.fatal(err)
	}

	return p, nil
}

func (c *conn) Do(req gotcp.Packet) (rsp gotcp.Packet, err error) {
	c.mu.Lock()
	//pending := c.pending
	c.pending = 0
	c.mu.Unlock()

	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	buf := req.Serialize()
	if err := c.writeBytes(buf); err != nil {
		return nil, c.fatal(err)
	}

	if err := c.bw.Flush(); err != nil {
		return nil, c.fatal(err)
	}

	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}

	if rsp, err = c.protocol.ReadPacket(c.conn); err != nil {
		return nil, c.fatal(err)
	}

	return rsp, err
}
