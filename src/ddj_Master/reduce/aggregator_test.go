package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"math"
	"testing"
)

func Test_Initialize(t *testing.T) {
	expected := 11
	Initialize()
	actual := len(aggregations)
	if actual != expected {
		t.Error("Expected: ", expected, " but got: ", actual)
	}
}

func Test_NonAggregation_Should_Return_Empty_Slice_When_Input_Is_Nil_Or_Empty(t *testing.T) {
	actual := NonAggregation(nil)
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_NonAggregation_Should_Panic_When_Input_Containse_nil(t *testing.T) {
	input := make([]Aggregates, 5)
	defer ExpectedPanic(t)
	NonAggregation(input)
}

func Test_NonAggregation_Should_Return_Same_Slice_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {

	expected := []Aggregates{dto.NewElement(1, 2, 0, 0.33)}
	input := expected

	actual := NonAggregation(input)
	AssertElementsEqual(expected, actual, t)
}

func Test_NonAggregation_Should_Return_Same_Slice_But_Sorted_As_Input_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	expected := []Aggregates{dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 4, 2, 0.33)}
	input := []Aggregates{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := NonAggregation(input)

	AssertElementsEqual(expected, actual, t)
}

func Test_NonAggregation_Should_Return_Sorted_Elements_From_All_Input_Slices(t *testing.T) {

	input := []Aggregates{
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33),
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33),
	}

	actual := NonAggregation(input)
	expected := []Aggregates{
		dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 2, 1, 0.33),
		dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 4, 2, 0.33),
	}
	AssertElementsEqual(expected, actual, t)
}

func Test_MaxAggregation_Should_Return_Empty_Slice_When_Input_Is_Nil_Or_Empty(t *testing.T) {
	actual := MaxAggregation(nil)
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_MaxAggregation_Should_Panic_When_Input_Containse_nil(t *testing.T) {
	input := make([]Aggregates, 5)
	defer ExpectedPanic(t)
	MaxAggregation(input)
}

func Test_MaxAggregation_Should_Return_Same_Slice_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {
	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{dto.NewElement(1, 2, 0, 0.33)}

	actual := MaxAggregation(input)
	AssertValuesEqual(expected, actual, t)
}

func Test_MaxAggregation_Should_Return_Maximum_From_Input_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := MaxAggregation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_MaxAggregation_Should_Return_Max_Element_From_All_Input_Slices(t *testing.T) {

	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.23),
		dto.NewElement(1, 4, 2, 0.13), dto.NewElement(1, 2, 1, 0.033), dto.NewElement(1, 0, 0, 0.3),
	}

	actual := MaxAggregation(input)

	AssertValuesEqual(expected, actual, t)

}

func Test_MinAggregation_Should_Panic_When_Input_Containse_nil(t *testing.T) {
	input := make([]Aggregates, 5)
	defer ExpectedPanic(t)
	MinAggregation(input)
}

func Test_MinAggregation_Should_Return_Empty_Slice_When_Input_Is_Nil_Or_Empty(t *testing.T) {
	actual := MinAggregation(nil)
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_MinAggregation_Should_Return_Same_Slice_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {
	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{dto.NewElement(1, 2, 0, 0.33)}

	actual := MinAggregation(input)
	AssertValuesEqual(expected, actual, t)
}

func Test_MinAggregation_Should_Return_Minimum_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := MinAggregation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_MinAggregation_Should_Return_Max_Element_From_All_Input_Slices(t *testing.T) {

	var value dto.Value
	value = 0.0
	expected := []*dto.Value{&value}
	input := []Aggregates{
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.23),
		dto.NewElement(1, 4, 2, 0.13), dto.NewElement(1, 2, 1, 0.0), dto.NewElement(1, 0, 0, 0.3),
	}

	actual := MinAggregation(input)

	AssertValuesEqual(expected, actual, t)

}

func Test_AddAggregation_Should_Panic_When_Input_Containse_nil(t *testing.T) {
	input := make([]Aggregates, 5)
	defer ExpectedPanic(t)
	AddAggregation(input)
}

