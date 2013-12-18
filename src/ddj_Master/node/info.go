package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Info struct {
	nodeId int32
	MemoryInfo
}

type MemoryInfo struct {
	MemoryTotal    int32
	MemoryFree     int32
	GpuMemoryFree  int32
	GpuMemoryTotal int32
}

func (this *Info) String() string {
	return fmt.Sprintf("#%d", this.nodeId)
}

func (this *Info) Size() int {
	return 24
}

func (this *MemoryInfo) Decode(buf []byte) error {
	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, this)

}

func (this *MemoryInfo) Encode() ([]byte, error) {
	return nil, fmt.Errorf("Not implemented yet!")
}
