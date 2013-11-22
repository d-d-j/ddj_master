package node

import (
	"net"
	"ddj_Master/dto"
)

type Communication struct {
	Incoming 	chan []byte
	Outgoing 	chan dto.Result
	Connection  net.Conn
	Stop     	chan bool
}
