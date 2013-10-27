package main

import (
	"bufio"
	"code.google.com/p/gorest"
	"container/list"
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

	gorest.RegisterService(new(InsertService)) //Register our service
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":8081", nil)

	nodeList := list.New()
	in2 := make(chan string)
	in := make(chan string)

	go ReadStdIn(in2)

	go node.IOHandler(in2, nodeList)
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

//Service Definition
type InsertService struct {
	gorest.RestService `root:"/tutorial/"`
	helloWorld         gorest.EndPoint `method:"GET" path:"/hello-world/" output:"string"`
	sayHello           gorest.EndPoint `method:"GET" path:"/hello/{name:string}" output:"string"`
}

func (serv InsertService) HelloWorld() string {
	return "Hello World"
}
func (serv InsertService) SayHello(name string) string {
	return "Hello " + name
}

func ReadStdIn(Incoming chan string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		Incoming <- line
		if err != nil {
			// You may check here if err == io.EOF
			break
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
