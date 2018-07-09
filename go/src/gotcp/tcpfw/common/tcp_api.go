package common

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gansidui/gotcp"
)

// Const
const (
	CMaxTryCount         = 3   // 最多尝试3次重试
	CRetSendFailed       = 501 // 发送报文失败
	CRetReadFailed       = 502 // 读取报文失败
	CRetUnMarshallFailed = 503 // UnMarashl失败
	CRetMarshallFailed   = 504 // Marashl失败
)

// 加密类型
const (
	CNoneKey    = 0 // 不加密
	CSessionKey = 1 // Session Key 加密
	CRandKey    = 2 // Rand    Key 加密
	CServKey    = 3 // Server  Key 加密
)

// 协议类型
const (
	CClineReq   = 0xF0 // 客户端请求
	CServAck    = 0xF1 // 服务器应答
	CServSend   = 0xF2 // 服务器主动发包
	CClientAck  = 0xF3 // 客户端应答
	CServToServ = 0xF4 // 服务器之间的包
)

const (
	CVerMmedia = 4 //版本号
)

// Error type
var (
	ErrConnClosing      = errors.New("use of closed network connection")
	ErrIpOrPortIllegal  = errors.New("ip or port illegal")
	ErrReachMaxTryCount = errors.New("reach max retry count")
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type TCPApi struct {
	DestIp       string
	DestPort     int
	ReadTimeOut  int            // 发送超时时间 单位秒
	WriteTimeOut int            // 接收超时时间 单位秒
	conn         *net.TCPConn   // the raw connection
	protocol     gotcp.Protocol // customize packet protocol
	extraData    interface{}    // to save extra data
	closeOnce    sync.Once      // close the conn, once, per instance
	closeFlag    int32          // close flag
	l            sync.Mutex     // sync mutex goroutines use only mutex
}

// newConn returns a wrapper of raw conn
func NewTCPApi(destIp string, destPort int, protocol gotcp.Protocol) *TCPApi {
	return &TCPApi{
		DestIp:       destIp,
		DestPort:     destPort,
		ReadTimeOut:  1,
		WriteTimeOut: 1,
		protocol:     protocol,
	}
}

// GetExtraData gets the extra data from the Conn
func (c *TCPApi) GetExtraData() interface{} {
	return c.extraData
}

// PutExtraData puts the extra data with the Conn
func (c *TCPApi) PutExtraData(data interface{}) {
	c.l.Lock()
	defer c.l.Unlock()

	c.extraData = data
}

// GetRawConn returns the raw net.TCPConn from the Conn
func (c *TCPApi) GetRawConn() *net.TCPConn {
	return c.conn
}

// Close closes the connection
func (c *TCPApi) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		c.conn.Close()
	})
}

// IsClosed indicates whether or not the connection is closed
func (c *TCPApi) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

// checker with out lock
func (c *TCPApi) isInit() bool {
	return c.conn != nil
}

// create conn with out lock
func (c *TCPApi) createConn() (err error) {
	if c.DestIp == "" || c.DestPort == 0 {
		err = ErrIpOrPortIllegal
		return err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", c.DestIp+":"+strconv.Itoa(c.DestPort))
	if err != nil {
		return err
	}

	c.conn, err = net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return err
	}
	err = c.conn.SetKeepAlive(true)
	if err != nil {
		return err
	}

	err = c.conn.SetKeepAlivePeriod(time.Minute)
	if err != nil {
		return err
	}

	return err
}

func (c *TCPApi) ReadPacket() (p gotcp.Packet, err error) {
	c.l.Lock()
	defer c.l.Unlock()
	if c.IsClosed() {
		err = ErrConnClosing
		return nil, err
	}

	if !c.isInit() {
		err = c.createConn()
		if err != nil { // create connect failed
			return
		}
	}
	// 设置接收等待的最长时间 ReadTimeOut s
	err = c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.ReadTimeOut) * time.Second))
	if err != nil {
		return nil, err
	}
	p, err = c.protocol.ReadPacket(c.conn)
	if err != nil { // 读取报文失败 关闭连接
		c.conn.Close()
		c.createConn()
		return nil, err
	}
	return
}

func (c *TCPApi) WritePacket(p gotcp.Packet) (err error) {
	c.l.Lock()
	defer c.l.Unlock()

	if c.IsClosed() {
		err = ErrConnClosing
		return err
	}

	if !c.isInit() {
		err = c.createConn()
		if err != nil { // create connect failed
			return err
		}
	}

	buf := p.Serialize()
	bufLen := len(buf)
	total := 0
	tryCount := 0
	for (tryCount < int(CMaxTryCount)) && (total < bufLen) {
		// 设置发送等待的最长时间 WriteTimeOut s
		err = c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.WriteTimeOut) * time.Second))
		if err != nil {
			return err
		}
		n, err := c.conn.Write(buf[total:])
		if err != nil { // 写入出错(包含写入超时) 新建一条连接从头开始发送报文
			c.conn.Close()
			err = c.createConn()
			if err != nil { // create connect failed
				return err
			}
			// 新建连接成功
			total = 0 // 从头开始发送
			tryCount += 1
			continue
		}
		total += n // total增加写入conn 的量
	}
	if tryCount == int(CMaxTryCount) { // 超过最大尝试次数
		err = ErrReachMaxTryCount
		return err
	} else { // 成功写入 返回nil
		err = nil
		return err
	}
}
