// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazanet

import (
	"net"
	"time"
)

// @return 上层回调返回false，则关闭UDPConnection
//
type OnReadUDPPacket func(b []byte, raddr *net.UDPAddr, err error) bool

type UDPConnectionOption struct {
	// 两种初始化方式：
	// 方式一：直接传入外部创建好的连接对象供内部使用
	Conn *net.UDPConn
	// 方式二：填入地址，内部创建连接对象
	// LAddr: 本地bind地址，如果设置为空，则自动选择可用端口
	//        比如作为客户端时，如果不想特别指定本地端口，可以设置为空
	//
	// RAddr: 如果为空，则只能使用func Write2Addr携带对端地址进行发送，不能使用func Write
	//        不为空的作用：作为客户端时，对端地址通常只有一个，在构造函数中指定，后续就不用每次发送都指定
	//        注意，对端地址需显式填写IP
	//        注意，即使使用方式一，也可以设置Rddr
	//
	LAddr string
	RAddr string

	MaxReadPacketSize int  // 读取数据时，内存块大小
	AllocEachRead     bool // 使用Read Loop时，是否每次读取都申请新的内存块，如果为false，则复用一块内存块
}

var defaultOption = UDPConnectionOption{
	MaxReadPacketSize: 1500,
	AllocEachRead:     true,
}

type UDPConnection struct {
	option UDPConnectionOption
	ruaddr *net.UDPAddr
}

type ModUDPConnectionOption func(option *UDPConnectionOption)

func NewUDPConnection(modOptions ...ModUDPConnectionOption) (*UDPConnection, error) {
	var err error

	c := &UDPConnection{}
	c.option = defaultOption
	for _, fn := range modOptions {
		fn(&c.option)
	}
	if c.ruaddr, err = net.ResolveUDPAddr(udpNetwork, c.option.RAddr); err != nil {
		return nil, err
	}
	if c.option.Conn != nil {
		return c, nil
	}

	if c.option.Conn, err = listenUDPWithAddr(c.option.LAddr); err != nil {
		return nil, err
	}
	return c, err
}

// 阻塞直至Read发生错误或上层回调函数返回false
//
// @return error: 如果外部调用Dispose，会返回error
//
func (c *UDPConnection) RunLoop(onRead OnReadUDPPacket) error {
	var b []byte
	if !c.option.AllocEachRead {
		b = make([]byte, c.option.MaxReadPacketSize)
	}
	for {
		if c.option.AllocEachRead {
			b = make([]byte, c.option.MaxReadPacketSize)
		}
		n, raddr, err := c.option.Conn.ReadFromUDP(b)
		if keepRunning := onRead(b[:n], raddr, err); !keepRunning {
			if err == nil {
				return c.Dispose()
			}
		}
		if err != nil {
			return err
		}
	}
}

// 直接读取数据，不使用RunLoop
//
func (c *UDPConnection) ReadWithTimeout(timeoutMS int) ([]byte, *net.UDPAddr, error) {
	if timeoutMS > 0 {
		if err := c.option.Conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)); err != nil {
			return nil, nil, err
		}
	}
	b := make([]byte, c.option.MaxReadPacketSize)
	n, raddr, err := c.option.Conn.ReadFromUDP(b)
	if err != nil {
		return nil, nil, err
	}
	return b[:n], raddr, nil
}

func (c *UDPConnection) Write(b []byte) error {
	_, err := c.option.Conn.WriteToUDP(b, c.ruaddr)
	return err
}

func (c *UDPConnection) Write2Addr(b []byte, ruaddr *net.UDPAddr) error {
	_, err := c.option.Conn.WriteToUDP(b, ruaddr)
	return err
}

func (c *UDPConnection) Dispose() error {
	return c.option.Conn.Close()
}
