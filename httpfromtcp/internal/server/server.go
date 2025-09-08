package server

import (
	"fmt"
	"io"
	"net"

	"github.com/ad1822/httpfromtcp/internal/request"
	"github.com/ad1822/httpfromtcp/internal/response"
)

type Server struct {
	closed  bool
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
	error      error
}

type Handler func(w *response.Writer, req *request.Request)

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	responseWriter := response.NewWriter(conn)

	request, err := request.RequestFromReader(conn)

	// NOTE: Bad Request
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
		return
	}

	s.handler(responseWriter, request)
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		fmt.Printf("Connection : %s\n", conn.RemoteAddr())
		if s.closed {
			return
		}
		if err != nil {
			return
		}
		go runConnection(s, conn)
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{closed: false, handler: handler}
	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
