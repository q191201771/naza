// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// package connection
//
// 对 net.Conn 接口的二次封装，目的有两个：
// 1. 在流媒体传输这种特定的长连接场景下提供更方便、高性能的接口
// 2. 便于后续将 TCPConn 替换成其他传输协议
package connection

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/q191201771/naza/pkg/nazalog"
)

var connectionErr = errors.New("naza.connection: fxxk")

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
	ReadLine() (line []byte, isPrefix bool, err error) // 只有设置了 ReadBufSize 才可以使用这个方法

	// 如果使用了 bufio 写缓冲，则将缓冲中的数据发送出去
	// 如果使用了 channel 异步发送，则阻塞等待，直到之前 channel 中的数据全部发送完毕
	// 一般在 Close 前，想要将剩余数据发送完毕时调用
	Flush() error

	// 阻塞直到连接关闭或发生错误
	// @return 返回 nil 则是本端主动调用 Close 关闭
	Done() <-chan error

	// TODO chef: 这几个接口是否不提供
	// Mod类型函数不加锁，需要调用方保证不发生竞态调用
	ModWriteChanSize(n int)
	ModWriteBufSize(n int)
	ModReadTimeoutMS(n int)
	ModWriteTimeoutMS(n int)
}

type Option struct {
	// 如果不为0，则之后每次读/写使用 bufio 的缓冲
	ReadBufSize  int
	WriteBufSize int

	// 如果不为0，则之后每次读/写都带超时
	ReadTimeoutMS  int
	WriteTimeoutMS int

	// 如果不为0，则写使用 channel 将数据发送到后台协程中发送
	WriteChanSize int
}

// 没有配置的属性，将按如下配置
var defaultOption = Option{
	ReadBufSize:    0,
	WriteBufSize:   0,
	ReadTimeoutMS:  0,
	WriteTimeoutMS: 0,
	WriteChanSize:  0,
}

type ModOption func(option *Option)

func New(conn net.Conn, modOptions ...ModOption) Connection {
	c := new(connection)
	c.doneChan = make(chan error, 1)
	c.Conn = conn

	c.option = defaultOption

	for _, fn := range modOptions {
		fn(&c.option)
	}

	if c.option.ReadBufSize > 0 {
		c.r = bufio.NewReaderSize(conn, c.option.ReadBufSize)
	} else {
		c.r = conn
	}

	if c.option.WriteBufSize > 0 {
		c.w = bufio.NewWriterSize(conn, c.option.WriteBufSize)
	} else {
		c.w = conn
	}

	if c.option.WriteBufSize > 0 {
		c.wChan = make(chan wMsg, c.option.WriteBufSize)
		c.flushDoneChan = make(chan struct{}, 1)
		c.exitChan = make(chan struct{}, 1)
		go c.runWriteLoop()
	}

	return c
}

type wMsgT int

const (
	_ wMsgT = iota
	wMsgTWrite
	wMsgTFlush
)

type wMsg struct {
	t wMsgT
	b []byte
}

type connection struct {
	Conn          net.Conn
	r             io.Reader
	w             io.Writer
	option        Option
	wChan         chan wMsg
	flushDoneChan chan struct{}
	exitChan      chan struct{}
	doneChan      chan error
	closeOnce     sync.Once
}

func (c *connection) ModWriteChanSize(n int) {
	if c.option.WriteChanSize > 0 {
		panic(connectionErr)
	}
	c.option.WriteChanSize = n
	c.wChan = make(chan wMsg, n)
	c.flushDoneChan = make(chan struct{}, 1)
	c.exitChan = make(chan struct{}, 1)
	go c.runWriteLoop()
}

func (c *connection) ModWriteBufSize(n int) {
	if c.option.WriteBufSize > 0 {
		// 如果之前已经设置过写缓冲，直接 panic
		// 这里改成 flush 后替换成新缓冲也行，暂时没这个必要
		panic(connectionErr)
	}
	c.option.WriteBufSize = n
	c.w = bufio.NewWriterSize(c.Conn, n)
}

func (c *connection) ModReadTimeoutMS(n int) {
	if c.option.ReadTimeoutMS > 0 {
		panic(connectionErr)
	}
	c.option.ReadTimeoutMS = n
}

func (c *connection) ModWriteTimeoutMS(n int) {
	if c.option.WriteTimeoutMS > 0 {
		panic(connectionErr)
	}
	c.option.WriteTimeoutMS = n
}

