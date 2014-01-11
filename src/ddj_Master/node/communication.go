package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
	"net"
)

//This structure handle network communication with node. It will send data that came on Incoming channel and pass read
//data on Outgoing
type Communication struct {
	Incoming   chan []byte
	Outgoing   chan dto.Result
	connection net.Conn
}

//Constructor for communication. It take net.Conn of established node connection
func NewCommunication(conn net.Conn) *Communication {
	com := new(Communication)
	com.Incoming = make(chan []byte)
	com.Outgoing = make(chan dto.Result)
	com.connection = conn
	return com
}

func makeCommunication(conn net.Conn) Communication {
	in := make(chan []byte)
	out := make(chan dto.Result)
	return Communication{in, out, conn}
}

// Defines a read function for a Node, reading from the connection into
// a buffer passed in. Returns true if read was successful, false otherwise
func (c *Communication) read(buffer []byte) bool {
	allBytesToRead := len(buffer)
	bytesAlreadyRead := 0
	for allBytesToRead != 0 {
		bytesRead, error := c.connection.Read(buffer[bytesAlreadyRead:])
		if error != nil {
			c.connection.Close()
			log.Error("Problem with connection: ", error)
			return false
		}
		allBytesToRead -= bytesRead
		bytesAlreadyRead += bytesRead
		log.Debug("Communication --> Read %d from %d bytes (%d bytes left)", bytesAlreadyRead, len(buffer), allBytesToRead)
	}
	log.Finest("Communication --> Read: %d", buffer)
	return true
}

// Defines a write function for a Node, write to the connection from
// a buffer passed in. Returns true if write was successful, false otherwise
func (c *Communication) write(buffer []byte) bool {
	bytesSend, error := c.connection.Write(buffer)
	if error != nil {
		c.connection.Close()
		log.Error("Problem with connection: ", error)
		return false
	}
	log.Debug("Written %d bytes", bytesSend)
	return true
}

// Closes a connection
func (c *Communication) close() {
	c.connection.Close()
	close(c.Outgoing)
}
