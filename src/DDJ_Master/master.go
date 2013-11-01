package main

import (
	"config"
	"container/list"
	"flag"
	"fmt"
	"log"
	"net"
	"node"
	"os"
	"rest"
	"strconv"
)

// Main: Starts a TCP server and waits infinitely for connections
func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println("Start Master")

	nodeList := list.New()
	in2 := rest.Channel.Get()
	in := make(chan string)

	go node.IOHandler(in2, nodeList)

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
