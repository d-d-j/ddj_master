package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"sort"
)

type Aggregator func([]*dto.Element) dto.Dtos

var aggregations map[int32]Aggregator

func Initialize() {
	aggregations = map[int32]Aggregator{
		common.AGGREGATION_NONE: NonAggregation,
		common.AGGREGATION_MAX:  MaxAggregation,
		common.AGGREGATION_MIN:  MinAggregation,
		common.AGGREGATION_ADD:  AddAggregation,
	}
}

func GetAggregator(aggregationType int32) Aggregator {
	aggregator := aggregations[aggregationType]
	if aggregator == nil {
		panic("Unknown aggregation.")
	}
	return aggregator
}

func NonAggregation(input []*dto.Element) dto.Dtos {

	output := make([]dto.Dto, 0, len(input))
	for _, element := range input {
		output = append(output, element)
	}

	sort.Sort(dto.ByTime(output))

	return output
}

func MaxAggregation(input []*dto.Element) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	max := input[0].Value

	for _, element := range input {
		if !element.Value.Less(max) {
			max = element.Value
		}
	}

	return dto.Dtos{&max}
}

func MinAggregation(input []*dto.Element) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	min := input[0].Value

	for _, element := range input {
		if element.Value.Less(min) {
			min = element.Value
		}
	}

	return dto.Dtos{&min}
}

func AddAggregation(input []*dto.Element) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	sum := dto.Value(0.0)

	for _, element := range input {
		sum += element.Value
	}

	return dto.Dtos{&sum}
}
