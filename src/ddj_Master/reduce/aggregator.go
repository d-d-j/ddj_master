package reduce

import (
	log "code.google.com/p/log4go"
	. "ddj_Master/common"
	"ddj_Master/dto"
	"math"
	"sort"
)

type Aggregator func([]Aggregates) dto.Dtos

type Aggregates dto.Dto

var aggregations map[int32]Aggregator

func Initialize() {
	aggregations = map[int32]Aggregator{
		AGGREGATION_NONE:         NonAggregation,
		AGGREGATION_MAX:          MaxAggregation,
		AGGREGATION_MIN:          MinAggregation,
		AGGREGATION_ADD:          AddAggregation,
		AGGREGATION_AVERAGE:      AverageAggregation,
		AGGREGATION_VARIANCE:     Variance,
		AGGREGATION_STDDEVIATION: StandartdDeviation,
		AGGREGATION_INTEGRAL:     Integral,
	}
}

func GetAggregator(aggregationType int32) Aggregator {
	aggregator := aggregations[aggregationType]
	if aggregator == nil {
		panic("Unknown aggregation.")
	}
	return aggregator
}

func NonAggregation(input []Aggregates) dto.Dtos {

	output := make([]dto.Dto, 0, len(input))
	for _, element := range input {
		e := element.(*dto.Element)
		output = append(output, e)
	}

	sort.Sort(dto.ByTime(output))

	return output
}

func MaxAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	max := input[0].(*dto.Element).Value

	for _, element := range input {
		e := element.(*dto.Element)
		if !e.Value.Less(max) {
			max = e.Value
		}
	}

	return dto.Dtos{&max}
}

func MinAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	min := input[0].(*dto.Element).Value

	for _, element := range input {
		e := element.(*dto.Element)
		if e.Value.Less(min) {
			min = e.Value
		}
	}

	return dto.Dtos{&min}
}

func AddAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}
	sum := dto.Value(0.0)

	for _, value := range input {
		v := *(value.(*dto.Value))
		sum += v
	}

	return dto.Dtos{&sum}
}

func AverageAggregation(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}

	var (
		average dto.Value
		count   int32
	)

	for _, variance := range input {
		v := variance.(*dto.AverageElement)
		average += v.Sum
		count += v.Count
	}

	average /= dto.Value(count)

	return dto.Dtos{&average}
}

func StandartdDeviation(input []Aggregates) dto.Dtos {

	length := len(input)
	if length < 1 {
		return dto.Dtos{}
	}

	data := make([]*dto.VarianceElement, length)

	for index, variance := range input {
		data[index] = variance.(*dto.VarianceElement)
	}

	σ := dto.Value(math.Sqrt(variance(data)))

	return dto.Dtos{&σ}
}

func Variance(input []Aggregates) dto.Dtos {

	length := len(input)
	if length < 1 {
		return dto.Dtos{}
	}

	data := make([]*dto.VarianceElement, length)

	for index, variance := range input {
		data[index] = variance.(*dto.VarianceElement)
	}

	σ := dto.Value(variance(data))

	return dto.Dtos{&σ}
}

func variance(data []*dto.VarianceElement) float64 {

	length := len(data)
	if length == 1 && data[0].Count <= 1 {
		return 0.0
	}
	var x dto.VarianceElement
	x = *data[0]

	for i := 1; i < length; i++ {
		y := data[i]

		n := dto.Value(x.Count + y.Count)
		delta := y.Mean - x.Mean
		delta2 := delta * delta
		mean := x.Mean + delta*dto.Value(y.Count)/n
		M2 := x.M2 + y.M2
		M2 += delta2 * dto.Value(x.Count*y.Count) / n

		x.Count = x.Count + y.Count
		x.Mean = mean
		x.M2 = M2
	}

	return float64(x.M2 / dto.Value(x.Count-1))
}

func Integral(input []Aggregates) dto.Dtos {

	var integral dto.Value
	length := len(input)
	if length < 1 {
		return dto.Dtos{&integral}
	}

	data := make([]*dto.IntegralElement, length)

	for index, variance := range input {
		data[index] = variance.(*dto.IntegralElement)
		log.Warn(data[index])
	}

	sort.Sort(dto.ByLeftTime(data))

	integral = data[0].Integral
	for i := 1; i < length; i++ {
		x := data[i-1]
		y := data[i]
		integral += y.Integral
		integral += (x.RightValue + y.LeftValue) * dto.Value(y.LeftTime-x.RightTime) / 2
	}

	return dto.Dtos{&integral}
}
