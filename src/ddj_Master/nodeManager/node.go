package nodeManager

// Imports required packages
import (
	log "code.google.com/p/log4go"
	"container/list"
	"net"
	"ddj_Master/dto"
)

// Defines a Node with a Id and connection object, and
// some channels for sending and receiving text.
// Also holds a pointer to the "global" list of all connected Nodes
type Node struct {
	Id       int32
	Incoming chan []byte
	Outgoing chan dto.Result
	Conn     net.Conn
	Quit     chan bool
	NodeList *list.List
}

// Defines a read function for a Node, reading from the connection into
// a buffer passed in. Returns true if read was successful, false otherwise
func (c *Node) Read(buffer []byte) bool {
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
func (c *Node) Close() {
	c.Quit <- true
	c.Conn.Close()
	c.RemoveMe()
}

// Comparison function to easily check equality with another Node
// based on the Id and connection
func (c *Node) Equal(other *Node) bool {
	if c.Id == other.Id {
		if c.Conn == other.Conn {
			return true
		}
	}
	return false
}

// Removes this Node from the Node list
func (c *Node) RemoveMe() {
	for entry := c.NodeList.Front(); entry != nil; entry = entry.Next() {
		Node := entry.Value.(Node)
		if c.Equal(&Node) {
			log.Debug("RemoveMe: ", c.Id)
			c.NodeList.Remove(entry)
		}
	}
}
