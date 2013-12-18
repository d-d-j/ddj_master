package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
	"net"
)

type Communication struct {
	Incoming   chan []byte
	Outgoing   chan dto.Result
	connection net.Conn
}

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
	bytesRead, error := c.connection.Read(buffer)
	if error != nil {
		c.connection.Close()
		log.Error("Problem with connection: %s", error)
		return false
	}
	log.Debug("Read %d bytes", bytesRead)
	return true
}

// Defines a write function for a Node, write to the connection from
// a buffer passed in. Returns true if write was successful, false otherwise
func (c *Communication) write(buffer []byte) bool {
	bytesSend, error := c.connection.Write(buffer)
	if error != nil {
		c.connection.Close()
		log.Error("Problem with connection: %s", error)
		return false
	}
	log.Debug("Written %d bytes", bytesSend)
	return true
}

// Closes a connection
func (c *Communication) close() {
	c.connection.Close()
}
