package reduce

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"testing"
)

var nonAggr NonAggregation

func Test_NonAggregation_Should_Return_Nil_When_Input_Is_Nil(t *testing.T) {
	actual := nonAggr.Aggregate(nil)
	if actual != nil {
		t.Error("Expected nil but got ", actual)
	}
}

func Test_NonAggregation_Should_Return_Empty_Slice_If_Input_Containse_Only_nil(t *testing.T) {
	input := make([]*dto.RestResponse, 5)

	actual := nonAggr.Aggregate(input)
	if len(actual) != 0 && actual != nil {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_NonAggregation_Should_Return_Same_Slice_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {
	input := make([]*dto.RestResponse, 1)
	expected := dto.Dtos{dto.NewElement(1, 2, 0, 0.33)}
	input[0] = dto.NewRestResponse("", 0, expected)

	actual := nonAggr.Aggregate(input)
	if actual.String() != expected.String() {
		t.Error("Expected ", expected, " but got ", actual)
	}
}

func Test_NonAggregation_Should_Return_Same_Slice_But_Sorted_As_Input_If_There_Was_Only_One_Slice_In_Input(t *testing.T) {
	input := make([]*dto.RestResponse, 1)
	input[0] = dto.NewRestResponse("", 0, dto.Dtos{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)})

	actual := nonAggr.Aggregate(input)
	expected := dto.Dtos{dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 4, 2, 0.33)}
	if actual.String() != expected.String() {
		t.Error("Expected ", expected, " but got ", actual)
	}
}

func Test_NonAggregation_Should_Return_Sorted_Elements_From_All_Input_Slices(t *testing.T) {
	input := make([]*dto.RestResponse, 2)
	input[0] = dto.NewRestResponse("", 0, dto.Dtos{
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33),
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33),
	})
	input[1] = dto.NewRestResponse("", 0, dto.Dtos{dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 0, 0, 0.33)})

	actual := nonAggr.Aggregate(input)
	expected := dto.Dtos{
		dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 0, 0, 0.33), dto.NewElement(1, 0, 0, 0.33),
		dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 2, 1, 0.33), dto.NewElement(1, 2, 1, 0.33),
		dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 4, 2, 0.33), dto.NewElement(1, 4, 2, 0.33),
	}
	if actual.String() != expected.String() {
		t.Error("Expected ", expected, " but got ", actual)
	}
}

func Test_GetAggregator(t *testing.T) {
	actual := GetAggregator(common.AGGREGATION_NONE)
	_, ok := actual.(NonAggregation)
	if !ok {
		t.Error("Expected NonAggregation")
	}
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered ", r)
		}
	}()
	actual = GetAggregator(common.CONST_UNINITIALIZED)
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
