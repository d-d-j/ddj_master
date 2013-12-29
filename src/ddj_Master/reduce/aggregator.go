package reduce

import (
	"ddj_Master/dto"
)

type Aggregator interface {
	Aggregate([]dto.Dtos) dto.Dtos
}

type NonAggregation struct{}

func (this NonAggregation) Aggregate(input []dto.Dtos) dto.Dtos {
	if len(input) == 1 {
		return input[0]
	}
	return nil
}
