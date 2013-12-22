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
	GpuMemoryTotal int32
	GpuMemoryFree  int32
}

func (this *MemoryInfo) String() string {
	return fmt.Sprintf("RAM: %d/%d\tGPU: %d/%d", this.MemoryFree, this.MemoryTotal, this.GpuMemoryFree, this.GpuMemoryTotal)
}

func (this *Info) String() string {
	return fmt.Sprintf("Node #%d\t", this.nodeId, this.MemoryInfo.String())
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
