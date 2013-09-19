package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
)

func ClientHandler(conn net.Conn) {
	remote := conn.RemoteAddr().String()
	for {
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal("Cannot read data from client")
		}
		log.Printf("%s > %s", remote, status)
		io.Copy(conn, bytes.NewBufferString(status))
	}
}

func main() {
	log.Printf("Start application")
	port := ":8080"
	log.Printf("Listening on %s", port)
	psock, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("Cannot listen on port %s", port)
	}

	for {
		conn, err := psock.Accept()
		if err != nil {
			log.Panic("Cannot accept connection")
		}
		log.Printf("Accepted connection with %s", conn.RemoteAddr().String())
		go ClientHandler(conn)
	}
}
