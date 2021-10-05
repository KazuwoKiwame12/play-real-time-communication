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
	lister, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := lister.Accept()
		if err != nil {
			continue
		}
		go func() {
			time.Sleep(10 * time.Second)
			conn.Close()
		}()
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	for {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		msgBuf := make([]byte, 1024)
		msgLen, err := conn.Read(msgBuf)
		if err != nil {
			fmt.Println("time out")
			return
		}

		msg := string(msgBuf[:msgLen])
		fmt.Printf("Recieve: %s\n", msg)
		msg = fmt.Sprintf("owl response: %s\n", msg)
		conn.Write([]byte(msg))
	}
}
