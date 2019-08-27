package main

import "net"

// MockServer 开启一个TCP监听，并从accept的连接上循环读取数据

type MockServer struct {
	l net.Listener
}

func (ms *MockServer) start(addr string) (err error) {
	ms.l, err = net.Listen("tcp", addr)
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

func (ms *MockServer) stop() {
	ms.l.Close()
}
