package rest

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"constants"
	"dto"
	"net/http"
	"sync/atomic"
)

var nextId int32

func getId() int32 {
	return atomic.AddInt32(&nextId, 1)
}

type InsertChannel struct {
	query chan dto.Query
}

func (s *InsertChannel) QueryChannel() chan dto.Query {
	return s.query
}

var Channel interface {
	QueryChannel() chan dto.Query
} = &InsertChannel{make(chan dto.Query)}

func StartApi(port string) {

	insertService := InsertService{}

	log.Info("Start REST API on " + port)
	gorest.RegisterService(&insertService) //Register our service
	gorest.RegisterMarshaller("application/json", gorest.NewJSONMarshaller())
	go http.Handle("/", gorest.Handle())
	go http.ListenAndServe(port, nil)
}

//Service Definition
type InsertService struct {
	gorest.RestService `root:"/"`
	insertData         gorest.EndPoint `method:"POST" path:"/data/" postdata:"dto.Element"`
	selectAll          gorest.EndPoint `method:"GET" path:"/data/" output:"[]dto.Dto"`
	getOptions         gorest.EndPoint `method:"OPTIONS" path:"/data"`
}

func (serv InsertService) InsertData(posted dto.Element) {
	serv.setHeader()
	log.Debug("Data to insert: ", posted)
	header := dto.TaskRequestHeader{getId(), constants.TASK_INSERT, 0}
	Channel.QueryChannel() <- dto.Query{header, nil, &posted}
}

func (serv InsertService) SelectAll() []dto.Dto {
	serv.setHeader()
	log.Debug("Selecting all data")

	response := make(chan []dto.Dto)
	header := dto.TaskRequestHeader{getId(), constants.TASK_SELECT_ALL, 0}
	Channel.QueryChannel() <- dto.Query{header, response, nil}

	return <-response
}

func (serv InsertService) GetOptions() {
	serv.setHeader()
	log.Debug("Return available options: ")
}

func (serv InsertService) setHeader() {
	rb := serv.ResponseBuilder()
	rb.CacheNoCache()
	rb.AddHeader("Access-Control-Allow-Origin", "*")
	rb.AddHeader("Allow", "GET, POST")
	rb.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept, x-requested-with")
}
