package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"testing"
)

func Test_Initialize(t *testing.T) {
	expected := 4
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
	input := make([]*dto.Element, 5)
	defer ExpectedPanic(t)
	NonAggregation(input)
}

func Test_NonAggregation_Should_Return_Same_Slice_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {

	expected := []*dto.Element{dto.NewElement(1, 2, 0, 0.33)}
	input := expected

	actual := NonAggregation(input)
	AssertElementsEqual(expected, actual, t)
}

func Test_NonAggregation_Should_Return_Same_Slice_But_Sorted_As_Input_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	expected := []*dto.Element{dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 4, 2, 0.33)}
	input := []*dto.Element{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := NonAggregation(input)

	AssertElementsEqual(expected, actual, t)
}

func Test_NonAggregation_Should_Return_Sorted_Elements_From_All_Input_Slices(t *testing.T) {

	input := []*dto.Element{
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33),
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33),
	}

	actual := NonAggregation(input)
	expected := []*dto.Element{
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
	input := make([]*dto.Element, 5)
	defer ExpectedPanic(t)
	MaxAggregation(input)
}

func Test_MaxAggregation_Should_Return_Same_Slice_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {
	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []*dto.Element{dto.NewElement(1, 2, 0, 0.33)}

	actual := MaxAggregation(input)
	AssertValuesEqual(expected, actual, t)
}

func Test_MaxAggregation_Should_Return_Maximum_From_Input_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []*dto.Element{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := MaxAggregation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_MaxAggregation_Should_Return_Max_Element_From_All_Input_Slices(t *testing.T) {

	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []*dto.Element{
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.23),
		dto.NewElement(1, 4, 2, 0.13), dto.NewElement(1, 2, 1, 0.033), dto.NewElement(1, 0, 0, 0.3),
	}

	actual := MaxAggregation(input)

	AssertValuesEqual(expected, actual, t)

}

func Test_MinAggregation_Should_Panic_When_Input_Containse_nil(t *testing.T) {
	input := make([]*dto.Element, 5)
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
	input := []*dto.Element{dto.NewElement(1, 2, 0, 0.33)}

	actual := MinAggregation(input)
	AssertValuesEqual(expected, actual, t)
}

func Test_MinAggregation_Should_Return_Minimum_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	var value dto.Value
	value = 0.33
	expected := []*dto.Value{&value}
	input := []*dto.Element{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := MinAggregation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_MinAggregation_Should_Return_Max_Element_From_All_Input_Slices(t *testing.T) {

	var value dto.Value
	value = 0.0
	expected := []*dto.Value{&value}
	input := []*dto.Element{
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.23),
		dto.NewElement(1, 4, 2, 0.13), dto.NewElement(1, 2, 1, 0.0), dto.NewElement(1, 0, 0, 0.3),
	}

	actual := MinAggregation(input)

	AssertValuesEqual(expected, actual, t)

}

func Test_AddAggregation_Should_Panic_When_Input_Containse_nil(t *testing.T) {
	input := make([]*dto.Element, 5)
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
	input := []*dto.Element{dto.NewElement(1, 2, 0, 0.33)}

	actual := AddAggregation(input)
	AssertValuesEqual(expected, actual, t)
}

func Test_AddAggregation_Should_Return_Sum_of_Slice_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	var value dto.Value
	value = 0.99
	expected := []*dto.Value{&value}
	input := []*dto.Element{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := AddAggregation(input)

	AssertValuesEqual(expected, actual, t)
}

func Test_AddAggregation_Should_Return_Sum_Of_Elements_From_All_Input_Slices(t *testing.T) {

	var value dto.Value
	value = 0.6
	expected := []*dto.Value{&value}
	input := []*dto.Element{
		dto.NewElement(1, 4, 2, 0.01), dto.NewElement(1, 2, 1, 0.03), dto.NewElement(1, 0, 0, 0.04),
		dto.NewElement(1, 4, 2, 0.02), dto.NewElement(1, 2, 1, 0.0), dto.NewElement(1, 0, 0, 0.5),
	}

	actual := AddAggregation(input)

	AssertValuesEqual(expected, actual, t)

}

func AssertElementsEqual(expected []*dto.Element, actual dto.Dtos, t *testing.T) {
	if len(expected) != len(actual) {
		t.Error("Wrong dimension. Expected ", len(expected), " but got ", len(actual))
	}

	for index, elem := range actual {
		if expected[index].String() != elem.String() {
			t.Error("Expected: ", expected, " but got: ", actual)
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
