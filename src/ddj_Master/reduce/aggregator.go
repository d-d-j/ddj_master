package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"sort"
)

type Aggregator interface {
	Aggregate([]*dto.RestResponse) dto.Dtos
}

type NonAggregation struct{}

func GetAggregator(aggregationType int32) Aggregator {
	switch aggregationType {
	case common.AGGREGATION_NONE:
		return NonAggregation{}
	}
	panic("Unknown aggregation")
}

func (this NonAggregation) Aggregate(input []*dto.RestResponse) dto.Dtos {

	if input == nil {
		return nil
	}

	totalSize := 0
	for _, element := range input {
		if element != nil {
			totalSize += len(element.Data)
		}
	}

	output := make([]dto.Dto, 0, totalSize)
	for _, element := range input {
		if element != nil {
			output = append(output, element.Data...)
		}
	}

	sort.Sort(dto.ByTime(output))

	return output
}
