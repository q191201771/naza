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
	"sync"
	"testing"

	"github.com/q191201771/naza/pkg/assert"

	"github.com/q191201771/naza/pkg/nazanet"
)

func TestUDPConnection(t *testing.T) {
	p := nazanet.NewAvailUDPConnPool(4000, 8000)
	srvConn, srvPort, err := p.Acquire()
	assert.Equal(t, nil, err)
	laddr := fmt.Sprintf(":%d", srvPort)
	srv := nazanet.NewUDPConnectionWithConn(srvConn)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		var count int
		err := srv.RunLoop(func(b []byte, raddr *net.UDPAddr, err error) bool {
			count++
			if count > 2 {
				return true
			}
			assert.Equal(t, []byte("hello"), b)
			err2 := srv.Write2Addr([]byte("world"), raddr)
			assert.Equal(t, nil, err2)
			return true
		})
		// 因为server loop是通过Dispose强行关闭的，所以这里error有值
		assert.IsNotNil(t, err)
	}()

	cli, err := nazanet.NewUDPConnection("", laddr)
	assert.Equal(t, nil, err)
	go func() {
		err := cli.Write([]byte("hello"))
		assert.Equal(t, nil, err)
		err = cli.RunLoop(func(b []byte, raddr *net.UDPAddr, err error) bool {
			assert.Equal(t, []byte("world"), b)
			return false
		})
		assert.Equal(t, nil, err)
		wg.Done()
	}()

	cli2, err := nazanet.NewUDPConnection("", laddr)
	assert.Equal(t, nil, err)
	go func() {
		err := cli2.Write([]byte("hello"))
		assert.Equal(t, nil, err)
		err = cli2.RunLoop(func(b []byte, raddr *net.UDPAddr, err error) bool {
			assert.Equal(t, []byte("world"), b)
			return false
		})
		assert.Equal(t, nil, err)
		wg.Done()
	}()

	wg.Wait()

	err = srv.Dispose()
	assert.Equal(t, nil, err)
}
