package node

import (
	"net"
	"ddj_Master/dto"
	log "code.google.com/p/log4go"
)

type Communication struct {
	Incoming 	chan []byte
	Outgoing 	chan dto.Result
	connection  net.Conn
}

func NewCommunication(conn net.Conn) *Communication {
	in := make(chan []byte)
	out := make(chan dto.Result)
	com := new(Communication)
	com.Incoming = in
	com.Outgoing = out
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
	log.Debug(c.Id, " try to read ", len(buffer), " bytes")
	bytesRead, error := c.connection.Read(buffer)
	if error != nil {
		c.Close()
		log.Error("Problem with connection: ", error)
		return false
	}
	log.Debug("Read ", bytesRead, " bytes")
	return true
}

// Defines a write function for a Node, write to the connection from
// a buffer passed in. Returns true if write was successful, false otherwise
func (c *Communication) send(buffer []byte) bool {
	bytesSend, error := c.connection.Write(buffer)
	if error != nil {
		c.Close()
		log.Error("Problem with connection: ", error)
		return false
	}
	log.Debug("Written ", bytesSend, " bytes")
	return true
}

// Closes a connection
func (c *Communication) close() {
	c.connection.Close()
}
