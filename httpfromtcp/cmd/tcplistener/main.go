package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ad1822/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		log.Fatalf("failed to listen on port 42069: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		fmt.Printf("Connection Accepted %s\n", conn.RemoteAddr())

		// NOTE: Now Reading from request, No File
		rl, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("RequestFromReader failed: %v", err)
			_ = conn.Close()
			continue
		}

		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n", rl.RequestLine.Method)
		fmt.Printf("- Target: %s\n", rl.RequestLine.RequestTarget)
		fmt.Printf("- HttpVersion: %s\n", rl.RequestLine.HttpVersion)

		fmt.Printf("Headers:\n")
		rl.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})
	}
}
