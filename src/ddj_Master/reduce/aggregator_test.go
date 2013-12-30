package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"testing"
)

func Test_NonAggregation_Should_Return_Nil_When_Input_Is_Nil(t *testing.T) {
	actual := NonAggregation(nil)
	if actual != nil {
		t.Error("Expected nil but got ", actual)
	}
}

func Test_NonAggregation_Should_Return_Empty_Slice_If_Input_Containse_Only_nil(t *testing.T) {
	input := make([]*dto.Element, 5)

	actual := NonAggregation(input)
	if len(actual) != 0 && actual != nil {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_NonAggregation_Should_Return_Same_Slice_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {

	expected := []*dto.Element{dto.NewElement(1, 2, 0, 0.33)}
	input := expected

	actual := NonAggregation(input)
	for index, elem := range actual {
		AssertEqual(expected[index], elem, t)
	}
}

func Test_NonAggregation_Should_Return_Same_Slice_But_Sorted_As_Input_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {

	expected := []*dto.Element{dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 4, 2, 0.33)}
	input := []*dto.Element{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)}

	actual := NonAggregation(input)

	for index, elem := range actual {
		AssertEqual(expected[index], elem, t)
	}
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
	for index, elem := range actual {
		AssertEqual(expected[index], elem, t)
	}
}

func AssertEqual(expected, actual dto.Dto, t *testing.T) {
	if expected.String() != actual.String() {
		t.Error("Expected: ", expected, " but got: ", actual)
	}
}

func Test_GetAggregator(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered ", r)
		}
	}()
	GetAggregator(common.CONST_UNINITIALIZED)
	t.Error("Expected panic")
}

func Test_GetAggregator_Should_Pannic_When_Pass_Invalid_Aggregation_Type(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()
	GetAggregator(common.CONST_UNINITIALIZED) //This should call panic
}
