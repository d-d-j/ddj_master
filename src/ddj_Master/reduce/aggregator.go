package reduce

import (
	"ddj_Master/dto"
)

type Aggregator interface {
	Aggregate([][]dto.Dto) []dto.Dto
}

type NonAggregation struct{}

func (this NonAggregation) Aggregate([][]dto.Dto) []dto.Dto {
	return nil
}
