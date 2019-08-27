package connstat2

import (
	"github.com/q191201771/nezha/assert"
	"net"
	"testing"
	"time"
)

type MockServer struct {
	l net.Listener
}

func (ms *MockServer)start() (err error) {
	ms.l, err = net.Listen("tcp", ":10027")
	if err != nil {
		return err
	}
	go func() {
		for {
			conn, err := ms.l.Accept()
			if err != nil {
				//fmt.Println(err)
				return
			}
			go func() {
				buf := make([]byte, 8)
				for {
					_, err = conn.Read(buf)
					if err != nil {
						//fmt.Println(err)
						return
					}
				}
			}()
		}
	}()
	return
}

func (ms *MockServer)stop() {
	ms.l.Close()
}

var n = 1000

func BenchmarkRaw(b *testing.B) {
	var ms MockServer
	ms.start()
	buf := make([]byte, 8)

	var conns []net.Conn

	for i := 0; i < n; i++ {
		conn, err := net.Dial("tcp", ":10027")
		assert.Equal(b, nil , err)
		conns = append(conns, conn)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			conns[j].Write(buf)
		}
	}

	ms.stop()
}

func BenchmarkDeadline(b *testing.B) {
	var ms MockServer
	ms.start()
	buf := make([]byte, 8)

	var conns []net.Conn

	for i := 0; i < n; i++ {
		conn, err := net.Dial("tcp", ":10027")
		assert.Equal(b, nil , err)
		conns = append(conns, conn)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			conns[j].SetWriteDeadline(time.Now().Add(1 * time.Second))
			conns[j].Write(buf)
		}
	}

	ms.stop()
}
