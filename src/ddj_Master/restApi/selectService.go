package restApi

import (
	log "code.google.com/p/log4go"
	"code.google.com/p/gorest"
	"ddj_Master/common"
	"ddj_Master/dto"
)


//Service Definition
type SelectService struct {
	gorest.RestService `root:"/"`
	selectAll       gorest.EndPoint `method:"GET" path:"/data/" output:"[]dto.Dto"`
	getOptions      gorest.EndPoint `method:"OPTIONS" path:"/data/selectOptions"`
	reqChan			chan<- RestRequest
}

func NewSelectService(c chan<- RestRequest) *SelectService {
	ss := new(SelectService)
	ss.reqChan = c
	return ss
}

func (serv SelectService) SelectAll() []dto.Dto {
	serv.setHeader()
	log.Debug("Selecting all data")

	responseChan := make(chan []dto.Dto)
	serv.reqChan <- RestRequest{common.TASK_SELECT_ALL, nil, responseChan}
	response := <-responseChan
	serv.set503HeaderWhenArgumentIsNil(response)
	return response
}

func (serv SelectService) setHeader() {
	rb := serv.ResponseBuilder()
	rb.CacheNoCache()
	rb.AddHeader("Access-Control-Allow-Origin", "*")
	rb.AddHeader("Allow", "GET, POST")
	rb.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept, x-requested-with")
}

func (serv SelectService) set503HeaderWhenArgumentIsNil(arg []dto.Dto) {
	if arg == nil {
		log.Error("Return HTTP 503")
		serv.ResponseBuilder().SetResponseCode(503).WriteAndOveride(
		[]byte("The server is currently unable to handle the request"))
	}
}
