//This package contains definition of structures that are adapter for node. It hides all network traffic and
//expose data via channels
package node

// Imports required packages
import (
	"bytes"
	log "code.google.com/p/log4go"
	"github.com/d-d-j/ddj_master/common"
	"github.com/d-d-j/ddj_master/dto"
	"encoding/binary"
	"net"
	"time"
)

// This is main structure that define node abstraction in master
type Node struct {
	Communication
	Id                int32
	Status            int32
	GpuIds            []int32
	stop              chan bool
	GetTaskChannel    chan dto.GetTaskRequest
	PreferredDeviceId int32
}

//Node constructor. It takes new node id, it's connection and channel to get task from node Manager
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

//This method wait for node's  login message and after login it start communication with this node.
func (n *Node) StartWork(balancerChannel chan<- []*dto.Info) {
	err := n.waitForLogin()
	if err == nil {
		go n.readerRoutine()
		go n.senderRoutine()
	}
	balancerChannel <- []*dto.Info{&dto.Info{NodeId: n.Id, MemoryInfo: dto.MemoryInfo{GpuId: n.GpuIds[0], MemoryTotal: 1, MemoryFree: 1, GpuMemoryTotal: 1, GpuMemoryFree: 1, DBMemoryFree: 1}}}
	log.Debug("Node %d READY", n.Id)
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
	log.Debug("Node received login data (CUDA GPU COUNT = %d) from node %d", cudaGpuCount, n.Id)

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
	log.Debug("Node %d is ready with devices %v", n.Id, n.GpuIds)
	n.Status = common.NODE_READY
	return err
}

func (n *Node) readerRoutine() {

	buffer := make([]byte, (&dto.Header{}).Size())
	for n.read(buffer) {
		log.Debug("Node reader received data from #%d", n.Id)
		r := &dto.Result{NodeId: n.Id}

		err := r.Decode(buffer)
		if err != nil {
			log.Error(err)
		}

		log.Fine("Response header: %s", r.Header.String())
		if r.DataSize == 0 {
			log.Finest("No additional data for %d", r.TaskId)
			r.Data = make([]byte, 0)
		} else {
			log.Finest("Reading response data %d bytes for task %d", r.DataSize, r.TaskId)
			r.Data = make([]byte, r.DataSize)
			n.read(r.Data)
			log.Finest("Response data of %d bytes read SUCCESS", r.DataSize)
		}

		log.Fine("Node is sending data to worker (taskId: %d)", r.TaskId)
		taskChan := make(chan *dto.Task)
		n.GetTaskChannel <- dto.GetTaskRequest{TaskId: r.TaskId, BackChan: taskChan}
		t := <-taskChan
		if t == nil {
			log.Critical(`Task %d does not exist.
				Some data could be lost.
				Omitting response from node %d because task %d is nil.
				Can't send data to worker`, r.TaskId, n.Id, r.TaskId)
			panic("Task not exist")
		}
		timeout := time.After(1 * time.Second)
		select {
		case t.ResultChan <- r:
			log.Fine("Task result sent to worker SUCCESS (task: ", t, "\"")
		case <-timeout:
			panic("TIMEOUT - sending task result to worker")
		}

	}

	log.Info("Node reader stopped for Node %s", n.Id)
	n.stop <- true
}

// Sending goroutine for Node - waits for data to be sent over Node.Incoming,
// then sends it over the socket
func (n *Node) senderRoutine() {
	for {
		select {
		case buffer := <-n.Incoming:
			log.Fine("Sending %v  to %d", buffer, n.Id)
			n.write(buffer)
		case <-n.stop:
			log.Info("Node #%d stopped", n.Id)
			NodeManager.DelChan <- n.Id
			n.close()
			break
		}
	}
}
