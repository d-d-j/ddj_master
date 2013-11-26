package node

// Imports required packages
import (
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
	"net"
)

// Defines a Node with Id, many GpuIds, a connection object, and
// some channels for sending and receiving data.
type Node struct {
	Communication
	Id       	int32
	GpuIds		[]int32
	Stats		Info
	stop		chan bool
}

func NewNode(id int32, gpuIds []int32, connection net.Conn) *Node {
	n := new(Node)
	n.Id = id
	n.GpuIds = gpuIds
	n.stop = make(chan bool)
	n.Communication = makeCommunication(connection)
	return n
}

func (n *Node) StartWork() {
	go n.readerRoutine()
	go n.senderRoutine()
}

func (n *Node) EndWork(){
	n.stop <- true
}

func (n *Node) readerRoutine() {
	var r dto.Result
	buffer := make([]byte, r.Header.Size())

	for n.Communication.read(buffer) {
		log.Debug("Node reader received data from ", n.Id)

		err := r.Decode(buffer)
		if err != nil {
			log.Error(err)
		}

		log.Fine("Response header: ", r.Header)
		if r.DataSize == 0 {
			r.Data = make([]byte, 0)
		} else {
			r.Data = make([]byte, r.DataSize)
			n.Communication.read(r.Data)
		}
		n.Outgoing <- r
	}

	log.Info("Node reader stopped for Node ", n.Id)
}

// Sending goroutine for Node - waits for data to be sent over Node.Incoming,
// then sends it over the socket
func (n *Node)senderRoutine() {
	for {
		select {
		case buffer := <-n.Communication.Incoming:
			log.Fine("Sending ", buffer, " to ", n.Id)
			n.Communication.write(buffer)
		case <-n.stop:
			log.Info("Node ", n.Id, " stopped")
			n.Communication.close()
			break
		}
	}
}
