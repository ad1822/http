package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	sender, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error in sending packets : ", err)
	}

	conn, err := net.DialUDP("udp", nil, sender)

	if err != nil {
		fmt.Println("Error : ", err)
	}

	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in ReadString : ", err)
		}
		fmt.Println(line)

		str, err := conn.Write([]byte(line))
		if err != nil {
			fmt.Println("Error in Write : ", err)
		}
		fmt.Println(str)
	}
}
