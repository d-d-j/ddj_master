package reduce

import (
	"ddj_Master/dto"
	"sort"
)

type Aggregator interface {
	Aggregate([][]*dto.RestResponse) dto.Dtos
}

type NonAggregation struct{}

func (this NonAggregation) Aggregate(input [][]*dto.RestResponse) dto.Dtos {

	if input == nil {
		return nil
	}

	output := make([]dto.Dto, 0, len(input))
	for j := 0; j < len(input); j++ {
		for i := 0; i < len(input[j]); i++ {
			output = append(output, input[j][i].Data...)
		}
	}

	sort.Sort(dto.ByTime(output))

	return output
}
