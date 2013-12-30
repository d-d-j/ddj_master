package task

import (
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
