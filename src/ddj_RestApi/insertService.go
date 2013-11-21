package ddj_RestApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"ddj_Dto"
	"ddj"
)

//Service Definition
type InsertService struct {
	gorest.RestService `root:"/"`
	insertData		gorest.EndPoint `method:"POST" path:"/data/" postdata:"dto.Element"`
	getOptions      gorest.EndPoint `method:"OPTIONS" path:"/data/insertOptions"`
	reqChan			chan<- Request
}

func NewInsertService(c <-chan Request) *InsertService {

	is := new(InsertService)
	is.reqChan = c
	return is
}

func (serv InsertService) InsertData(posted ddj_Dto.Element) {

	serv.setHeader()
	log.Debug("Data to insert: ", posted)
	responseChan := make(chan []ddj_Dto.Dto)
	serv.reqChan <- Request{ddj.TASK_INSERT, &posted, responseChan}
	response := <-responseChan
	serv.set503HeaderWhenArgumentIsNil(response)
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

func (serv InsertService) set503HeaderWhenArgumentIsNil(arg []ddj_Dto.Dto) {

	if arg == nil {
		log.Error("Return HTTP 503")
		serv.ResponseBuilder().SetResponseCode(503).WriteAndOveride(
		[]byte("The server is currently unable to handle the request"))
	}
}

