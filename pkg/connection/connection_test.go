// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package connection_test

import (
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/connection"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
)

func TestWriteTimeout(t *testing.T) {
	// 开启一个 tcp 服务器，只accept一个连接，之后对这个连接不做任何读写
	// 使用 Connection 设置写超时后，死循环往服务器发送数据
	ch := make(chan struct{}, 1)
	l, err := net.Listen("tcp", ":10027")
	assert.Equal(t, nil, err)
	defer l.Close()
	go func() {
		srvConn, _ := l.Accept()
		defer srvConn.Close()
		<-ch
	}()
	conn, err := net.Dial("tcp", ":10027")
	c := connection.New(conn, func(opt *connection.Option) {
		opt.WriteTimeoutMs = 1000
	})
	assert.Equal(t, nil, err)
	b := make([]byte, 128*1024)
	for {
		n, err := c.Write(b)
		nazalog.Infof("%d %+v", n, err)
		if err != nil {
			break
		}
	}
	c.Close()
	ch <- struct{}{}
}

func TestWrite(t *testing.T) {
	// TODO(chef):
	// - 使用testWithConnPair
	// - 提高覆盖率

	var sentN uint32
	var sentDone uint32

	rand.Seed(time.Now().Unix())
	l, err := net.Listen("tcp", ":10027")
	assert.Equal(t, nil, err)
	defer l.Close()
	go func() {
		c, err := l.Accept()
		srvConn := connection.New(c, func(option *connection.Option) {
			option.WriteChanSize = 1024
			//option.WriteBufSize = 256
			option.WriteTimeoutMs = 10000
		})
		assert.Equal(t, nil, err)
		for i := 0; i < 10; i++ {
			b := make([]byte, rand.Intn(4096))
			n, err := srvConn.Write(b)
			if err == nil {
				nazalog.Debugf("sent. i=%d, n=%d", i, n)
			}
			assert.Equal(t, nil, err)
			atomic.AddUint32(&sentN, uint32(n))
		}
		err = srvConn.Flush()
		assert.Equal(t, nil, err)
		nazalog.Debugf("total sent:%d", sentN)
		atomic.StoreUint32(&sentDone, 1)
		srvConn.Close()
	}()

	conn, err := net.Dial("tcp", ":10027")
	assert.Equal(t, nil, err)
	b := make([]byte, 4096)
	var readN uint32
	for {
		n, _ := conn.Read(b)
		readN += uint32(n)
		nazalog.Debugf("total read:%d", readN)
		if atomic.LoadUint32(&sentDone) == 1 && atomic.LoadUint32(&sentN) == readN {
			break
		}
	}
	conn.Close()
}

func TestConnection_Writev(t *testing.T) {
	goldenB := []byte{'a', 'b', 'c', 'd', 'e'}

	testWithConnPair(t, func(srvConn, cliConn net.Conn) {
		c := connection.New(cliConn, func(option *connection.Option) {
		})
		n, err := c.Writev(net.Buffers{goldenB[:2], goldenB[2:5]})
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		recvB := make([]byte, 4096)
		n, err = srvConn.Read(recvB)
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, goldenB[:5], recvB[:5])

		c.Close()
		srvConn.Close()
	})

	testWithConnPair(t, func(srvConn, cliConn net.Conn) {
		c := connection.New(cliConn, func(option *connection.Option) {
			option.WriteChanSize = 128
		})
		n, err := c.Writev(net.Buffers{goldenB[:2], goldenB[2:5]})
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		recvB := make([]byte, 4096)
		n, err = srvConn.Read(recvB)
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, goldenB[:5], recvB[:5])

		c.Close()
		srvConn.Close()
	})

	testWithConnPair(t, func(srvConn, cliConn net.Conn) {
		c := connection.New(cliConn, func(option *connection.Option) {
			option.WriteBufSize = 1024
		})
		n, err := c.Writev(net.Buffers{goldenB[:2], goldenB[2:5]})
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		err = c.Flush()
		assert.Equal(t, nil, err)
		recvB := make([]byte, 4096)
		n, err = srvConn.Read(recvB)
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, goldenB[:5], recvB[:5])

		c.Close()
		srvConn.Close()
	})
}

func testWithConnPair(t *testing.T, cb func(srvConn, cliConn net.Conn)) {
	l, err := net.Listen("tcp", ":10027")
	assert.Equal(t, nil, err)
	defer l.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	var srvConn net.Conn
	go func() {
		var err error
		srvConn, err = l.Accept()
		assert.Equal(t, nil, err)
		wg.Done()
	}()

	cliConn, err := net.Dial("tcp", ":10027")
	assert.Equal(t, nil, err)
	wg.Wait()

	cb(srvConn, cliConn)
}
