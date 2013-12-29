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
	GpuId			 int32
	MemoryTotal      int32
	MemoryFree       int32
	GpuMemoryTotal   int32
	GpuMemoryFree    int32
}

func (this *MemoryInfo) String() string {
	return fmt.Sprintf("GPUId: %d RAM: %d/%d\tGPU: %d/%d", this.GpuId, this.MemoryFree, this.MemoryTotal, this.GpuMemoryFree, this.GpuMemoryTotal)
}

func (this *Info) String() string {
	return fmt.Sprintf("Node #%d %s\t", this.nodeId, this.MemoryInfo.String())
}

func (this *Info) Size() int {
	return 28
}

func (this *MemoryInfo) Size() int {
	return 20
}

func (this *MemoryInfo) Decode(buf []byte) error {
	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, this)

}

func (this *MemoryInfo) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, this.GpuId)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, this.MemoryTotal)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, this.MemoryFree)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, this.GpuMemoryTotal)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, this.GpuMemoryFree)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
