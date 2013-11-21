package ddj_RestApi

import (
	"ddj_Dto"
	"fmt"
	"dto"
)

type Request struct {
	Type     int32
	Data	 ddj_Dto.Dto
	Response <-chan ddj_Dto.Dto
}

func (r *Request) String() string {

	return fmt.Sprintf("Request of type: %d, with data size: %d", r.Type, r.DataSize)
}
