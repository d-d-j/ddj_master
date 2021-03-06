package restApi

import (
	"code.google.com/p/gorest"
	log "code.google.com/p/log4go"
	"github.com/d-d-j/ddj_master/common"
	"github.com/d-d-j/ddj_master/dto"
	"fmt"
	"strconv"
	"strings"
        "sort"
)

//Select Service definition
type SelectService struct {
	gorest.RestService    `root:"/" consumes:"application/json" produces:"application/json"`
	selectQuery           gorest.EndPoint `method:"GET" path:"/data/metric/{metrics:string}/tag/{tags:string}/time/{times:string}/aggregation/{aggr:string}" output:"RestResponse"`
	histogramByValueQuery gorest.EndPoint `method:"GET" path:"/data/metric/{metrics:string}/tag/{tags:string}/time/{times:string}/aggregation/histogramByValue/from/{from:float32}/to/{to:float32}/buckets/{buckets:int32}" output:"RestResponse"`
	histogramByTimeQuery  gorest.EndPoint `method:"GET" path:"/data/metric/{metrics:string}/tag/{tags:string}/time/{times:string}/aggregation/histogramByTime/from/{from:int64}/to/{to:int64}/buckets/{buckets:int32}" output:"RestResponse"`
	interpolateQuery      gorest.EndPoint `method:"GET" path:"/data/metric/{metrics:string}/tag/{tags:string}/time/from/{from:int64}/to/{to:int64}/aggregation/series/sum/samples/{samples:int32}" output:"RestResponse"`
}

