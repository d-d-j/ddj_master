package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Info struct {
	NodeId int32
	MemoryInfo
}

type MemoryInfo struct {
	GpuId           int32
	MemoryTotal     int32
	MemoryFree      int32
	GpuMemoryTotal  int32
	GpuMemoryFree   int32
	_               int32
	DBMemoryFree	uint64
}

func (this *MemoryInfo) String() string {
	return fmt.Sprintf("GPUId: %d RAM: %d/%d\tGPU: %d/%d\tDB: %d", this.GpuId, this.MemoryFree, this.MemoryTotal, this.GpuMemoryFree, this.GpuMemoryTotal, this.DBMemoryFree)
}

func (this *Info) String() string {
	return fmt.Sprintf("Node #%d %s\t", this.NodeId, this.MemoryInfo.String())
}

func (this *Info) Size() int {
	return 40
}

func (this *MemoryInfo) Size() int {
	return 32
}

func NewMemoryInfo(gpuId int32, memoryTotal int32, memoryFree int32, gpuMemoryTotal int32, gpuMemoryFree int32, dbMemoryfree uint64) *MemoryInfo {
	memoryInfo := new(MemoryInfo)

	memoryInfo.GpuId = gpuId
	memoryInfo.MemoryTotal = memoryTotal
	memoryInfo.MemoryFree = memoryFree
	memoryInfo.GpuMemoryTotal = gpuMemoryTotal
	memoryInfo.GpuMemoryFree = gpuMemoryFree
	memoryInfo.DBMemoryFree = dbMemoryfree

	return memoryInfo
}

func (this *MemoryInfo) Decode(buf []byte) error {
	buffer := bytes.NewBuffer(buf)

	return binary.Read(buffer, binary.LittleEndian, this)
}

func (this *MemoryInfo) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
