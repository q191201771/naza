// package connection
//
// 对 net.Conn 接口的二次封装，目的有两个：
// 1. 在流媒体传输这种特定的长连接场景下提供更方便、高性能的接口
// 2. 便于后续将 TCPConn 替换成其他传输协议
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
	// Read
	// Write
	// Close
	// LocalAddr
	// RemoteAddr
	// SetDeadline
	// SetReadDeadline
	// SetWriteDeadline
	net.Conn

	ReadAtLeast(buf []byte, min int) (n int, err error)
	ReadLine() (line []byte, isPrefix bool, err error)

	// TODO chef: 这个接口是否不提供
	Printf(fmt string, v ...interface{}) (n int, err error)

	// 如果使用了 bufio 写缓冲，则将缓冲中的数据发送出去
	// 如果使用了 channel 异步发送，则阻塞等待，直到之前 channel 中的数据全部发送完毕
	Flush() error

	// TODO chef: 这几个接口是否不提供
	ModWriteBufSize(n int)
	ModReadTimeoutMS(n int)
	ModWriteTimeoutMS(n int)
}

type Config struct {
	// 如果不为0，则之后每次读/写使用 bufio 的缓冲
	ReadBufSize  int
	WriteBufSize int

	// 如果不为0，则之后每次读/写都带超时
	ReadTimeoutMS  int
	WriteTimeoutMS int

	// 如果不过0，则写使用 channel 将数据发送到后台协程中发送
	WChanSize int
}

type wMsgT int

const (
	_ wMsgT = iota
	wMsgTWrite
	wMsgTFlush
	wMsgTClose
)

type wmsg struct {
	t wMsgT
	b []byte
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
	if config.WChanSize > 0 {
		c.wChan = make(chan wmsg, config.WChanSize)
		c.flushDoneChan = make(chan struct{}, 1)
		go c.runWriteLoop()
	}
	c.config = config
	return &c
}

type connection struct {
	Conn   net.Conn
	r      io.Reader
	w      io.Writer
	wChan  chan wmsg
	flushDoneChan chan struct{}
	config Config
}

// Mod类型函数不加锁

// 由调用方保证不和写操作并发执行
func (c *connection) ModWriteBufSize(n int) {
	if c.config.WriteBufSize > 0 {
		// 如果之前已经设置过写缓冲，直接 panic
		// 这里改成 flush 后替换成新缓冲也行，暂时没这个必要
		panic(connectionErr)
	}
	c.config.WriteBufSize = n
	c.w = bufio.NewWriterSize(c.Conn, n)
}

func (c *connection) ModReadTimeoutMS(n int) {
	if c.config.ReadTimeoutMS > 0 {
		panic(connectionErr)
	}
	c.config.ReadTimeoutMS = n
}

func (c *connection) ModWriteTimeoutMS(n int) {
	if c.config.WriteTimeoutMS > 0 {
		panic(connectionErr)
	}
	c.config.WriteTimeoutMS = n
}

func (c *connection) ReadAtLeast(buf []byte, min int) (n int, err error) {
	if c.config.ReadTimeoutMS > 0 {
		// TODO chef: 超时的错误返回
		_ = c.Conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
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
		_ = c.Conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
	}
	return bufioReader.ReadLine()
}

func (c *connection) Printf(format string, v ...interface{}) (n int, err error) {
	if c.config.WriteTimeoutMS > 0 {
		_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
	}
	return fmt.Fprintf(c.Conn, format, v...)
}

func (c *connection) Read(b []byte) (n int, err error) {
	if c.config.ReadTimeoutMS > 0 {
		_ = c.Conn.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
	}
	return c.r.Read(b)
}

func (c *connection) Write(b []byte) (n int, err error) {
	if c.config.WChanSize > 0 {
		c.wChan <- wmsg{}
		return len(b), nil
	}
	return c.write(b)
}

func (c *connection) write(b []byte) (n int, err error) {
	if c.config.WriteTimeoutMS > 0 {
		_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
	}
	return c.w.Write(b)
}

func (c *connection) runWriteLoop() {
	for {
		msg, ok := <- c.wChan
		if !ok {
			return
		}
		switch msg.t {
		case wMsgTWrite:
			if _, err := c.write(msg.b); err != nil {
				_ = c.Close()
			}
		case wMsgTFlush:
			c.flush()
			c.flushDoneChan <- struct{}{}
		case wMsgTClose:
			// TODO chef: 是否需要
		}
	}
}

func (c *connection) Flush() error {
	if c.config.WChanSize > 0 {
		c.wChan <- wmsg{t:wMsgTFlush}
		<- c.flushDoneChan
		return nil
	}

	return c.flush()
}

func (c *connection) flush() error {
	w, ok := c.w.(*bufio.Writer)
	if ok {
		if c.config.WriteTimeoutMS > 0 {
			_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
		}
		if err := w.Flush(); err != nil {
			_ = c.Close()
		}
	}
	return nil
}

// 调用方需保证不和 Write 接口并发调用
func (c *connection) Close() error {
	if c.config.WChanSize > 0 {
		close(c.wChan)
	}

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
