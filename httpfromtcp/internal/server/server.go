package server

import (
	"fmt"
	"io"
	"net"

	"github.com/ad1822/httpfromtcp/internal/response"
)

type Server struct {
	closed bool
}

func runConnection(conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(10)
	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, headers)
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
