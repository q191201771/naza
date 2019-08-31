//package connection 对 net.Conn 接口的二次封装
package connection

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var connectionErr = errors.New("connection: fxxk")

type Connection interface {
	// 包含 interface net.Conn 的所有方法
	net.Conn

	// 额外提供的读方法
	ReadAtLeast(buf []byte, min int) (n int, err error)
	ReadLine() (line []byte, isPrefix bool, err error)

	// 额外提供的写方法
	Printf(fmt string, v ...interface{}) (n int, err error)
}

type Config struct {
	 // 如果不为0，则使用 buffer
	ReadBufSize  int
	WriteBufSize int

	// 如果不为0，则之后每次读/写都带超时
	ReadTimeoutMS int
	WriteTimeoutMS int
}

func New(conn net.Conn, config Config) Connection {
	var c connection
	c.Conn = conn
	if config.ReadBufSize > 0 {
		c.r = bufio.NewReaderSize(conn, config.ReadBufSize)
	} else {
		c.r = conn
	}
	if config.WriteBufSize > 0 {
		c.w = bufio.NewWriterSize(conn, config.WriteBufSize)
	} else {
		c.w = conn
	}
	c.config = config
	return &c
}

type connection struct {
	Conn   net.Conn
	r      io.Reader
	w      io.Writer
	config Config
}

// 为保证运行时性能，调用方需保证：
// 1. 调用 Config 方法时，不并行调用其它接口造成竞态读写属性
// 2. 只能在使用无缓冲时切换成缓冲，如果已经有缓冲，则不能再次切换
func (c *connection) Config(config Config) {
	c.config = config
}

func (c *connection) ReadAtLeast(buf []byte, min int) (n int, err error) {
	if c.config.ReadTimeoutMS > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
	}
	return io.ReadAtLeast(c.r, buf, min)
}

func (c *connection) ReadLine() (line []byte, isPrefix bool, err error) {
	bufioReader, ok := c.r.(*bufio.Reader)
	if !ok {
		// 目前只有使用了 bufio.Reader 时才能执行 ReadLine 操作
		panic(connectionErr)
	}
	if c.config.ReadTimeoutMS > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
	}
	return bufioReader.ReadLine()
}

func (c *connection) Printf(format string, v ...interface{}) (n int, err error) {
	if c.config.WriteTimeoutMS > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
	}
	return fmt.Fprintf(c.Conn, format, v...)
}

func (c *connection) Read(b []byte) (n int, err error) {
	if c.config.ReadTimeoutMS > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
	}
	return c.r.Read(b)
}

func (c *connection) Write(b []byte) (n int, err error) {
	if c.config.WriteTimeoutMS > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
	}
	return c.w.Write(b)
}

func (c *connection) Close() error {
	return c.Conn.Close()
}

func (c *connection) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *connection) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

func (c *connection) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}
func (c *connection) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}
