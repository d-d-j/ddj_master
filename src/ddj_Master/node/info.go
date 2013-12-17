package node

import (
	"fmt"
)

type Info struct {
	nodeId int32
}

func (this Info) String() string {
	return fmt.Sprintf("#%d", this.nodeId)
}
