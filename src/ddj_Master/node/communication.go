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
