package connstat2

import (
	"github.com/q191201771/nezha/assert"
	"net"
	"testing"
)

func startMockServer(t assert.TestingT) {
	l, err := net.Listen("tcp", ":10027")
	assert.Equal(t, nil , err)
	go func() {
		//for {
			conn, err := l.Accept()
			assert.Equal(t, nil , err)
			go func() {
				buf := make([]byte, 8)
				for {
					_, err = conn.Read(buf)
					assert.Equal(t, nil , err)
				}
			}()
		//}
	}()
}

func Benchmark(b *testing.B) {
	startMockServer(b)
	//time.Sleep(time.Duration(1) * time.Second)

	conn, err := net.Dial("tcp", ":10027")
	assert.Equal(b, nil , err)
	buf := make([]byte, 8)
	for i := 0; i < b.N; i++ {
		conn.Write(buf)
	}
}