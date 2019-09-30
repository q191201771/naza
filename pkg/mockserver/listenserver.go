package mockserver

import (
	"github.com/q191201771/naza/pkg/nazalog"
	"net"
	"time"
)

// 建立一个server端的监听，在内部创建n个连接快速消耗掉listen队列，达到对外模拟不处理连接的情况

type MockListenServer struct {
	l net.Listener
}

func (s *MockListenServer) Run(addr string) (err error) {
	s.l, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	for i := 0; ; i++ {
		if _, err := net.DialTimeout("tcp", addr, time.Duration(200)*time.Millisecond); err != nil {
			nazalog.Infof("Dial failed. i=%d, err=%+v", i, err)
			break
		}

	}
	return
}

func (s *MockListenServer) Dispose() {
	_ = s.l.Close()
}
