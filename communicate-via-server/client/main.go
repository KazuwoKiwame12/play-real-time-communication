package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		stdin := bufio.NewScanner(os.Stdin)
		for stdin.Scan() {
			conn.Write([]byte(stdin.Text()))
		}
	}()

	for {
		msgBuf := make([]byte, 1024)
		msgLen, err := conn.Read(msgBuf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(msgBuf[:msgLen]))
	}
}