func Test_AddAggregation_Should_Return_Empty_Slice_When_Input_Is_Nil_Or_Empty(t *testing.T) {
	actual := AddAggregation(nil)
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_AddAggregation_Should_Return_Same_Value_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {
	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{&value}

	actual := AddAggregation(input)
	AssertValuesEqual(expected, actual, t)
}

func Test_AddAggregation_Should_Return_Sum_of_Slice_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	var value dto.Value
	value = 0.99
	exp := 3 * value
	expected := []*dto.Value{&exp}
	input := []Aggregates{&value, &value, &value}

	actual := AddAggregation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_AverageAggregation_Should_Panic_When_Input_Containse_nil(t *testing.T) {
	input := make([]Aggregates, 5)
	defer ExpectedPanic(t)
	AverageAggregation(input)
}

func Test_AverageAggregation_Should_Return_Empty_Slice_When_Input_Is_Nil_Or_Empty(t *testing.T) {
	actual := AverageAggregation(nil)
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_AverageAggregation_Should_Return_Same_Value_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {
	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{&dto.AverageElement{1, 1 * value}}

	actual := AverageAggregation(input)
	AssertValuesEqual(expected, actual, t)
}

func Test_AverageAggregation_Should_Return_Avg_of_Slice(t *testing.T) {

	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []Aggregates{&dto.AverageElement{6, 12 * value}, &dto.AverageElement{6, 0}}

	actual := AverageAggregation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_StandardDeviation_Should_Return_StdDev_of_Slice(t *testing.T) {

	var value dto.Value
	value = dto.Value(math.Sqrt(10.0 / 3.0))
	expected := []*dto.Value{&value}
	input := []Aggregates{&dto.VarianceElement{1, 5, 0}, &dto.VarianceElement{1, 6, 0}, &dto.VarianceElement{1, 8, 0}, &dto.VarianceElement{1, 9, 0}}

	actual := StandartdDeviation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_StandardDeviation_Should_Return_StdDev_of_Slice_2(t *testing.T) {
	//http://en.wikipedia.org/wiki/Standard_deviation#Basic_examples
	var value dto.Value
	value = dto.Value(2.138089935)
	expected := []*dto.Value{&value}
	input := []Aggregates{&dto.VarianceElement{1, 2, 0}, &dto.VarianceElement{3, 4, 0}, &dto.VarianceElement{2, 5, 0}, &dto.VarianceElement{1, 7, 0}, &dto.VarianceElement{1, 9, 0}}

	actual := StandartdDeviation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_StandardDeviation_Should_Return_0_For_One_Element_And_Count_1(t *testing.T) {
	var value dto.Value
	value = dto.Value(0.0)
	expected := []*dto.Value{&value}
	input := []Aggregates{&dto.VarianceElement{1, 5, 0}}

	actual := StandartdDeviation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_Integral_Should_Return_Integral_of_Slice(t *testing.T) {

	var value dto.Value
	value = dto.Value(2.2)
	expected := []*dto.Value{&value}
	input := []Aggregates{&dto.IntegralElement{0.5, 0.2, 1, 0.1, 2}, &dto.IntegralElement{0.5, 0.0, 10, 0.2, 20}, &dto.IntegralElement{0.5, 0.1, 3, 0.2, 4}}

	actual := Integral(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_Histogram_Should_Return_Sum_Of_All_Histograms(t *testing.T) {

	histogram := dto.Histogram{[]int32{0, 1, 1, 5, 7, 2}}
	expected := dto.Histogram{[]int32{0, 3, 3, 15, 21, 6}}
	input := []Aggregates{&histogram, &histogram, &histogram}

	actual := Histogram(input)

	if expected.String() != actual[0].String() {
		t.Error("Expected: ", expected, " but got: ", actual)
	}
}

func AssertElementsEqual(expected []Aggregates, actual dto.Dtos, t *testing.T) {
	if len(expected) != len(actual) {
		t.Error("Wrong dimension. Expected ", len(expected), " but got ", len(actual))
	}

	for index, elem := range actual {
		e := expected[index].(*dto.Element)
		if e.String() != elem.String() {
			t.Error("Expected: ", e, " but got: ", actual)
		}
	}
}

func AssertValuesEqual(expected []*dto.Value, actual dto.Dtos, t *testing.T) {
	if len(expected) != len(actual) {
		t.Error("Wrong dimension. Expected ", len(expected), " but got ", len(actual))
	}

	for index, elem := range actual {
		if expected[index].String() != elem.String() {
			t.Error("Expected: ", expected, " but got: ", actual)
		}
	}
}

func Test_GetAggregator_Should_Panic_When_Unknow_Type_Pass(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered ", r)
		}
	}()
	GetAggregator(common.CONST_UNINITIALIZED)
	t.Error("Expected panic")
}

func Test_GetAggregator_Should_Pannic_When_Pass_Invalid_Aggregation_Type(t *testing.T) {
	defer ExpectedPanic(t)
	GetAggregator(common.CONST_UNINITIALIZED) //This should call panic
}

func ExpectedPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Error("Expected panic")
	}
}
