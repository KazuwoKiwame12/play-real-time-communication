package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type rooms map[string]*room

type room struct {
	users  []user
	ctx    context.Context
	cancel context.CancelFunc
}

type user struct {
	name string
	conn net.Conn
}

func (r *room) broadCast(msg []byte) {
	for _, user := range r.users {
		user.conn.Write(msg)
	}
}

func (r *room) close() {
	for _, user := range r.users {
		user.conn.Close()
	}
}

var rs rooms = make(rooms, 1)

func main() {
	lister, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := lister.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	var roomKey string = ""
	var name string = "null"

	for {
		msgBuf := make([]byte, 1024)
		msgLen, err := conn.Read(msgBuf)
		if err != nil {
			fmt.Println("time out")
			return
		}

		strs := strings.Split(string(msgBuf[:msgLen]), " ")
		switch strs[0] {
		case "CREATE":
			roomKey = createKey(20)
			name = strs[1]
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			rs[roomKey] = &room{
				users: []user{
					{
						name: name,
						conn: conn,
					},
				},
				ctx:    ctx,
				cancel: cancel,
			}
			conn.Write([]byte(fmt.Sprintf("key is %s", roomKey)))
		case "ENTER":
			roomKey = strs[2]
			name = strs[1]
			r, ok := rs[roomKey]
			if !ok {
				conn.Write([]byte("inputted key doesn't exist"))
				continue
			}
			r.users = append(r.users, user{
				name: name,
				conn: conn,
			})
			conn.Write([]byte("connected!"))
		default:
			continue
		}
		break
	}

	for {
		msgBuf := make([]byte, 1024)
		msgLen, err := conn.Read(msgBuf)
		if err != nil {
			rs[roomKey].cancel()
			conn.Close()
			return
		}
		select {
		case <-rs[roomKey].ctx.Done():
			rs[roomKey].broadCast([]byte("system: tim up\n"))
			rs[roomKey].close()
			return
		default:
			rs[roomKey].broadCast([]byte(fmt.Sprintf("%s: %s\n", name, msgBuf[:msgLen])))
		}
	}
}

func createKey(length int) string {
	rand.Seed(time.Now().UnixNano())
	const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = LETTERS[rand.Intn(len(LETTERS))]
	}
	return string(result)
}
