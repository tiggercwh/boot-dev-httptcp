package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:2025")
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}

	defer conn.Close()

	for {
		fmt.Println(">")
		reader := bufio.NewReader(os.Stdin)
		b, err := reader.ReadSlice('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		conn.Write(b)
	}
}
