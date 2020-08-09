// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazanet_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazanet"
)

func TestUDPConnection(t *testing.T) {
	port, err := nazanet.NewAvailUDPConnPool(8000, 16000).Peek()
	assert.Equal(t, nil, err)

	addr := fmt.Sprintf(":%d", port)
	conn := nazanet.NewUDPConnectionWithLocalAddr(addr, func(b []byte, remoteAddr net.Addr, err error) {
		nazalog.Debugf("%d, %v, %v", len(b), remoteAddr, err)
	})
	err = conn.Listen()
	assert.Equal(t, nil, err)
	go conn.RunLoop()
	err = conn.Dispose()
	assert.Equal(t, nil, err)
	time.Sleep(100 * time.Millisecond)
}

func TestNewUDPConnectionWithConn(t *testing.T) {
	conn, _, err := nazanet.NewAvailUDPConnPool(8000, 16000).Acquire()
	assert.Equal(t, nil, err)

	connection := nazanet.NewUDPConnectionWithConn(conn, func(b []byte, remoteAddr net.Addr, err error) {

	})
	go connection.RunLoop()
	_ = connection.Dispose()
}
