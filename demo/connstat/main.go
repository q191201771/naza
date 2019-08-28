package main

import (
	"flag"
	"fmt"
	"github.com/q191201771/nezha/pkg/errors"
	"github.com/q191201771/nezha/pkg/log"
	"net"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

// 测试 net.Conn.SetWriteDeadline 的性能开销
//
// 使用 <numOfConn> 个客户端端连接，并行向服务端发送数据
// 每个客户端发送 <numOfMsgPerConn> 个消息
// 服务端 MockServer 只接收数据，不做其他逻辑

const (
	addr = ":10027"
)

var (
	numOfConn       int
	numOfMsgPerConn int
)

func raw(writeTimeoutSec int) {
	log.Infof("> raw. writeTimeoutSec=%d", writeTimeoutSec)
	var ms MockServer
	var conns []net.Conn
	buf := make([]byte, 8)

	ms.start(addr)
	for i := 0; i < numOfConn; i++ {
		conn, err := net.Dial("tcp", addr)
		errors.PanicIfErrorOccur(err)
		conns = append(conns, conn)
	}

	var wg sync.WaitGroup
	wg.Add(numOfConn)

	b := time.Now()
	log.Infof("b:%+v", b)
	fp, err := os.Create(fmt.Sprintf("profile.out.%d", writeTimeoutSec))
	errors.PanicIfErrorOccur(err)
	defer fp.Close()
	err = pprof.StartCPUProfile(fp)
	errors.PanicIfErrorOccur(err)
	for i := 0; i < numOfConn; i++ {
		go func(ii int) {
			for j := 0; j < numOfMsgPerConn; j++ {
				if writeTimeoutSec != 0 {
					err := conns[ii].SetWriteDeadline(time.Now().Add(time.Duration(writeTimeoutSec) * time.Second))
					errors.PanicIfErrorOccur(err)
				}
				_, err := conns[ii].Write(buf)
				errors.PanicIfErrorOccur(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	pprof.StopCPUProfile()
	log.Infof("cost=%v", time.Now().Sub(b))

	for i := 0; i < numOfConn; i++ {
		err := conns[i].Close()
		errors.PanicIfErrorOccur(err)
	}
	ms.stop()
	log.Info("< raw.")
}

func main() {
	c := flag.Int("c", 0, "num of conn")
	n := flag.Int("n", 0, "num of msg per conn")
	flag.Parse()
	if *c == 0 || *n == 0 {
		flag.Usage()
		return
	}
	numOfConn = *c
	numOfMsgPerConn = *n

	raw(0)
	raw(5)
}
