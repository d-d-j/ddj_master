package dto

import (
	"fmt"
)

type RestRequest struct {
	Type     int32
	Data     Dto
	Response chan *RestResponse
}

func (r *RestRequest) String() string {

	return fmt.Sprintf("Request of type: %d", r.Type)
}