//This method handle  standard select query
func (serv SelectService) SelectQuery(metrics, tags, times, aggr string) dto.RestResponse {
	log.Finest("Selecting data")
	serv.setHeader()
	responseChan := make(chan *dto.RestResponse)
	data, err := prepareQuery(metrics, tags, times, aggr)
	if err != nil {
		log.Error("Return HTTP 400", err)
		msg := fmt.Sprintf("The request could not be understood by the server due to malformed syntax. You SHOULD NOT repeat the request without modifications.\n%s", err)
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(
			[]byte(msg))
		return dto.RestResponse{}
	}
	log.Fine("Query: ", &data)
	restRequestChannel <- dto.RestRequest{Type: common.TASK_SELECT, Data: &data, Response: responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

//This method handle query to aggregate with series interpolation
func (serv SelectService) InterpolateQuery(metrics, tags string, from, to int64, samples int32) dto.RestResponse {
	log.Finest("Selecting Series")
	serv.setHeader()
	responseChan := make(chan *dto.RestResponse)
	times := fmt.Sprintf("%d-%d", from, to)
	query, err := prepareQueryWithoutAggregationType(metrics, tags, times)
	if err != nil || query.MetricsCount == 0 || query.TagsCount == 0 {
		log.Error("Return HTTP 400")
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(
			[]byte("The request could not be understood by the server due to malformed syntax. You SHOULD NOT repeat the request without modifications."))
		return dto.RestResponse{}
	}
	query.AggregationType = common.AGGREGATION_SERIES_SUM
	query.AdditionalData = dto.InterpolatedData(samples)
	log.Fine("Query: ", &query)
	restRequestChannel <- dto.RestRequest{Type: common.TASK_SELECT, Data: &query, Response: responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

//This method handle histogram by value select
func (serv SelectService) HistogramByValueQuery(metrics, tags, times string, from, to float32, buckets int32) dto.RestResponse {
	log.Finest("Selecting Histogram")
	serv.setHeader()
	responseChan := make(chan *dto.RestResponse)
	query, err := prepareQueryWithoutAggregationType(metrics, tags, times)
	if err != nil {
		log.Error("Return HTTP 400")
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(
			[]byte("The request could not be understood by the server due to malformed syntax. You SHOULD NOT repeat the request without modifications."))
		return dto.RestResponse{}
	}
	query.AggregationType = common.AGGREGATION_HISTOGRAM_BY_VALUE
	query.AdditionalData = dto.HistogramValueData{Min: from, Max: to, BucketCount: buckets}
	log.Fine("Query: ", &query)
	restRequestChannel <- dto.RestRequest{Type: common.TASK_SELECT, Data: &query, Response: responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

//This method handle histogram by time select
func (serv SelectService) HistogramByTimeQuery(metrics, tags, times string, from, to int64, buckets int32) dto.RestResponse {
	log.Finest("Selecting Histogram")
	serv.setHeader()
	responseChan := make(chan *dto.RestResponse)
	query, err := prepareQueryWithoutAggregationType(metrics, tags, times)
	if err != nil {
		log.Error("Return HTTP 400")
		serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride(
			[]byte("The request could not be understood by the server due to malformed syntax. You SHOULD NOT repeat the request without modifications."))
		return dto.RestResponse{}
	}
	query.AggregationType = common.AGGREGATION_HISTOGRAM_BY_TIME
	query.AdditionalData = dto.HistogramTimeData{Min: from, Max: to, BucketCount: buckets}
	log.Fine("Query: ", &query)
	restRequestChannel <- dto.RestRequest{Type: common.TASK_SELECT, Data: &query, Response: responseChan}
	response := <-responseChan
	serv.setSelectHeaderErrors(response)
	return *response
}

func prepareQuery(metrics, tags, times, aggr string) (dto.Query, error) {

	query, err := prepareQueryWithoutAggregationType(metrics, tags, times)
	if err != nil {
		return dto.Query{}, err
	}

	query.AggregationType, err = prepareAggregationType(aggr)
	if err != nil {
		return dto.Query{}, err
	}

	ok, err := validateQuery(query)
	if !ok {
		return dto.Query{}, err
	}

	return query, nil
}

func validateQuery(query dto.Query) (bool, error) {

	if query.AggregationType == common.AGGREGATION_INTEGRAL && (query.MetricsCount != 1 || query.TagsCount != 1) {
		return false, fmt.Errorf("Integral can be done only for one series")
	}

	for i := int32(0); i < query.TimeSpansCount; i += 2 {
		if query.TimeSpans[i] > query.TimeSpans[i+1] {
			return false, fmt.Errorf("Beginning of time span must be less or equal then end")
		}
	}

	return true, nil
}

func prepareQueryWithoutAggregationType(metrics, tags, times string) (dto.Query, error) {

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

	return dto.Query{MetricsCount: int32(len(metricsArr)), Metrics: metricsArr, TagsCount: int32(len(tagsArr)), Tags: tagsArr, TimeSpansCount: int32(len(timesArr) / 2), TimeSpans: timesArr}, nil
}

func prepareAggregationType(aggregation string) (int32, error) {
	switch aggregation {
	case "none":
		return common.AGGREGATION_NONE, nil
	case "sum":
		return common.AGGREGATION_ADD, nil
	case "max":
		return common.AGGREGATION_MAX, nil
	case "min":
		return common.AGGREGATION_MIN, nil
	case "avg":
		return common.AGGREGATION_AVERAGE, nil
	case "std":
		return common.AGGREGATION_STDDEVIATION, nil
	case "var":
		return common.AGGREGATION_VARIANCE, nil
	case "int":
		return common.AGGREGATION_INTEGRAL, nil
	default:
		return common.TASK_ERROR, fmt.Errorf("Unknown aggregation type")
	}
}

const ALL string = "all"

func prepareTimeSpans(times string) ([]int64, error) {

	if times == ALL {
		return []int64{}, nil
	}

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

	if input == ALL {
		return make([]int32, 0), nil
	}

	set := make(map[int]bool)

	inputSplited := strings.Split(input, ",")
	for _, element := range inputSplited {
		value, err := strconv.Atoi(element)
		if err != nil {
			return nil, err
		}
		set[value] = true
	}

	unsorted := []int{}
	for element := range set {
		unsorted = append(unsorted, element)
	}

	sort.Ints(unsorted)
	inputArr := make([]int32, 0, len(set))
	for _, element := range unsorted {
		inputArr = append(inputArr, int32(element))
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
	} else if response.Error != "" {
		log.Error("Return HTTP 500")
		serv.ResponseBuilder().SetResponseCode(500).WriteAndOveride(
			[]byte(fmt.Sprintf("TaskId: %d, Error: %s", response.TaskId, response.Error)))
	}
}
