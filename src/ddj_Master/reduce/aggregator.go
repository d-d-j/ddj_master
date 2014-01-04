package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"sort"
)

type Aggregator func([]Aggregates) dto.Dtos

type Aggregates interface{}

var aggregations map[int32]Aggregator

func Initialize() {
	aggregations = map[int32]Aggregator{
		common.AGGREGATION_NONE:    NonAggregation,
		common.AGGREGATION_MAX:     MaxAggregation,
		common.AGGREGATION_MIN:     MinAggregation,
		common.AGGREGATION_ADD:     AddAggregation,
		common.AGGREGATION_AVERAGE: AverageAggregation,
	}
}

func GetAggregator(aggregationType int32) Aggregator {
	aggregator := aggregations[aggregationType]
	if aggregator == nil {
		panic("Unknown aggregation.")
	}
	return aggregator
}

func NonAggregation(input []Aggregates) dto.Dtos {

	output := make([]dto.Dto, 0, len(input))
	for _, element := range input {
		e := element.(*dto.Element)
		output = append(output, e)
	}

	sort.Sort(dto.ByTime(output))

	return output
}

func MaxAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	max := input[0].(*dto.Element).Value

	for _, element := range input {
		e := element.(*dto.Element)
		if !e.Value.Less(max) {
			max = e.Value
		}
	}

	return dto.Dtos{&max}
}

func MinAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	min := input[0].(*dto.Element).Value

	for _, element := range input {
		e := element.(*dto.Element)
		if e.Value.Less(min) {
			min = e.Value
		}
	}

	return dto.Dtos{&min}
}

func AddAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	sum := dto.Value(0.0)

	for _, value := range input {
		v := *(value.(*dto.Value))
		sum += v
	}

	return dto.Dtos{&sum}
}

func AverageAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}

	var (
		average dto.Value
		count   int32
	)

	for _, variance := range input {
		v := variance.(*dto.AverageElement)
		average += v.Sum
		count += v.Count
	}

	average /= dto.Value(count)

	return dto.Dtos{&average}
}
