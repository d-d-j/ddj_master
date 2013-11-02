package main

import (
	"config"
	"container/list"
	"fmt"
	"log"
	"net"
	"node"
	"rest"
)

// Main: Starts a TCP server and waits infinitely for connections
func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println("Start Master")

	nodeList := list.New()

	go node.IOHandler(rest.Channel.QueryChannel(), nodeList)

	log.Print("Load configuration: ")
	cfg, err := config.Load()
	if err != nil {
		log.Panic("Problem with configuration: ", err)
	}
	log.Println(cfg)

	portApi := fmt.Sprintf(":%d", cfg.Ports.Api)
	rest.StartApi(portApi)

	service := fmt.Sprintf("127.0.0.1:%d", cfg.Ports.NodeCommunication)
	tcpAddr, error := net.ResolveTCPAddr("tcp", service)
	if error != nil {
		log.Panic("Error: Could not resolve address")
	}

	log.Println("Listening on: ", tcpAddr.String())
	netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
	if error != nil {
		log.Panic(error)
	}
	defer netListen.Close()

	WaitForNodes(netListen, nodeList)

}

func WaitForNodes(netListen net.Listener, nodes *list.List) {
	in := make(chan string)
	for {
		log.Println("Waiting for nodes")
		connection, error := netListen.Accept()
		if error != nil {
			log.Println("node error: ", error)
		} else {
			log.Println("Accept node: ", connection.RemoteAddr())
			go node.NodeHandler(connection, in, nodes)
		}
	}
}
