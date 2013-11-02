//http://dev.badgerr.co.uk/erlsrv/file/4f624b6bc22c/gostuff/chatserv.go
package node

// Imports required packages
import (
	"bytes"
	"container/list"
	"dto"
	"log"
	"net"
)

// Defines a Node with a name and connection object, and
// some channels for sending and receiving text.
// Also holds a pointer to the "global" list of all connected Nodes
type Node struct {
	Name     string
	Incoming chan []byte
	Outgoing chan string
	Conn     net.Conn
	Quit     chan bool
	NodeList *list.List
}

// Defines a read function for a Node, reading from the connection into
// a buffer passed in. Returns true if read was successful, false otherwise
func (c *Node) Read(buffer []byte) bool {
	bytesRead, error := c.Conn.Read(buffer)
	if error != nil {
		c.Close()
		log.Println(error)
		return false
	}
	log.Println("Read ", bytesRead, " bytes")
	return true
}

// Closes a Node connection and removes it from the Node list
func (c *Node) Close() {
	c.Quit <- true
	c.Conn.Close()
	c.RemoveMe()
}

// Comparison function to easily check equality with another Node
// based on the name and connection
func (c *Node) Equal(other *Node) bool {
	if bytes.Equal([]byte(c.Name), []byte(other.Name)) {
		if c.Conn == other.Conn {
			return true
		}
	}
	return false
}

// Removes this Node from the Node list
func (c *Node) RemoveMe() {
	for entry := c.NodeList.Front(); entry != nil; entry = entry.Next() {
		Node := entry.Value.(Node)
		if c.Equal(&Node) {
			log.Println("RemoveMe: ", c.Name)
			c.NodeList.Remove(entry)
		}
	}
}

// Server listener goroutine - waits for data from the incoming channel
// (each Node.Outgoing stores this), and passes it to each Node.Incoming channel
func IOHandler(Query <-chan dto.Query, NodeList *list.List) {
	for {
		query := <-Query
		log.Println("Query", query)
		header := query.TaskRequestHeader

		var (
			buf       []byte
			headerBuf []byte
			err       error
		)

		if query.Load != nil {
			buf, err = query.Load.Encode()
			if err != nil {
				log.Println(err)
				continue
			}
		}

		header.Size = (int32)(len(buf))
		headerBuf, err = header.Encode()
		if err != nil {
			log.Println(err)
			continue
		}

		complete := make([]byte, 100)
		copy(complete, headerBuf)
		copy(complete[len(headerBuf):], buf)

		//TODO: Replace this with StoreManager
		for e := NodeList.Front(); e != nil; e = e.Next() {
			Node := e.Value.(Node)

			Node.Incoming <- complete
		}
	}
}

// Node reading goroutine - reads incoming data from the tcp socket,
// sends it to the Node.Outgoing channel (to be picked up by IOHandler)
func NodeReader(Node *Node) {
	buffer := make([]byte, 2048)

	for Node.Read(buffer) {
		log.Println("NodeReader received ", Node.Name, "> ", string(buffer))
		for i := 0; i < 2048; i++ {
			buffer[i] = 0x00
		}
	}

	log.Println("NodeReader stopped for ", Node.Name)
	Node.Close()
}

// Node sending goroutine - waits for data to be sent over Node.Incoming
// (from IOHandler), then sends it over the socket
func NodeSender(Node *Node) {
	for {
		select {
		case buffer := <-Node.Incoming:
			log.Println("NodeSender sending ", buffer, " to ", Node.Name)
			Node.Conn.Write(buffer)
		case <-Node.Quit:
			log.Println("Node ", Node.Name, " quitting")
			Node.Conn.Close()
			break
		}
	}
}

// Creates a new Node object for each new connection using the name sent by the Node,
// then starts the NodeSender and NodeReader goroutines to handle the IO
func NodeHandler(conn net.Conn, ch chan string, NodeList *list.List) {
	buffer := make([]byte, 1024)
	bytesRead, error := conn.Read(buffer)
	if error != nil {
		log.Println("Node connection error: ", error)
	}
	name := string(buffer[0:bytesRead])
	newNode := &Node{name, make(chan []byte), make(chan string), conn, make(chan bool), NodeList}
	go NodeSender(newNode)
	//go NodeReader(newNode)
	NodeList.PushBack(*newNode)
	// ch <- string(name + " has joined")
}
