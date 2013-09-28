package main

import (
	"client"
	"container/list"
	"log"
	"net"
	"os"
)

// Main: Starts a TCP server and waits infinitely for connections
func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println("Start Master")

	getService()

	clientList := list.New()
	in := make(chan string)
	go client.IOHandler(in, clientList)
	service := "127.0.0.1:8080"
	tcpAddr, error := net.ResolveTCPAddr("tcp", service)
	if error != nil {
		log.Println("Error: Could not resolve address")
	} else {
		log.Println("Listening on: ", tcpAddr.String())
		netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
		if error != nil {
			log.Println(error)
		} else {
			defer netListen.Close()
			for {
				log.Println("Waiting for clients")
				connection, error := netListen.Accept()
				if error != nil {
					log.Println("Client error: ", error)
				} else {
					log.Println("Accept client")
					go client.ClientHandler(connection, in, clientList)
				}
			}
		}
	}
}

func getService() string {

	log.Println("Program arguments: ", os.Args[1:])

	return os.Args[1]
}
