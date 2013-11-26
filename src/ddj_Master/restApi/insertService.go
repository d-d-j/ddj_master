package restApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
	"ddj_Master/common"
	"fmt"
)

//Service Definition
type InsertService struct {
	gorest.RestService `root:"/"`
	insertData		gorest.EndPoint `method:"POST" path:"/data/" postdata:"dto.Element"`
	getOptions      gorest.EndPoint `method:"OPTIONS" path:"/data/insertOptions"`
	reqChan			chan<- RestRequest
}

func NewInsertService(c chan<- RestRequest) *InsertService {
	is := new(InsertService)
	is.reqChan = c
	return is
}

func (serv InsertService) InsertData(posted dto.Element) {
	log.Finest("Inserting data - data to insert: ", posted)
	serv.setHeader()
	responseChan := make(chan *RestResponse)
	serv.reqChan <- RestRequest{common.TASK_INSERT, &posted, responseChan}
	response := <-responseChan
	serv.setInsertResponse(response)
}

func (serv InsertService) GetOptions() {
	serv.setHeader()
	log.Debug("Return available options")
}

func (serv InsertService) setHeader() {
	rb := serv.ResponseBuilder()
	rb.CacheNoCache()
	rb.AddHeader("Access-Control-Allow-Origin", "*")
	rb.AddHeader("Allow", "GET, POST")
	rb.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept, x-requested-with")
}

func (serv InsertService) setInsertResponse(response *RestResponse) {
	if response == nil {
		log.Error("Return HTTP 503")
		serv.ResponseBuilder().SetResponseCode(503).WriteAndOveride(
		[]byte("The server is currently unable to handle the request"))
	} else if response.TaskId == 0 {
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride([]byte("Server error - sorry:("))
	} else {
		serv.ResponseBuilder().SetResponseCode(202).WriteAndOveride([]byte(fmt.Sprintf("/data/task/%d/status", response.TaskId)))
	}
}

