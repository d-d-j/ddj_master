package restApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
	"net/http"
)

var restRequestChannel chan dto.RestRequest = make(chan dto.RestRequest)

type NetworkApi interface {
	StartApi() <-chan dto.RestRequest
}

type Server struct {
	Port string
}

func (sv Server) StartApi() <-chan dto.RestRequest {

	insertService := InsertService{}
	selectService := SelectService{}

	if sv.Port == "" {
		sv.Port = "8888"
	}

	log.Info("Start REST API on port number " + sv.Port)
	gorest.RegisterService(&insertService)
	gorest.RegisterService(&selectService)
	gorest.RegisterMarshaller("application/json", gorest.NewJSONMarshaller())
	go http.Handle("/", gorest.Handle())
	go http.ListenAndServe(sv.Port, nil)
	return restRequestChannel
}
