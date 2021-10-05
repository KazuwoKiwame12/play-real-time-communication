package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	addr := os.Getenv("ADDRESS")
	fmt.Println(addr)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var i int = 0
	for {
		i += 1
		time.Sleep(1 * time.Second)
		msg := fmt.Sprintf("hello %d", i)
		conn.Write([]byte(msg))
		time.Sleep(1 * time.Second)
		msgBuf := make([]byte, 1024)
		msgLen, err := conn.Read(msgBuf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Recieve: %s\n", string(msgBuf[:msgLen]))
	}
}
