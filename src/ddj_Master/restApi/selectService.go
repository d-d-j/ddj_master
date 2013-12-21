package restApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"strconv"
	"strings"
)

//Service Definition
type SelectService struct {
	gorest.RestService `root:"/" consumes:"application/json" produces:"application/json"`
	selectQuery        gorest.EndPoint `method:"GET" path:"/data/metric/{metrics:string}/tag/{tags:string}/time/{times:string}/aggregation/{aggr:string}" output:"RestResponse"`
}

func (serv SelectService) SelectQuery(metrics, tags, times, aggr string) RestResponse {
	log.Finest("Selecting data")
	serv.setHeader()
	responseChan := make(chan *RestResponse)
	data, err := prepareQuery(metrics, tags, times, aggr)
	if err != nil {
		log.Error("Return HTTP 400")
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(
			[]byte("The request could not be understood by the server due to malformed syntax. You SHOULD NOT repeat the request without modifications."))
		return RestResponse{}
	}
	log.Fine("Query: ", &data)
	restRequestChannel <- RestRequest{common.TASK_SELECT_ALL, &data, responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

func prepareQuery(metrics, tags, times, aggr string) (dto.Query, error) {

	const ALL string = "all"

	timesSplited := strings.Split(times, ",")

	timesArr := make([]int64, len(timesSplited)*2)

	var (
		metricsArr []int32
		tagsArr    []int32
	)
	if metrics == ALL {
		metricsArr = make([]int32, 0)
	} else {
		metricsSplited := strings.Split(metrics, ",")
		metricsArr = make([]int32, len(metricsSplited))
		for i, element := range metricsSplited {
			value, err := strconv.Atoi(element)
			if err != nil {
				return dto.Query{}, err
			}
			metricsArr[i] = int32(value)
		}
	}

	if tags == ALL {
		tagsArr = make([]int32, 0)
	} else {
		tagsSplited := strings.Split(tags, ",")
		tagsArr = make([]int32, len(tagsSplited))
		for i, element := range tagsSplited {
			value, err := strconv.Atoi(element)
			if err != nil {
				return dto.Query{}, err
			}
			tagsArr[i] = int32(value)
		}
	}

	for i := 0; i < len(timesSplited); i++ {
		time := strings.Split(timesSplited[i], "-")
		value, err := strconv.ParseInt(time[0], 10, 64)
		if err != nil {
			return dto.Query{}, err
		}
		timesArr[2*i] = int64(value)
		value, err = strconv.ParseInt(time[1], 10, 64)
		if err != nil {
			return dto.Query{}, err
		}
		timesArr[2*i+1] = int64(value)
	}

	var aggregation int32
	switch aggr {
	case "none":
		aggregation = common.AGGREGATION_NONE
	case "sum":
		aggregation = common.AGGREGATION_ADD
	case "average":
		aggregation = common.AGGREGATION_AVERAGE
	}

	return dto.Query{int32(len(metricsArr)), metricsArr, int32(len(tagsArr)), tagsArr, int32(len(timesSplited)), timesArr, aggregation}, nil
}

func (serv SelectService) setHeader() {
	rb := serv.ResponseBuilder()
	rb.CacheNoCache()
	rb.AddHeader("Access-Control-Allow-Origin", "*")
	rb.AddHeader("Allow", "GET, POST")
	rb.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept, x-requested-with")
}

func (serv SelectService) setSelectHeaderErrors(response *RestResponse) {
	if response == nil {
		log.Error("Return HTTP 503")
		serv.ResponseBuilder().SetResponseCode(503).WriteAndOveride(
			[]byte("The server is currently unable to handle the request"))
	} // TODO: Set more errors if response.Error != "" or TaskId == 0
}
