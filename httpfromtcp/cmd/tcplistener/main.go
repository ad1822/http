package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ad1822/httpfromtcp/internal/request"
)

func main() {
	listener, _ := net.Listen("tcp", "localhost:42069")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Connection Accepted", conn.RemoteAddr())

		// NOTE: Now Reading from request, No File
		rl, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Errorf("Error in RequestFromReader method: %v", err)
		}

		fmt.Println("Request Line :")
		fmt.Println("- Method: ", rl.RequestLine.Method)
		fmt.Println("- Target: ", rl.RequestLine.RequestTarget)
		fmt.Println("- HttpVersion: ", rl.RequestLine.HttpVersion)

		fmt.Println("Headers :")
		rl.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})
	}

}
