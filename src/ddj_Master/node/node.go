package node

// Imports required packages
import (
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
	"net"
	"ddj_Master/common"
	"encoding/binary"
	"bytes"
)

// Defines a Node with Id, many GpuIds, a connection object, and
// some channels for sending and receiving data.
type Node struct {
	Communication
	Id       	int32
	Status		int32
	GpuIds		[]int32
	Stats		Info
	stop		chan bool
}

func NewNode(id int32, connection net.Conn) *Node {
	n := new(Node)
	n.Status = common.NODE_CONNECTED
	n.Id = id
	n.stop = make(chan bool)
	n.Communication = makeCommunication(connection)
	return n
}

func (n *Node) StartWork(balancerChannel chan<- Info) {
	err := n.waitForLogin()
	if err == nil {
		go n.readerRoutine()
		go n.senderRoutine()
	}
	balancerChannel <- Info{n.Id}
	log.Debug("Node %d READY", n.Id)
}

func (n *Node) EndWork(){
	n.stop <- true
}

func (n *Node) waitForLogin() error  {
	buffer := make([]byte, 4)
	n.Communication.read(buffer)
	buf := bytes.NewBuffer(buffer)
	var cudaGpuCount int32
	err := binary.Read(buf, binary.LittleEndian, &cudaGpuCount)
	if err != nil {
		log.Error("Node login error for node ", n.Id)
		n.Status = common.NODE_ERROR
		return err
	}
	log.Debug("Node received login data (CUDA GPU COUNT = ", cudaGpuCount, ") from node ", n.Id)

	buffer = make([]byte, 4*cudaGpuCount)
	n.Communication.read(buffer)
	buf = bytes.NewBuffer(buffer)
	n.GpuIds = make([]int32, cudaGpuCount)
	for i:=int32(0); i<cudaGpuCount; i++ {
		err = binary.Read(buf, binary.LittleEndian, &(n.GpuIds[i]))
		if err != nil {
			log.Error("Node login error for node ", n.Id)
			n.Status = common.NODE_ERROR
			return err
		}
	}
	log.Debug("Node ", n.Id, " is ready with devices ", n.GpuIds)
	n.Status = common.NODE_READY
	return err
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
