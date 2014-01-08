package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/reduce"
)

func parseResults(results []*dto.Result, aggregationType int32) []reduce.Aggregates {
	switch aggregationType {
	case common.AGGREGATION_ADD:
		return parseResultsToValues(results)
	case common.AGGREGATION_AVERAGE:
		return parseResultsToAverage(results)
	case common.AGGREGATION_STDDEVIATION:
		return parseResultsToVariance(results)
	case common.AGGREGATION_VARIANCE:
		return parseResultsToVariance(results)
	case common.AGGREGATION_INTEGRAL:
		return parseResultsToIntegralElements(results)
	case common.AGGREGATION_HISTOGRAM_BY_TIME:
		return parseResultsToHistograms(results)
	case common.AGGREGATION_HISTOGRAM_BY_VALUE:
		return parseResultsToHistograms(results)
	default:
		return parseResultsToElements(results)
	}
}

func parseResultsToElements(results []*dto.Result) []reduce.Aggregates {
	elementSize := (&dto.Element{}).Size()
	resultsCount := len(results)
	elements := make([]reduce.Aggregates, 0, resultsCount)

	for _, result := range results {
		length := len(result.Data) / elementSize
		for j := 0; j < length; j++ {
			var e dto.Element
			err := e.Decode(result.Data[j*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			elements = append(elements, &e)
		}
	}
	return elements
}

func parseResultsToValues(results []*dto.Result) []reduce.Aggregates {
	elementSize := (dto.Value(0)).Size()
	resultsCount := len(results)
	values := make([]reduce.Aggregates, 0, resultsCount)
	log.Fine("Parsing %d results", resultsCount)
	for _, result := range results {
		length := len(result.Data) / elementSize
		for j := 0; j < length; j++ {
			var e dto.Value
			err := e.Decode(result.Data[j*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			values = append(values, &e)
		}
	}
	return values
}

func parseResultsToAverage(results []*dto.Result) []reduce.Aggregates {
	elementSize := (&dto.AverageElement{}).Size()
	resultsCount := len(results)
	values := make([]reduce.Aggregates, 0, resultsCount)
	for _, result := range results {
		length := len(result.Data) / elementSize
		for j := 0; j < length; j++ {
			var e dto.AverageElement
			err := e.Decode(result.Data[j*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			values = append(values, &e)
		}
	}
	return values
}

func parseResultsToVariance(results []*dto.Result) []reduce.Aggregates {
	elementSize := (&dto.VarianceElement{}).Size()
	resultsCount := len(results)
	values := make([]reduce.Aggregates, 0, resultsCount)
	for _, result := range results {
		length := len(result.Data) / elementSize
		for j := 0; j < length; j++ {
			var e dto.VarianceElement
			err := e.Decode(result.Data[j*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			values = append(values, &e)
		}
	}
	return values
}

func parseResultsToIntegralElements(results []*dto.Result) []reduce.Aggregates {
	elementSize := dto.IntegralElement{}.Size()
	resultsCount := len(results)
	values := make([]reduce.Aggregates, 0, resultsCount)
	for _, result := range results {
		length := len(result.Data) / elementSize
		for j := 0; j < length; j++ {
			var e dto.IntegralElement
			err := e.Decode(result.Data[j*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			values = append(values, &e)
		}
	}
	return values
}

func parseResultsToHistograms(results []*dto.Result) []reduce.Aggregates {
	panic("Not implemented")
}
