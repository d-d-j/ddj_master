package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"sort"
)

type Aggregator interface {
	Aggregate([]*dto.Element) dto.Dtos
}

type NonAggregation struct{}

func GetAggregator(aggregationType int32) Aggregator {
	switch aggregationType {
	case common.AGGREGATION_NONE:
		return NonAggregation{}
	}
	panic("Unknown aggregation")
}

func (this NonAggregation) Aggregate(input []*dto.Element) dto.Dtos {

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
