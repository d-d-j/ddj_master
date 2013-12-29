package reduce

import (
	"ddj_Master/dto"
	"sort"
)

type Aggregator interface {
	Aggregate([]*dto.RestResponse) dto.Dtos
}

type NonAggregation struct{}

func (this NonAggregation) Aggregate(input []*dto.RestResponse) dto.Dtos {

	if input == nil {
		return nil
	}

	output := make([]dto.Dto, 0, len(input))
	for _, element := range input {
		if element != nil {
			output = append(output, element.Data...)
		}
	}

	sort.Sort(dto.ByTime(output))

	return output
}
