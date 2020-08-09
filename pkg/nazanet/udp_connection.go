// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazanet

import "net"

// TODO chef:
//   1. 长度可能需要提供接口供业务方设置
//   2. 每次读取数据时的buf，是新建还是复用，需要提供接口供业务方设置

const maxUDPPacketLength = 1500

// @param err: 注意，当err不为nil时，read loop将结束并退出（该语义后续可能发生变化，具体见代码）
type OnReadUDPPacket func(b []byte, remoteAddr net.Addr, err error)

type UDPConnection struct {
	localAddr       string
	conn            *net.UDPConn
	onReadUDPPacket OnReadUDPPacket
}

func NewUDPConnectionWithLocalAddr(localAddr string, onReadUDPPacket OnReadUDPPacket) *UDPConnection {
	return &UDPConnection{
		localAddr:       localAddr,
		onReadUDPPacket: onReadUDPPacket,
	}
}

// 直接使用已绑定好监听的net.UDPConn对象
func NewUDPConnectionWithConn(conn *net.UDPConn, onReadUDPPacket OnReadUDPPacket) *UDPConnection {
	return &UDPConnection{
		conn:            conn,
		onReadUDPPacket: onReadUDPPacket,
	}
}

// 配合func NewUDPConnectionWithLocalAddr使用
func (u *UDPConnection) Listen() error {
	udpAddr, err := net.ResolveUDPAddr("udp", u.localAddr)
	if err != nil {
		return err
	}
	u.conn, err = net.ListenUDP("udp", udpAddr)
	return err
}

// 开启读取事件循环，读取到数据时通过回调返回给上层
func (u *UDPConnection) RunLoop() {
	for {
		b := make([]byte, maxUDPPacketLength)
		length, remoteAddr, err := u.conn.ReadFrom(b)
		u.onReadUDPPacket(b[:length], remoteAddr, err)
		if err != nil {
			break
		}
	}
}

func (u *UDPConnection) Write(b []byte) error {
	_, err := u.conn.Write(b)
	return err
}

func (u *UDPConnection) Dispose() error {
	return u.conn.Close()
}
