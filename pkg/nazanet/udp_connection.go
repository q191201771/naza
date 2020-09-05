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

const maxReadSizeOfUDPConnection = 1500

// @return 上层回调返回false，则关闭UDPConnection
//
type OnReadUDPPacket func(b []byte, raddr *net.UDPAddr, err error) bool

type UDPConnection struct {
	conn   *net.UDPConn
	ruaddr *net.UDPAddr
}

// @param laddr: 本地bind地址，如果设置为空，则自动选择可用端口
//               比如作为客户端时，如果不想特别指定本地端口，可以设置为空
//
// @param raddr: 如果为空，则只能使用func Write2Addr携带对端地址进行发送，不能使用func Write
//               好处是作为客户端时，对端地址通常只有一个，在构造函数中指定，后续就不用每次发送都指定
//
func NewUDPConnection(laddr string, raddr string) (c *UDPConnection, err error) {
	c = &UDPConnection{}
	conn, err := listenUDPWithAddr(laddr)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	if c.ruaddr, err = net.ResolveUDPAddr(udpNetwork, raddr); err != nil {
		return nil, err
	}

	return c, nil
}

// @param conn: 直接传入外部创建好的连接对象供内部使用
func NewUDPConnectionWithConn(conn *net.UDPConn) (c *UDPConnection) {
	return &UDPConnection{
		conn: conn,
	}
}

// 阻塞直至Read发生错误或上层回调函数返回false
//
// @return error: 如果外部调用Dispose，会返回error
//
func (c *UDPConnection) RunLoop(onRead OnReadUDPPacket) error {
	// TODO chef: 外部可以指定，是否复用
	b := make([]byte, maxReadSizeOfUDPConnection)
	for {
		n, raddr, err := c.conn.ReadFromUDP(b)
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
		if err := c.conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond)); err != nil {
			return nil, nil, err
		}
	}
	b := make([]byte, maxReadSizeOfUDPConnection)
	n, raddr, err := c.conn.ReadFromUDP(b)
	if err != nil {
		return nil, nil, err
	}
	return b[:n], raddr, nil
}

func (c *UDPConnection) Write(b []byte) error {
	_, err := c.conn.WriteToUDP(b, c.ruaddr)
	return err
}

func (c *UDPConnection) Write2Addr(b []byte, ruaddr *net.UDPAddr) error {
	_, err := c.conn.WriteToUDP(b, ruaddr)
	return err
}

func (c *UDPConnection) Dispose() error {
	return c.conn.Close()
}
