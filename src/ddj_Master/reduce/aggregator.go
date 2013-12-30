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
		return NonAggregation
	}
	panic("Unknown aggregation")
}

func NonAggregation(input []*dto.Element) dto.Dtos {

	if input == nil {
		return nil
	}

	output := make([]dto.Dto, 0, len(input))
	for _, element := range input {
		if element != nil {
			output = append(output, element)
		}
	}

	sort.Sort(dto.ByTime(output))

	return output
}

func MaxAggregation(input []*dto.Element) dto.Dtos {
	if input == nil {
		return nil
	}
	var ok bool
	x := dto.Value(common.CONST_INT_MIN_VALUE)
	for _, y := range input {
		if y != nil && y.Value > x {
			x = y.Value
			ok = true
		}
	}

	if ok {
		return dto.Dtos{&x}
	}
	return dto.Dtos{}
}

func MinAggregation(input []*dto.Element) dto.Dtos {
	if input == nil {
		return nil
	}
	var ok bool
	x := dto.Value(common.CONST_INT_MAX_VALUE)
	for _, y := range input {
		if y != nil && y.Value < x {
			x = y.Value
			ok = true
		}
	}

	if ok {
		return dto.Dtos{&x}
	}
	return dto.Dtos{}
}
