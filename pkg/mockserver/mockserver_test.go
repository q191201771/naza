package mockserver

import (
	"github.com/q191201771/nezha/pkg/assert"
	"net"
	"testing"
	"time"
)

var addr = ":10027"

func TestMockListenServer(t *testing.T) {
	//var s MockListenServer
	//err := s.Run(addr)
	//assert.Equal(t, nil, err)
	//defer s.Dispose()
	//_, err = net.DialTimeout("tcp", addr, time.Duration(1000) * time.Millisecond)
	//assert.IsNotNil(t, err)
}

func TestMockAcceptServer(t *testing.T) {
	var s MockAcceptServer
	var conns []net.Conn
	go s.Run(addr)
	for i := 0; i < 16; i++ {
		c, err := net.DialTimeout("tcp", addr, time.Duration(1000)*time.Millisecond)
		if err != nil {
			break
		}
		//assert.Equal(t, nil, err)
		conns = append(conns, c)
	}
	s.Dispose()
}

func TestCorner(t *testing.T) {
	var s MockListenServer
	err := s.Run("wrong addr")
	assert.IsNotNil(t, err)
}
