package restApi

import (
	log "code.google.com/p/log4go"
	"code.google.com/p/gorest"
	"ddj_Master/common"
)


//Service Definition
type SelectService struct {
	gorest.RestService `root:"/ddj/" consumes:"application/json" produces:"application/json"`
	selectAll       gorest.EndPoint `method:"GET" path:"/data/" output:"RestResponse"`
	reqChan			chan<- RestRequest
}

func (serv SelectService) SelectAll() RestResponse {
	log.Finest("Selecting all data")
	serv.setHeader()
	responseChan := make(chan *RestResponse)
	serv.reqChan <- RestRequest{common.TASK_SELECT_ALL, nil, responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

func (serv SelectService) setHeader() {
	rb := serv.ResponseBuilder()
	rb.CacheNoCache()
	rb.AddHeader("Access-Control-Allow-Origin", "*")
	rb.AddHeader("Allow", "GET, POST")
	rb.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept, x-requested-with")
}

func (serv SelectService) setSelectHeaderErrors(response *RestResponse) {
	if response == nil {
		log.Error("Return HTTP 503")
		serv.ResponseBuilder().SetResponseCode(503).WriteAndOveride(
		[]byte("The server is currently unable to handle the request"))
	} // TODO: Set more errors if response.Error != "" or TaskId == 0
}
