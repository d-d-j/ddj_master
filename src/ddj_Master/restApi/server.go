package restApi

import (
	log "code.google.com/p/log4go"
	"code.google.com/p/gorest"
	"net/http"
)

type NetworkApi interface {
	StartApi() <- chan RestRequest
}

type Server struct {
	Port string
}

func (sv Server) StartApi() <-chan RestRequest {

	c := make(chan RestRequest)
	insertService := InsertService{}
	insertService.reqChan = c
	selectService := SelectService{}
	selectService.reqChan = c

	if sv.Port == "" {
		sv.Port = "8888"
	}

	log.Info("Start REST API on port number " + sv.Port)
	gorest.RegisterService(&insertService)
	gorest.RegisterService(&selectService)
	gorest.RegisterMarshaller("application/json", gorest.NewJSONMarshaller())
	go http.Handle("/", gorest.Handle())
	go http.ListenAndServe(sv.Port, nil)
	return c
}
