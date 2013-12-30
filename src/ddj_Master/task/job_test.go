package task

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"testing"
)

func Test_ParseResultsToElements_Should_Return_Empty_Slice_For_Nil_Imput(t *testing.T) {
	actual := parseResultsToElements(nil)
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_ParseResultsToElements_Should_Return_Empty_Slice_For_Empty_Imput(t *testing.T) {
	actual := parseResultsToElements([]*dto.Result{})
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_ParseResultsToElements_Should_Return_Empty_Slice_For_Empty_Data_Imput(t *testing.T) {
	data := []byte{}
	result := dto.NewResult(0, common.TASK_SELECT, 0, data)
	actual := parseResultsToElements([]*dto.Result{result, result, result})
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_ParseResultsToElements_Should_Return_One_Elements_When_Called_With_One_In_Slice(t *testing.T) {
	// PREPARE DATA FOR TEST
	expected := dto.NewElement(1, 2, 3, 0.33)
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := dto.NewResult(0, common.TASK_SELECT, int32(expected.Size()), data)
	actual := parseResultsToElements([]*dto.Result{result})
	if actual[0].String() != expected.String() {
		t.Error("Expected empty slice but got ", actual)
	}
}
