package reduce

import (
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
	input := make([]dto.Dtos, 5)

	actual := nonAggr.Aggregate(input)
	if len(actual) != 0 && actual != nil {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_NonAggregation_Should_Return_Same_Slice_But_As_Input_If_There_Was_Only_One_Element_In_Input_Slice(t *testing.T) {
	input := make([]dto.Dtos, 1)
	expected := dto.Dtos{dto.NewElement(1, 2, 0, 0.33)}
	input[0] = expected

	actual := nonAggr.Aggregate(input)
	if actual.String() != expected.String() {
		t.Error("Expected ", expected, " but got ", actual)
	}
}
