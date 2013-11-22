package node

// Imports required packages
import (
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
)

// Defines a Node with Id, many GpuIds, a connection object, and
// some channels for sending and receiving data.
type Node struct {
	Communication
	Id       	int32
	GpuIds		[]int32
	Stats		Info
	stop		<-chan bool
}

// Comparison function to easily check equality with another Node
// based on the Id and connection
func (n *Node) Equal(other *Node) bool {
	if n.Id == other.Id {
		//if c.Communication == other.Communication {
			return true
		//}
	}
	return false
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
	for n.Read(buffer) {
		log.Debug("Node reader received data from ", n.Id)
		err := r.DecodeHeader(buffer)
		if err != nil {
			log.Error(err)
		}
		log.Fine("Response header: ", r.Header)
		if r.DataSize == 0 {
			r.Data = make([]byte, 0)
			n.Outgoing <- r
			continue
		}
		r.Data = make([]byte, r.DataSize)
		n.Read(r.Data)

		// TODO: MOVE FROM HERE
		/*
		if r.Type == common.TASK_SELECT_ALL {
			length := int(r.LoadSize / 24)
			load := make([]dto.Dto, length)
			for i := 0; i < length; i++ {
				var e dto.Element
				err = e.Decode(buffer[i*(e.Size()+4):])
				if err != nil {
					log.Error(err)
					continue
				}
				load[i] = &e
			}
			r.Load = load
		}
		*/

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
			n.Communication.Write(buffer)
		case <-n.Stop:
			log.Info("Node ", n.Id, " stopped")
			n.Communication.close()
			break
		}
	}
}
