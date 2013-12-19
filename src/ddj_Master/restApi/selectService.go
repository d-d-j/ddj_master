package restApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
)

//Service Definition
type SelectService struct {
	gorest.RestService `root:"/" consumes:"application/json" produces:"application/json"`
	selectQuery        gorest.EndPoint `method:"GET" path:"/data/metric/{metrics:string}/tag/{tags:string}/time/{times:string}/aggregation/{aggr:string}" output:"RestResponse"`
}

func (serv SelectService) SelectQuery(metrics, tags, times, aggr string) RestResponse {
	log.Finest("Selecting data")
	serv.setHeader()
	responseChan := make(chan *RestResponse)
	restRequestChannel <- RestRequest{common.TASK_SELECT_ALL, nil, responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

func prepareQuery(metrics, tags, times, aggr string) (dto.Query, error) {
	return dto.Query{}, nil
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
