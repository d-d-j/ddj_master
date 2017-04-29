//This package ocntains function used to reduce aggregation results from nodes
package reduce

import (
	. "github.com/d-d-j/ddj_master/common"
	"github.com/d-d-j/ddj_master/dto"
	"math"
	"sort"
)

//This type is function that will handle aggregation of given input
type Aggregator func([]Aggregates) dto.Dtos

//This is alias for argument types of aggregation method
type Aggregates dto.Dto

var aggregations map[int32]Aggregator

//Initialize aggregations. Should be called before GetAggregator
func Initialize() {
	aggregations = map[int32]Aggregator{
		AGGREGATION_NONE:               nonAggregation,
		AGGREGATION_MAX:                maxAggregation,
		AGGREGATION_MIN:                minAggregation,
		AGGREGATION_ADD:                addAggregation,
		AGGREGATION_AVERAGE:            averageAggregation,
		AGGREGATION_VARIANCE:           variance,
		AGGREGATION_STDDEVIATION:       standartdDeviation,
		AGGREGATION_INTEGRAL:           integral,
		AGGREGATION_HISTOGRAM_BY_TIME:  histogram,
		AGGREGATION_HISTOGRAM_BY_VALUE: histogram,
		AGGREGATION_SERIES_SUM:         seriesSum,
	}
}

//Return Aggregator proper for given aggregation type
func GetAggregator(aggregationType int32) Aggregator {
	aggregator := aggregations[aggregationType]
	if aggregator == nil {
		panic("Unknown aggregation.")
	}
	return aggregator
}

func nonAggregation(input []Aggregates) dto.Dtos {

	output := make([]dto.Dto, 0, len(input))
	for _, element := range input {
		e := element.(*dto.Element)
		output = append(output, e)
	}

	sort.Sort(dto.ByTime(output))

	return output
}

func maxAggregation(input []Aggregates) dto.Dtos {

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

func minAggregation(input []Aggregates) dto.Dtos {

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

func addAggregation(input []Aggregates) dto.Dtos {

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

func averageAggregation(input []Aggregates) dto.Dtos {

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

func standartdDeviation(input []Aggregates) dto.Dtos {

	length := len(input)
	if length < 1 {
		return dto.Dtos{}
	}

	data := make([]*dto.VarianceElement, length)

	for index, variance := range input {
		data[index] = variance.(*dto.VarianceElement)
	}

	σ := dto.Value(math.Sqrt(varianceHelper(data)))

	return dto.Dtos{&σ}
}

func variance(input []Aggregates) dto.Dtos {

	length := len(input)
	if length < 1 {
		return dto.Dtos{}
	}

	data := make([]*dto.VarianceElement, length)

	for index, variance := range input {
		data[index] = variance.(*dto.VarianceElement)
	}

	σ := dto.Value(varianceHelper(data))

	return dto.Dtos{&σ}
}

func varianceHelper(data []*dto.VarianceElement) float64 {

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

func integral(input []Aggregates) dto.Dtos {

	var integral dto.Value
	length := len(input)
	if length < 1 {
		return dto.Dtos{&integral}
	}

	data := make([]*dto.IntegralElement, length)

	for index, variance := range input {
		data[index] = variance.(*dto.IntegralElement)
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

func histogram(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}

	length := len(input[0].(*dto.Histogram).Data)
	histogram := dto.Histogram{make([]int32, length)}

	for _, h := range input {
		for index, value := range h.(*dto.Histogram).Data {
			histogram.Data[index] += value
		}
	}

	return dto.Dtos{&histogram}
}

func seriesSum(input []Aggregates) dto.Dtos {

	if len(input) < 1 {
		return dto.Dtos{}
	}

	length := len(input[0].(*dto.InterpolateElement).Data)
	seriesSum := dto.InterpolateElement{make([]dto.Value, length)}

	for _, h := range input {
		for index, value := range h.(*dto.InterpolateElement).Data {
			seriesSum.Data[index] += value
		}
	}

	return dto.Dtos{&seriesSum}
}
