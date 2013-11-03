//http://dev.badgerr.co.uk/erlsrv/file/4f624b6bc22c/gostuff/chatserv.go
package node

// Imports required packages
import (
	"container/list"
	"dto"
	"log"
	"net"
	"sync/atomic"
)

// Defines a Node with a Id and connection object, and
// some channels for sending and receiving text.
// Also holds a pointer to the "global" list of all connected Nodes
type Node struct {
	Id       int32
	Incoming chan []byte
	Outgoing chan dto.Result
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
// based on the Id and connection
func (c *Node) Equal(other *Node) bool {
	if c.Id == other.Id {
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
			log.Println("RemoveMe: ", c.Id)
			c.NodeList.Remove(entry)
		}
	}
}

func IOHandler(Query <-chan dto.Query, Result <-chan dto.Result, NodeList *list.List) {
	taskResponse := make(map[int32]chan []dto.Dto)
	for {
		select {
		case query := <-Query:
			log.Println("Query", query)
			header := query.TaskRequestHeader

			var (
				buf       []byte
				headerBuf []byte
				err       error
			)

			if query.Response != nil {
				taskResponse[query.Id] = query.Response
			}

			if query.Load != nil {
				buf, err = query.Load.Encode()
				if err != nil {
					log.Println(err)
					continue
				}
			}

			header.LoadSize = (int32)(len(buf))
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
		case result := <-Result:
			log.Println("Result: ", result.String(), result.Load)
			ch := taskResponse[result.Id]
			log.Println("Pass result data to proper client")
			ch <- result.Load
		}

	}
}

// Node reading goroutine - reads incoming data from the tcp socket,
// sends it to the Node.Outgoing channel (to be picked up by IOHandler)
func NodeReader(Node *Node) {
	buffer := make([]byte, 2048)

	for Node.Read(buffer) {
		log.Println("NodeReader received data from", Node.Id)
		var r dto.Result
		err := r.DecodeHeader(buffer)
		if err != nil {
			log.Println(err)
		}
		buffer = buffer[r.TaskRequestHeader.Size():]
		if r.Code == 2 {
			length := int(r.LoadSize / 24)
			load := make([]dto.Dto, length)
			for i := 0; i < length; i++ {
				var e dto.Element
				err = e.Decode(buffer[i*(e.Size()+4):])
				if err != nil {
					log.Println(err)
					continue
				}
				load[i] = &e
			}
			r.Load = load
		}
		log.Println("Send response to IOHandler")
		Node.Outgoing <- r
	}

	log.Println("NodeReader stopped for ", Node.Id)
	Node.Close()
}

// Node sending goroutine - waits for data to be sent over Node.Incoming
// (from IOHandler), then sends it over the socket
func NodeSender(Node *Node) {
	for {
		select {
		case buffer := <-Node.Incoming:
			log.Println("NodeSender sending ", buffer, " to ", Node.Id)
			Node.Conn.Write(buffer)
		case <-Node.Quit:
			log.Println("Node ", Node.Id, " quitting")
			Node.Conn.Close()
			break
		}
	}
}

var nextId int32

func getId() int32 {
	return atomic.AddInt32(&nextId, 1)
}

// Creates a new Node object for each new connection using the Id sent by the Node,
// then starts the NodeSender and NodeReader goroutines to handle the IO
func NodeHandler(conn net.Conn, ch chan dto.Result, NodeList *list.List) {

	newNode := &Node{getId(), make(chan []byte), ch, conn, make(chan bool), NodeList}
	go NodeSender(newNode)
	go NodeReader(newNode)
	NodeList.PushBack(*newNode)
}
