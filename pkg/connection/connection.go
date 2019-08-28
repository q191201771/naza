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

	// 带超时的读/写方法
	ReadWithTimeout(b []byte, timeoutMS int) (n int, err error)
	ReadAtLeastWithTimeout(buf []byte, min int, timeoutMS int) (n int, err error)
	ReadLineWithTimeout(timeoutMS int) (line []byte, isPrefix bool, err error)
	WriteWithTimeout(b []byte, timeoutMS int) (n int, err error)
	PrintfWithTimeout(timeoutMS int, fmt string, v ...interface{}) (n int, err error)
}

// 可配置是直接从 net.Conn 上读/写数据，还是中间添加一层buffer缓冲
type Config struct {
	ReadBufSize  int
	WriteBufSize int
}

func New(conn net.Conn, config *Config) Connection {
	var c connection
	c.Conn = conn
	if config != nil {
		if config.ReadBufSize > 0 {
			c.r = bufio.NewReaderSize(conn, config.ReadBufSize)
		}
		if config.WriteBufSize > 0 {
			c.w = bufio.NewWriterSize(conn, config.WriteBufSize)
		}
		c.config = *config
	}
	if c.r == nil {
		c.r = conn
	}
	if c.w == nil {
		c.w = conn
	}

	return &c
}

type connection struct {
	Conn   net.Conn
	r      io.Reader
	w      io.Writer
	config Config
}

func (c *connection) ReadAtLeastWithTimeout(buf []byte, min int, timeoutMS int) (n int, err error) {
	if timeoutMS > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond))
	}
	return io.ReadAtLeast(c.r, buf, min)
}

func (c *connection) ReadLineWithTimeout(timeoutMS int) (line []byte, isPrefix bool, err error) {
	if timeoutMS > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond))
	}
	return c.ReadLine()
}

func (c *connection) ReadWithTimeout(b []byte, timeoutMS int) (n int, err error) {
	if timeoutMS > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond))
	}
	return c.Conn.Read(b)
}

func (c *connection) WriteWithTimeout(b []byte, timeoutMS int) (n int, err error) {
	if timeoutMS > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond))
	}
	return c.Conn.Write(b)
}

func (c *connection) PrintfWithTimeout(timeoutMS int, format string, v ...interface{}) (n int, err error) {
	if timeoutMS > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond))
	}
	return fmt.Fprintf(c.Conn, format, v...)
}

func (c *connection) ReadAtLeast(buf []byte, min int) (n int, err error) {
	return io.ReadAtLeast(c.r, buf, min)
}

func (c *connection) ReadLine() (line []byte, isPrefix bool, err error) {
	bufioReader, ok := c.r.(*bufio.Reader)
	if !ok {
		// 目前只有使用了 bufio.Reader 时才能执行 ReadLine 操作
		return nil, false, connectionErr
	}
	return bufioReader.ReadLine()
}

func (c *connection) Printf(format string, v ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Conn, format, v...)
}

func (c *connection) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *connection) Write(b []byte) (n int, err error) {
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
