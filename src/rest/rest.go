package rest

import (
	"code.google.com/p/gorest"
	"dto"
	"log"
	"net/http"
)

type InsertChannel struct {
	channel chan dto.Element
}

func (s *InsertChannel) Get() chan dto.Element {
	return s.channel
}

var Channel interface {
	Get() chan dto.Element
} = &InsertChannel{make(chan dto.Element)}

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
	insertData         gorest.EndPoint `method:"POST" path:"/series/id/{id:int32}/data/" postdata:"dto.Element"`
}

func (serv InsertService) InsertData(posted dto.Element, id int32) {
	log.Println("Inserting new data to series: ", id)
	log.Println("Data to insert: ", posted)

	Channel.Get() <- posted

}
