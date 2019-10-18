// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package connection

import (
	"net"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
)

// TODO chef: 补充单元测试

func TestWriteTimeout(t *testing.T) {
	// 开启一个 tcp 服务器，只accept一个连接，之后对这个连接不做任何读写
	// 使用 Connection 设置写超时后，死循环往服务器发送数据
	ch := make(chan struct{}, 1)
	l, err := net.Listen("tcp", ":10027")
	assert.Equal(t, nil, err)
	defer l.Close()
	go func() {
		conn, _ := l.Accept()
		defer conn.Close()
		<-ch
	}()
	conn, err := net.Dial("tcp", ":10027")
	c := New(conn, func(opt *Option) {
		opt.WriteTimeoutMS = 1000
	})
	assert.Equal(t, nil, err)
	b := make([]byte, 128)
	for {
		n, err := c.Write(b)
		nazalog.Infof("%d %+v", n, err)
		if err != nil {
			break
		}
	}
	ch <- struct{}{}
}
