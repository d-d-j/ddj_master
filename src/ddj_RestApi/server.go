package ddj_RestApi

import (
	log "code.google.com/p/log4go"
	"code.google.com/p/gorest"
	"net/http"
)

type NetworkApi interface {
	StartApi() <- chan Request
}

type Server struct {
	port string
}

func (sv Server) StartApi() <- chan Request {

	c := make(chan Request)
	insertService := NewInsertService(c)
	selectService := NewSelectService(c)

	if(sv.port == nil) {
		sv.port = 8888
	}

	log.Info("Start REST API on port number " + sv.port)
	gorest.RegisterService(insertService) //Register insert service
	gorest.RegisterService(selectService) //Register select service
	gorest.RegisterMarshaller("application/json", gorest.NewJSONMarshaller())
	go http.Handle("/", gorest.Handle())
	go http.ListenAndServe(sv.port, nil)
	return c
}
