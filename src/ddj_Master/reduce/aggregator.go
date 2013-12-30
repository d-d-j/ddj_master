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
	return nil
}
