package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	TaskId   int64
	Type     int32
	DataSize int32
	DeviceId int32
}

func (h *Header) String() string {
	return fmt.Sprintf("#%d Code: %d [%d] - dev#%d", h.TaskId, h.Type, h.DataSize, h.DeviceId)
}

func (h *Header) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, h)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *Header) Decode(buf []byte) error {
	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, h)

}

func (h *Header) Size() int {
	return binary.Size(h)
}
