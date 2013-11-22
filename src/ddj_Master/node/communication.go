package node

import (
	"net"
	"ddj_Master/dto"
	log "code.google.com/p/log4go"
)

type Communication struct {
	Incoming 	chan []byte
	Outgoing 	chan dto.Result
	Connection  net.Conn
	Stop     	chan bool
}

// Defines a read function for a Node, reading from the connection into
// a buffer passed in. Returns true if read was successful, false otherwise
func (c *Communication) Read(buffer []byte) bool {
	log.Debug(c.Id, " try to read ", len(buffer), " bytes")
	bytesRead, error := c.Conn.Read(buffer)
	if error != nil {
		c.Close()
		log.Error("Problem with connection: ", error)
		return false
	}
	log.Debug("Read ", bytesRead, " bytes")
	return true
}

// Closes a Node connection and removes it from the Node list
func (c *Communication) Close() {
	c.Stop <- true
	c.Conn.Close()
}

func NodeReader(Node *Node) {

	var r dto.Result
	buffer := make([]byte, r.TaskRequestHeader.Size())
	for Node.Read(buffer) {
		log.Debug("NodeReader received data from", Node.Id)
		err := r.DecodeHeader(buffer)
		if err != nil {
			log.Error(err)
		}
		log.Fine("Response header: ", r.TaskRequestHeader)
		if r.LoadSize == 0 {
			r.Load = make([]dto.Dto, 0)
			Node.Outgoing <- r
			continue
		}
		buffer := make([]byte, r.LoadSize)
		Node.Read(buffer)
		if r.Code == constants.TASK_SELECT_ALL {
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
		log.Debug("Send response to IOHandler")
		Node.Outgoing <- r
		buffer = make([]byte, r.TaskRequestHeader.Size())
	}

	log.Info("NodeReader stopped for ", Node.Id)
}

// Node sending goroutine - waits for data to be sent over Node.Incoming
// (from IOHandler), then sends it over the socket
func NodeSender(Node *Node) {
	for {
		select {
		case buffer := <-Node.Incoming:
			log.Fine("NodeSender sending ", buffer, " to ", Node.Id)
			Node.Conn.Write(buffer)
		case <-Node.Quit:
			log.Info("Node ", Node.Id, " quitting")
			Node.Conn.Close()
			break
		}
	}
}
