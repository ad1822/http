package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}

			buffer = buffer[:n]
			if i := bytes.IndexByte(buffer, '\n'); i != -1 {
				str += string(buffer[:i])
				buffer = buffer[i+1:]
				out <- str
				str = ""
			}
			str += string(buffer)
		}
	}()
	return out
}

func main() {
	// openFile, err := os.Open("message.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer openFile.Close()
	//
	// lines := getLinesChannel(openFile)
	// for line := range lines {
	// 	fmt.Printf("read: %s\n", line)
	// }
	listener, _ := net.Listen("tcp4", "localhost:42069")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Connection Accepted", conn.RemoteAddr())

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
