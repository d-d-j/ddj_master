package rest

import (
	"code.google.com/p/gorest"
	"constants"
	"dto"
	"log"
	"net/http"
)

type InsertChannel struct {
	insert chan dto.Element
	query  chan dto.Query
}

func (s *InsertChannel) Get() chan dto.Element {
	return s.insert
}

func (s *InsertChannel) QueryChannel() chan dto.Query {
	return s.query
}

var Channel interface {
	Get() chan dto.Element
	QueryChannel() chan dto.Query
} = &InsertChannel{make(chan dto.Element), make(chan dto.Query)}

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

	Channel.Get() <- posted

}

func (serv InsertService) SelectAll() string {
	log.Println("Selecting all data")

	response := make(chan string)
	//TODO add ID
	Channel.QueryChannel() <- dto.Query{1, constants.TASK_SELECT_ALL, response}

	return <-response
}