func (c *connection) ReadAtLeast(buf []byte, min int) (n int, err error) {
	if c.option.ReadTimeoutMS > 0 {
		err = c.SetReadDeadline(time.Now().Add(time.Duration(c.option.ReadTimeoutMS) * time.Millisecond))
		if err != nil {
			nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
			return 0, err
		}
	}
	n, err = io.ReadAtLeast(c.r, buf, min)
	if err != nil {
		nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
		c.close(err)
	}
	return n, err
}

// TODO chef: 测试 bufio 设置的大小 < 换行符位置时的情况
func (c *connection) ReadLine() (line []byte, isPrefix bool, err error) {
	bufioReader, ok := c.r.(*bufio.Reader)
	if !ok {
		// 目前只有使用了 bufio.Reader 时才能执行 ReadLine 操作
		panic(connectionErr)
	}
	if c.option.ReadTimeoutMS > 0 {
		err = c.SetReadDeadline(time.Now().Add(time.Duration(c.option.ReadTimeoutMS) * time.Millisecond))
		if err != nil {
			nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
			return nil, false, err
		}
	}
	line, isPrefix, err = bufioReader.ReadLine()
	if err != nil {
		nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
		c.close(err)
	}
	return line, isPrefix, err
}

func (c *connection) Read(b []byte) (n int, err error) {
	if c.option.ReadTimeoutMS > 0 {
		err = c.SetReadDeadline(time.Now().Add(time.Duration(c.option.ReadTimeoutMS) * time.Millisecond))
		if err != nil {
			nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
			return 0, err
		}
	}
	n, err = c.r.Read(b)
	if err != nil {
		nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
		c.close(err)
	}
	return n, err
}

func (c *connection) Write(b []byte) (n int, err error) {
	if c.option.WriteChanSize > 0 {
		c.wChan <- wMsg{t: wMsgTWrite, b: b}
		return len(b), nil
	}
	return c.write(b)
}

func (c *connection) write(b []byte) (n int, err error) {
	if c.option.WriteTimeoutMS > 0 {
		err = c.SetWriteDeadline(time.Now().Add(time.Duration(c.option.WriteTimeoutMS) * time.Millisecond))
		if err != nil {
			nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
			return 0, err
		}
	}
	n, err = c.w.Write(b)
	if err != nil {
		nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
		c.close(err)
	}
	return n, err
}

func (c *connection) runWriteLoop() {
	for {
		select {
		case <-c.exitChan:
			nazalog.Debug("exitChan recv, exit write loop.")
			return
		case msg := <-c.wChan:
			switch msg.t {
			case wMsgTWrite:
				if _, err := c.write(msg.b); err != nil {
					nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
					return
				}
			case wMsgTFlush:
				if err := c.flush(); err != nil {
					nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
					c.flushDoneChan <- struct{}{}
					return
				}
				c.flushDoneChan <- struct{}{}
			}
		}
	}
}

func (c *connection) Flush() error {
	if c.option.WriteChanSize > 0 {
		c.wChan <- wMsg{t: wMsgTFlush}
		<-c.flushDoneChan
		return nil
	}

	return c.flush()
}

func (c *connection) flush() error {
	w, ok := c.w.(*bufio.Writer)
	if ok {
		if c.option.WriteTimeoutMS > 0 {
			err := c.SetWriteDeadline(time.Now().Add(time.Duration(c.option.WriteTimeoutMS) * time.Millisecond))
			if err != nil {
				nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
				return err
			}
		}
		if err := w.Flush(); err != nil {
			nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
			c.close(err)
			return err
		}
	}
	return nil
}

func (c *connection) Close() error {
	nazalog.Debugf("naza connection Close. conn=%p", c)
	c.close(nil)
	return nil
}

func (c *connection) close(err error) {
	nazalog.Debugf("naza connection close. err=%v, conn=%p", err, c)
	c.closeOnce.Do(func() {
		if c.option.WriteChanSize > 0 {
			c.exitChan <- struct{}{}
		}
		c.doneChan <- err
		_ = c.Conn.Close()
	})
}

func (c *connection) Done() <-chan error {
	return c.doneChan
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
		nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
		c.close(err)
	}
	return err
}

func (c *connection) SetReadDeadline(t time.Time) error {
	err := c.Conn.SetReadDeadline(t)
	if err != nil {
		nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
		c.close(err)
	}
	return err
}

func (c *connection) SetWriteDeadline(t time.Time) error {
	err := c.Conn.SetWriteDeadline(t)
	if err != nil {
		nazalog.Debugf("naza connection. error=%v, conn=%p", err, c)
		c.close(err)
	}
	return err
}
