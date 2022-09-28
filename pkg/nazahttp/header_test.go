// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazahttp_test

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/q191201771/naza/pkg/assert"

	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/nazahttp"
)

func TestHeader(t *testing.T) {
	// TODO(chef): 抽象出可获取可用监听端口
	for port := 8080; port != 8090; port++ {
		addr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}

		go func() {
			_, _ = nazahttp.GetHttpFile(fmt.Sprintf("http://%s/test", addr), 100)
		}()

		conn, err := ln.Accept()
		r := bufio.NewReader(conn)
		fl, hs, err := nazahttp.ReadHttpHeader(r)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, len(hs) > 0)
		nazalog.Debugf("first line:%s", fl)
		nazalog.Debugf("header fields:%+v", hs)

		m, u, v, err := nazahttp.ParseHttpRequestLine(fl)
		assert.Equal(t, nil, err)
		nazalog.Debugf("method:%s, uri:%s, version:%s", m, u, v)
		assert.Equal(t, "GET", m)
		assert.Equal(t, "/test", u)
		assert.Equal(t, "HTTP/1.1", v)

		_ = conn.Close()
		_ = ln.Close()
		break
	}
}

func TestReadHttpResponseMessage(t *testing.T) {
	// 兼容性case1
	for port := 8080; port != 8090; port++ {
		addr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			conn, err := net.Dial("tcp", addr)
			assert.Equal(t, nil, err)
			conn.Write([]byte("PLAY rtsp://127.0.0.1:5544/live/test110 RTSP/1.0\r\nUser-Agent: lal/0.26.0\r\nSession: 191201771\r\nRange: npt=0.000-\r\nCSeq: 5\r\n\r\n"))
			r := bufio.NewReader(conn)
			ctx, err := nazahttp.ReadHttpResponseMessage(r)
			assert.Equal(t, nil, err)
			assert.Equal(t, "url=track_id=0;seq=63248;rtptime=0,url=track_id=1;seq=56208;rtptime=0", ctx.Headers.Get("Rtp-Info"))
			wg.Done()
		}()

		conn, err := ln.Accept()
		r := bufio.NewReader(conn)
		_, err = nazahttp.ReadHttpResponseMessage(r)
		assert.Equal(t, nil, err)

		_, _ = conn.Write([]byte("RTSP/1.0 200 OK\r\nCSeq: 5\r\nSession: ac5a1f04\r\nRTP-Info: url=track_id=0;seq=63248;rtptime=0,\r\nurl=track_id=1;seq=56208;rtptime=0\r\n\r\n"))
		wg.Wait()

		_ = conn.Close()
		_ = ln.Close()
		break
	}
}

func TestParseHttpStatusLine(t *testing.T) {
	v, c, r, e := nazahttp.ParseHttpStatusLine("HTTP/1.0 200 OK")
	assert.Equal(t, nil, e)
	assert.Equal(t, "HTTP/1.0", v)
	assert.Equal(t, "200", c)
	assert.Equal(t, "OK", r)

	v, c, r, e = nazahttp.ParseHttpStatusLine("HTTP/1.1 400 Bad Request")
	assert.Equal(t, nil, e)
	assert.Equal(t, "HTTP/1.1", v)
	assert.Equal(t, "400", c)
	assert.Equal(t, "Bad Request", r)

	//statusLine := "HTTP/1.1 400 "
	//for i := 0; i <= len(statusLine); i++ {
	//	sl := statusLine[0:i]
	//	_, _, _, e = nazahttp.ParseHttpStatusLine(sl)
	//	assert.IsNotNil(t, e, sl)
	//}

	v, c, r, e = nazahttp.ParseHttpStatusLine("HTTP/1.1 475 ")
	assert.Equal(t, nil, e)
	assert.Equal(t, "HTTP/1.1", v)
	assert.Equal(t, "475", c)
	assert.Equal(t, "", r)

	v, c, r, e = nazahttp.ParseHttpStatusLine("HTTP/1.1 475")
	assert.Equal(t, nil, e)
	assert.Equal(t, "HTTP/1.1", v)
	assert.Equal(t, "475", c)
	assert.Equal(t, "", r)

	v, c, r, e = nazahttp.ParseHttpStatusLine("HTTP/1.1 475  ")
	assert.Equal(t, nil, e)
	assert.Equal(t, "HTTP/1.1", v)
	assert.Equal(t, "475", c)
	assert.Equal(t, " ", r)

	// 测试解析错误的情况
	v, c, r, e = nazahttp.ParseHttpStatusLine("fxxk")
	assert.IsNotNil(t, e)
}
