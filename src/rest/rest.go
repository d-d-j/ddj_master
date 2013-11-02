package rest

import (
	"code.google.com/p/gorest"
	"constants"
	"dto"
	"log"
	"net/http"
	"sync/atomic"
)

var nextId int32

func getId() int32 {
	return atomic.AddInt32(&nextId, 1)
}

type InsertChannel struct {
	insert chan dto.Query
	query  chan dto.Query
}

func (s *InsertChannel) Get() chan dto.Query {
	return s.insert
}

func (s *InsertChannel) QueryChannel() chan dto.Query {
	return s.query
}

var Channel interface {
	Get() chan dto.Query
	QueryChannel() chan dto.Query
} = &InsertChannel{make(chan dto.Query), make(chan dto.Query)}

func StartApi(port string) {

	insertService := InsertService{}

	log.Print("Start REST API on " + port)
	gorest.RegisterService(&insertService) //Register our service
	gorest.RegisterMarshaller("application/json", gorest.NewJSONMarshaller())
	go http.Handle("/", gorest.Handle())
	go http.ListenAndServe(port, nil)
}

//Service Definition
type InsertService struct {
	gorest.RestService `root:"/"`
	insertData         gorest.EndPoint `method:"POST" path:"/data/" postdata:"dto.Element"`
	selectAll          gorest.EndPoint `method:"GET" path:"/data/" output:"string"`
}

func (serv InsertService) InsertData(posted dto.Element) {
	log.Println("Data to insert: ", posted)
	header := dto.TaskRequestHeader{getId(), constants.TASK_INSERT, 0}
	Channel.Get() <- dto.Query{header, nil, &posted}

}

func (serv InsertService) SelectAll() string {
	log.Println("Selecting all data")

	response := make(chan string)
	//TODO add ID
	header := dto.TaskRequestHeader{getId(), constants.TASK_SELECT_ALL, 0}
	Channel.QueryChannel() <- dto.Query{header, response, nil}

	return <-response
}
