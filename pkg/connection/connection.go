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
	"github.com/q191201771/nezha/pkg/log"
	"io"
	"net"
	"sync"
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

	// 阻塞直到连接主动或被动关闭
	// @return 返回 nil 则是本端主动调用 Close 关闭
	Done() <-chan error

	// TODO chef: 这几个接口是否不提供
	ModWriteChanSize(n int)
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
	wMsgTClose // TODO chef: 没有使用
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
	c.doneChan = make(chan error, 1)
	c.exitChan = make(chan struct{}, 1)
	c.config = config
	return &c
}

type connection struct {
	Conn          net.Conn
	r             io.Reader
	w             io.Writer
	config        Config
	wChan         chan wmsg
	flushDoneChan chan struct{}
	doneChan      chan error
	exitChan      chan struct{}
	closeOnce     sync.Once
}

// Mod类型函数不加锁

// 由调用方保证不和写操作并发执行
func (c *connection) ModWriteChanSize(n int) {
	if c.config.WChanSize > 0 {
		panic(connectionErr)
	}
	c.config.WChanSize = n
	c.wChan = make(chan wmsg, n)
	c.flushDoneChan = make(chan struct{}, 1)
	go c.runWriteLoop()
}

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
		err = c.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
		if err != nil {
			log.Debugf("error=%v", err)
			return 0, err
		}
	}
	n, err = io.ReadAtLeast(c.r, buf, min)
	if err != nil {
		log.Debugf("error=%v", err)
		c.close(err)
	}
	return n, err
}

func (c *connection) ReadLine() (line []byte, isPrefix bool, err error) {
	bufioReader, ok := c.r.(*bufio.Reader)
	if !ok {
		// 目前只有使用了 bufio.Reader 时才能执行 ReadLine 操作
		panic(connectionErr)
	}
	if c.config.ReadTimeoutMS > 0 {
		err = c.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
		if err != nil {
			log.Debugf("error=%v", err)
			return nil, false, err
		}
	}
	line, isPrefix, err = bufioReader.ReadLine()
	if err != nil {
		log.Debugf("error=%v", err)
		c.close(err)
	}
	return line, isPrefix, err
}

func (c *connection) Printf(format string, v ...interface{}) (n int, err error) {
	if c.config.WriteTimeoutMS > 0 {
		_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
	}
	return fmt.Fprintf(c.Conn, format, v...)
}

func (c *connection) Read(b []byte) (n int, err error) {
	if c.config.ReadTimeoutMS > 0 {
		err = c.SetReadDeadline(time.Now().Add(time.Duration(c.config.ReadTimeoutMS) * time.Millisecond))
		if err != nil {
			log.Debugf("error=%v", err)
			return 0, err
		}
	}
	n, err = c.r.Read(b)
	if err != nil {
		log.Debugf("error=%v", err)
		c.close(err)
	}
	return n, err
}

func (c *connection) Write(b []byte) (n int, err error) {
	if c.config.WChanSize > 0 {
		c.wChan <- wmsg{t: wMsgTWrite, b: b}
		return len(b), nil
	}
	return c.write(b)
}

func (c *connection) write(b []byte) (n int, err error) {
	if c.config.WriteTimeoutMS > 0 {
		err = c.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
		if err != nil {
			log.Debugf("error=%v", err)
			return 0, err
		}
	}
	n, err = c.w.Write(b)
	if err != nil {
		log.Debugf("error=%v", err)
		c.close(err)
	}
	return n, err
}

func (c *connection) runWriteLoop() {
	for {
		select {
		case <-c.exitChan:
			log.Debug("exitChan recv, exit write loop.")
			return
		case msg := <-c.wChan:
			switch msg.t {
			case wMsgTWrite:
				if _, err := c.write(msg.b); err != nil {
					log.Debugf("error=%v", err)
					return
				}
			case wMsgTFlush:
				if err := c.flush(); err != nil {
					log.Debugf("error=%v", err)
					c.flushDoneChan <- struct{}{}
					return
				}
				c.flushDoneChan <- struct{}{}
			case wMsgTClose:
				// TODO chef: 是否需要
			}
		}
	}
}

func (c *connection) Flush() error {
	if c.config.WChanSize > 0 {
		c.wChan <- wmsg{t: wMsgTFlush}
		<-c.flushDoneChan
		return nil
	}

	return c.flush()
}

func (c *connection) flush() error {
	w, ok := c.w.(*bufio.Writer)
	if ok {
		if c.config.WriteTimeoutMS > 0 {
			err := c.SetWriteDeadline(time.Now().Add(time.Duration(c.config.WriteTimeoutMS) * time.Millisecond))
			if err != nil {
				log.Debugf("error=%v", err)
				return err
			}
		}
		if err := w.Flush(); err != nil {
			log.Debugf("error=%v", err)
			c.close(err)
			return err
		}
	}
	return nil
}

func (c *connection) Close() error {
	log.Debugf("Close.")
	c.close(nil)
	return nil
}

func (c *connection) close(err error) {
	log.Debugf("close. err=%v", err)
	c.closeOnce.Do(func() {
		if c.config.WChanSize > 0 {
			c.exitChan <- struct{}{}
		}
		c.doneChan <- err
		_ = c.Conn.Close()
	})
}

func (c *connection) Done() <-chan error {
	return c.doneChan
	//err := <-c.doneChan
	//log.Debugf("Done. err=%v", err)
	//if err != nil {
	//	c.close(err)
	//}
	//
	//ch := make(chan error, 1)
	//ch <- err
	//return ch
}

func (c *connection) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *connection) SetDeadline(t time.Time) error {
	err := c.Conn.SetDeadline(t)
	if err != nil {
		log.Debugf("error=%v", err)
		c.close(err)
	}
	return err
}

func (c *connection) SetReadDeadline(t time.Time) error {
	err := c.Conn.SetReadDeadline(t)
	if err != nil {
		log.Debugf("error=%v", err)
		c.close(err)
	}
	return err
}

func (c *connection) SetWriteDeadline(t time.Time) error {
	err := c.Conn.SetWriteDeadline(t)
	if err != nil {
		log.Debugf("error=%v", err)
		c.close(err)
	}
	return err
}
