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
	Id          		int32
	Status      		int32
	GpuIds      		[]int32
	Stats       		Info
	stop        		chan bool
	GetTaskChannel 		chan dto.GetTaskRequest
	PreferredDeviceId	int32
}

func NewNode(id int32, connection net.Conn, taskChannel chan dto.GetTaskRequest) *Node {
	n := new(Node)
	n.Status = common.NODE_CONNECTED
	n.Id = id
	n.stop = make(chan bool)
	n.Communication = makeCommunication(connection)
	n.GetTaskChannel = taskChannel
	n.PreferredDeviceId = common.CONST_UNINITIALIZED
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
		if err != nil || len(n.GpuIds) == 0 {
			log.Error("Node login error for node ", n.Id)
			n.Status = common.NODE_ERROR
			return err
		}
	}
	n.PreferredDeviceId = n.GpuIds[0]
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


func (n *Node) processResult(result dto.Result) {

	//Create return channel on which we will get channel to send response
	taskChan := make(chan *dto.Task)

	//Send GetTaskReuest to get channel on which we will return result
	n.GetTaskChannel <- dto.GetTaskRequest{result.TaskId, taskChan}

	//Wait for task
	t := <-taskChan

	//Create rest response from dto.Result
	var responseData []dto.Dto
	var err error
	switch result.Type {
	case common.TASK_INFO:
		var nodeInfo Info
		err = nodeInfo.MemoryInfo.Decode(result.Data)
		nodeInfo.nodeId = n.Id
		log.Debug("Node info %s", nodeInfo.String())
		responseData = []dto.Dto{&nodeInfo}
	case common.TASK_SELECT:
		elementSize := (&dto.Element{}).Size()
		length := len(result.Data) / elementSize
		elements := make([]dto.Dto, length)
		for i := 0; i < length; i++ {
			var e dto.Element
			err = e.Decode(result.Data[i*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			elements[i] = &e
		}
		responseData = elements
	default:
		log.Error("Cannot parse task result data - wrong task type")
		t.ResultChan <- dto.NewRestResponse("Wrong task type", result.TaskId, nil)
		return
	}
	if err != nil || responseData == nil {
		log.Error("Cannot parse task result data", err)
		t.ResultChan <- dto.NewRestResponse("Task result data error", result.TaskId, nil)
	} else {
		t.ResultChan <- dto.NewRestResponse("", result.TaskId, responseData)
	}
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
