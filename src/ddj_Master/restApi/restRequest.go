package restApi

import (
	"fmt"
	"ddj_Master/dto"
)

type RestRequest struct {
	Type     int32
	Data	 dto.Dto
	Response chan *RestResponse
}

func (r *RestRequest) String() string {

	return fmt.Sprintf("Request of type: %d", r.Type)
}
