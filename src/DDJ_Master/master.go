package main

import (
	"container/list"
	"flag"
	"log"
	"net"
	"node"
	"os"
	"strconv"
)

// Main: Starts a TCP server and waits infinitely for connections
func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println("Start Master")

	nodeList := list.New()
	in := make(chan string)
	go node.IOHandler(in, nodeList)
	service := "127.0.0.1:" + getPortFromArgument()
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
				log.Println("Waiting for nodes")
				connection, error := netListen.Accept()
				if error != nil {
					log.Println("node error: ", error)
				} else {
					log.Println("Accept node")
					go node.NodeHandler(connection, in, nodeList)

				}
			}
		}
	}
}

func getPortFromArgument() string {

	log.Println("Program arguments: ", os.Args[1:])
	port := flag.Int("port", 8080, "port number")
	flag.Parse()
	log.Println("Port: ", *port)
	return strconv.Itoa(*port)
}
