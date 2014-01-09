package restApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"fmt"
)

//Service Definition
type InsertService struct {
	gorest.RestService `root:"/" consumes:"application/json" produces:"application/json"`
	insertData         gorest.EndPoint `method:"POST" path:"/data/" postdata:"ddj_Master.dto.Element"`
	flushBuffer        gorest.EndPoint `method:"POST" path:"/data/flush" postdata:"string"`
	getOptions         gorest.EndPoint `method:"OPTIONS" path:"/data/{...:string}"`
}

func (serv InsertService) GetOptions(_ ...string) {
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

func (serv InsertService) setInsertResponse(response *dto.RestResponse) {
	if response == nil {
		log.Error("Return HTTP 503")
		serv.ResponseBuilder().SetResponseCode(503).WriteAndOveride(
			[]byte("The server is currently unable to handle the request"))
	} else if response.TaskId == 0 {
		log.Error("Return HTTP 500")
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride(
			[]byte("Server error - sorry:("))
	} else if response.Error != "" {
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride(
			[]byte(fmt.Sprintf("TaskId: %d, Error: %s", response.TaskId, response.Error)))
	} else {
		log.Finest("Return URI to the result")
		serv.ResponseBuilder().SetResponseCode(202).WriteAndOveride(
			[]byte(fmt.Sprintf("/data/task/%d/status", response.TaskId)))
	}
}

func (serv InsertService) InsertData(PostData dto.Element) {
	serv.setHeader()
	log.Finest("Inserting data - data to insert: ", PostData)
	responseChan := make(chan *dto.RestResponse)
	restRequestChannel <- dto.RestRequest{Type: common.TASK_INSERT, Data: &PostData, Response: responseChan}
	response := <-responseChan
	log.Finest("Insert Result: ", response)
	serv.setInsertResponse(response)
}

func (serv InsertService) FlushBuffer(_ string) {
	serv.setHeader()
	log.Finest("Flush GPU buffer")
	responseChan := make(chan *dto.RestResponse)
	restRequestChannel <- dto.RestRequest{Type: common.TASK_FLUSH, Data: new(dto.EmptyElement), Response: responseChan}
	response := <-responseChan
	log.Finest("Flush Result: ", response)
	serv.setInsertResponse(response)
}
