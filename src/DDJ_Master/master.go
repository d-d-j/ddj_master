package main

import (
	"code.google.com/p/gorest"
	"container/list"
	"dto"
	"flag"
	"log"
	"net"
	"net/http"
	"node"
	"os"
	"strconv"
)

// Main: Starts a TCP server and waits infinitely for connections
func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println("Start Master")

	nodeList := list.New()
	in2 := Channel.Get()
	in := make(chan string)

	go node.IOHandler(in2, nodeList)

	insertService := InsertService{}

	log.Println("Channel: ", in2)

	gorest.RegisterService(&insertService) //Register our service
	gorest.RegisterMarshaller("application/json", gorest.NewJSONMarshaller())
	go http.Handle("/", gorest.Handle())
	go http.ListenAndServe(":8081", nil)

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

type InsertChannel struct {
	channel chan dto.Element
}

func (s *InsertChannel) Get() chan dto.Element {
	return s.channel
}

var Channel interface {
	Get() chan dto.Element
} = &InsertChannel{make(chan dto.Element)}

//Service Definition
type InsertService struct {
	gorest.RestService `root:"/"`
	insertData         gorest.EndPoint `method:"POST" path:"/series/id/{id:int32}/data/" postdata:"dto.Element"`
}

func (serv InsertService) InsertData(posted dto.Element, id int32) {
	log.Println("Inserting new data to series: ", id)
	log.Println("Data to insert: ", posted)
	log.Println("Channel: ", Channel.Get())

	Channel.Get() <- posted

}

func getPortFromArgument() string {

	log.Println("Program arguments: ", os.Args[1:])
	port := flag.Int("port", 8080, "port number")
	flag.Parse()
	log.Println("Port: ", *port)
	return strconv.Itoa(*port)
}
