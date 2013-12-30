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
	x := input[0].Value

	for _, y := range input {
		if !y.Value.Less(x) {
			x = y.Value
		}
	}

	return dto.Dtos{&x}
}

func MinAggregation(input []*dto.Element) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	x := input[0].Value

	for _, y := range input {
		if y.Value.Less(x) {
			x = y.Value
		}
	}

	return dto.Dtos{&x}
}
