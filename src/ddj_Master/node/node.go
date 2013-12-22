package node

// Imports required packages
import (
	"bytes"
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"encoding/binary"
	"net"
)

// Defines a Node with Id, many GpuIds, a connection object, and
// some channels for sending and receiving data.
type Node struct {
	Communication
	Id          int32
	Status      int32
	GpuIds      []int32
	Stats       Info
	stop        chan bool
	TaskChannel chan dto.GetTaskRequest
}

func NewNode(id int32, connection net.Conn, taskChannel chan dto.GetTaskRequest) *Node {
	n := new(Node)
	n.Status = common.NODE_CONNECTED
	n.Id = id
	n.stop = make(chan bool)
	n.Communication = makeCommunication(connection)
	n.TaskChannel = taskChannel
	return n
}

func (n *Node) StartWork(balancerChannel chan<- Info) {
	err := n.waitForLogin()
	if err == nil {
		go n.readerRoutine()
		go n.senderRoutine()
	}
	balancerChannel <- Info{n.Id, MemoryInfo{1, 1, 1, 1}}
	log.Debug("Node %d READY", n.Id)
}

func (n *Node) EndWork() {
	n.stop <- true
}

func (n *Node) waitForLogin() error {
	buffer := make([]byte, 4)
	if !n.read(buffer) {
		n.stop <- true
	}
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
	if !n.read(buffer) {
		n.stop <- true
	}
	buf = bytes.NewBuffer(buffer)
	n.GpuIds = make([]int32, cudaGpuCount)
	for i := int32(0); i < cudaGpuCount; i++ {
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

	for n.read(buffer) {
		log.Debug("Node reader received data from #%d", n.Id)

		err := r.Decode(buffer)
		if err != nil {
			log.Error(err)
		}

		log.Fine("Response header: %s", r.Header.String())
		if r.DataSize == 0 {
			r.Data = make([]byte, 0)
		} else {
			log.Finest("Reading response data %d bytes", r.DataSize)
			r.Data = make([]byte, r.DataSize)
			n.read(r.Data)
		}
		go n.processResult(r)
	}

	log.Info("Node reader stopped for Node ", n.Id)
	n.stop <- true
}

//FIXME
func (n *Node) processResult(result dto.Result) {
	taskChan := make(chan chan *dto.RestResponse)
	n.TaskChannel <- dto.GetTaskRequest{result.TaskId, taskChan}
	r := <-taskChan
	var nodeInfo Info
	err := nodeInfo.MemoryInfo.Decode(result.Data)
	if err != nil {
		log.Error("Cannot parse node info ", err)
	}
	nodeInfo.nodeId = n.Id
	log.Debug("Node info %s", nodeInfo.String())
	log.Finest(r)
	r <- dto.NewRestResponse("", result.TaskId, []dto.Dto{&nodeInfo})
}

// Sending goroutine for Node - waits for data to be sent over Node.Incoming,
// then sends it over the socket
func (n *Node) senderRoutine() {
	for {
		select {
		case buffer := <-n.Incoming:
			log.Fine("Sending ", buffer, " to ", n.Id)
			n.write(buffer)
		case <-n.stop:
			log.Info("Node #%d stopped", n.Id)
			n.close()
			NodeManager.DelChan <- n.Id
			break
		}
	}
}
