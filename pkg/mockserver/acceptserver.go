package mockserver

import (
	"net"
	"sync"
)

type MockAcceptServer struct {
	l net.Listener
	conns []net.Conn
	m sync.Mutex
}

func (s *MockAcceptServer) Run(addr string) (err error) {
	s.m.Lock()
	s.l, err = net.Listen("tcp", addr)
	s.m.Unlock()
	if err != nil {
		return
	}
	c, err := s.l.Accept()
	if err != nil {
		return
	}
	s.conns = append(s.conns, c)
	return
}

func (s *MockAcceptServer) Dispose() {
	s.m.Lock()
	if s.l != nil {
		s.l.Close()

	}
	s.m.Unlock()
}

