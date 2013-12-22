package restApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"fmt"
	"strconv"
	"strings"
)

//Service Definition
type SelectService struct {
	gorest.RestService `root:"/" consumes:"application/json" produces:"application/json"`
	selectQuery        gorest.EndPoint `method:"GET" path:"/data/metric/{metrics:string}/tag/{tags:string}/time/{times:string}/aggregation/{aggr:string}" output:"RestResponse"`
}

func (serv SelectService) SelectQuery(metrics, tags, times, aggr string) dto.RestResponse {
	log.Finest("Selecting data")
	serv.setHeader()
	responseChan := make(chan *dto.RestResponse)
	data, err := prepareQuery(metrics, tags, times, aggr)
	if err != nil {
		log.Error("Return HTTP 400")
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(
			[]byte("The request could not be understood by the server due to malformed syntax. You SHOULD NOT repeat the request without modifications."))
		return dto.RestResponse{}
	}
	log.Fine("Query: ", &data)
	restRequestChannel <- dto.RestRequest{common.TASK_SELECT_ALL, &data, responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

func prepareQuery(metrics, tags, times, aggr string) (dto.Query, error) {

	metricsArr, err := prepareTagsOrMetrics(metrics)
	if err != nil {
		return dto.Query{}, err
	}
	tagsArr, err := prepareTagsOrMetrics(tags)
	if err != nil {
		return dto.Query{}, err
	}
	timesArr, err := prepareTimeSpans(times)
	if err != nil {
		return dto.Query{}, err
	}

	aggregation, err := prepareAggregationType(aggr)
	if err != nil {
		return dto.Query{}, err
	}

	return dto.Query{int32(len(metricsArr)), metricsArr, int32(len(tagsArr)), tagsArr, int32(len(timesArr) / 2), timesArr, aggregation}, nil
}

func prepareAggregationType(aggregation string) (int32, error) {
	switch aggregation {
	case "none":
		return common.AGGREGATION_NONE, nil
	case "sum":
		return common.AGGREGATION_ADD, nil
	case "average":
		return common.AGGREGATION_AVERAGE, nil
	default:
		return common.TASK_ERROR, fmt.Errorf("Unknown aggregation type")
	}
}

func prepareTimeSpans(times string) ([]int64, error) {
	timesSplited := strings.Split(times, ",")
	timesArr := make([]int64, len(timesSplited)*2)
	for i := 0; i < len(timesSplited); i++ {
		time := strings.Split(timesSplited[i], "-")
		value, err := strconv.ParseInt(time[0], 10, 64)
		if err != nil {
			return nil, err
		}
		timesArr[2*i] = int64(value)
		value, err = strconv.ParseInt(time[1], 10, 64)
		if err != nil {
			return nil, err
		}
		timesArr[2*i+1] = int64(value)
	}
	return timesArr, nil
}

func prepareTagsOrMetrics(input string) ([]int32, error) {
	const ALL string = "all"
	if input == ALL {
		return make([]int32, 0), nil
	}

	inputSplited := strings.Split(input, ",")
	inputArr := make([]int32, len(inputSplited))
	for i, element := range inputSplited {
		value, err := strconv.Atoi(element)
		if err != nil {
			return nil, err
		}
		inputArr[i] = int32(value)
	}

	return inputArr, nil
}

func (serv SelectService) setHeader() {
	rb := serv.ResponseBuilder()
	rb.CacheNoCache()
	rb.AddHeader("Access-Control-Allow-Origin", "*")
	rb.AddHeader("Allow", "GET, POST")
	rb.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept, x-requested-with")
}

func (serv SelectService) setSelectHeaderErrors(response *dto.RestResponse) {
	if response == nil {
		log.Error("Return HTTP 503")
		serv.ResponseBuilder().SetResponseCode(503).WriteAndOveride(
			[]byte("The server is currently unable to handle the request"))
	} // TODO: Set more errors if response.Error != "" or TaskId == 0
}
