package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(conn io.ReadWriteCloser) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World! ")
	conn.Write(out)
	conn.Close()
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		fmt.Printf("Connection : %s\n", conn.RemoteAddr())
		if s.closed {
			return
		}
		if err != nil {
			fmt.Errorf("Can't accept connection : %s", err)
		}
		go runConnection(conn)
	}
}

func Serve(port uint16) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{closed: false}
	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) listen() {
}
