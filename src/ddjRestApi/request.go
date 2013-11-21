package ddjRestApi

import (
	"fmt"
	"dto"
)

type Request struct {
	Type     int32
	Data	 dto.Dto
	Response <-chan dto.Dto
}

func (r *Request) String() string {

	return fmt.Sprintf("Request of type: %d, with data size: %d", r.Type, r.DataSize)
}
