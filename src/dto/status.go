package dto

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Status struct {
	totalRam, usedRam int32
	gpuTemperature    int32
}

func NewStatus(totalRam int32, usedRam int32, gpuTemp int32) *Status {
	s := new(Status)
	s.totalRam = totalRam
	s.usedRam = usedRam
	s.gpuTemperature = gpuTemp
	return s
}

func (s *Status) Equal(other *Status) bool {
	left, err := s.GobEncode()
	if err != nil {
		return false
	}
	right, err := s.GobEncode()
	if err != nil {
		return false
	}
	if bytes.Equal(left, right) {
		return true
	}
	return false
}

func (s *Status) String() string {
	return fmt.Sprintf("Ram: %d/%d Temp: %d℃", s.usedRam, s.totalRam, s.gpuTemperature)
}

func (s *Status) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(s.totalRam)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(s.usedRam)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(s.gpuTemperature)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (s *Status) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&s.totalRam)
	if err != nil {
		return err
	}
	err = decoder.Decode(&s.usedRam)
	if err != nil {
		return err
	}
	return decoder.Decode(&s.gpuTemperature)
}
