package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"sort"
)

type Aggregator func([]*dto.Element) dto.Dtos

func GetAggregator(aggregationType int32) Aggregator {
	switch aggregationType {
	case common.AGGREGATION_NONE:
		return NonAggregation
	case common.AGGREGATION_MAX:
		return MaxAggregation
	case common.AGGREGATION_MIN:
		return MinAggregation
	case common.AGGREGATION_ADD:
		return AddAggregation
	}
	panic("Unknown aggregation")
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
