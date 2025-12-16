package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	_ = err

	conn.Write([]byte("hello"))
}