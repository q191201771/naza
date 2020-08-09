// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazanet_test

import (
	"net"
	"testing"

	"github.com/q191201771/naza/pkg/assert"

	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/nazanet"
)

func TestAvailUDPConnPool_Acquire(t *testing.T) {
	var conns []*net.UDPConn
	aucp := nazanet.NewAvailUDPConnPool(8000, 8005)
	closedOnlyOnceFlag := false

	// 循环次数大于端口范围，测试后面的获取是否返回错误
	for i := 0; i < 10; i++ {
		conn, port, err := aucp.Acquire()
		nazalog.Debugf("%d: %p, %d, %v", i, conn, port, err)

		// 关闭一次，看下次是否能复用
		if !closedOnlyOnceFlag {
			if err == nil {
				err = conn.Close()
				assert.Equal(t, nil, err)
				closedOnlyOnceFlag = true
			}
			continue
		}

		conns = append(conns, conn)
	}

	for _, conn := range conns {
		if conn != nil {
			conn.Close()
		}
	}
}

func TestAvailUDPConnPool_Acquire2(t *testing.T) {
	aucp := nazanet.NewAvailUDPConnPool(8000, 8005)
	closedOnlyOnceFlag := false

	// 循环次数大于端口范围，测试后面的获取是否返回错误
	for i := 0; i < 10; i++ {
		conn1, port1, conn2, port2, err := aucp.Acquire2()
		nazalog.Debugf("%d: %p, %d, %p, %d, %v", i, conn1, port1, conn2, port2, err)

		// 关闭一次，看下次是否能复用
		if !closedOnlyOnceFlag {
			if err == nil {
				err = conn1.Close()
				assert.Equal(t, nil, err)
				closedOnlyOnceFlag = true
			}
		}
	}
}

func TestAvailUDPConnPool_Peek(t *testing.T) {
	aucp := nazanet.NewAvailUDPConnPool(8000, 8005)
	nazalog.Debug(aucp.Peek())
}
