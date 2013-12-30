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

func Test_ParseResultsToElements_Should_Return_One_Element_When_Called_With_One_In_Slice(t *testing.T) {
	// PREPARE DATA FOR TEST
	expected := dto.NewElement(1, 2, 3, 0.33)
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := dto.NewResult(0, common.TASK_SELECT, int32(expected.Size()), data)
	actual := parseResultsToElements([]*dto.Result{result})
	for _, elem := range actual {
		AssertEqual(expected, elem, t)
	}
}

func Test_ParseResultsToElements_Should_Return_All_Elements_From_Single_Input(t *testing.T) {
	// PREPARE DATA FOR TEST
	expected := dto.Dtos{dto.NewElement(1, 2, 3, 0.33), dto.NewElement(4, 5, 6, 0.66), dto.NewElement(7, 8, 9, 0.99)}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := dto.NewResult(0, common.TASK_SELECT, int32(expected.Size()), data)
	actual := parseResultsToElements([]*dto.Result{result})
	// ASSERTIONS

	if len(actual) != 3 {
		t.Error("Wrong data returned. Expected ", len(expected), " values")
	}

	for index, elem := range actual {
		AssertEqual(expected[index], elem, t)
	}
}

func Test_ParseResultsToElements_Should_Return_All_Elements_From_Input(t *testing.T) {
	// PREPARE DATA FOR TEST
	expected := dto.Dtos{
		dto.NewElement(1, 2, 3, 0.33), dto.NewElement(4, 5, 6, 0.66), dto.NewElement(7, 8, 9, 0.99),
	}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := dto.NewResult(0, common.TASK_SELECT, int32(expected.Size()), data)
	actual := parseResultsToElements([]*dto.Result{result, result, result, result})
	// ASSERTIONS

	expected = append(expected, expected...)
	expected = append(expected, expected...)

	if len(actual) != len(expected) {
		t.Error("Wrong data returned. Expected ", len(expected), " values", " but got ", len(actual))
	}

	for index, elem := range actual {
		AssertEqual(expected[index], elem, t)
	}
}

func Test_ParseResultsToInfos_Should_Return_Empty_Slice_For_Nil_Imput(t *testing.T) {
	actual := parseResultsToInfos(nil)
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_ParseResultsToInfos_Should_Return_Empty_Slice_For_Empty_Imput(t *testing.T) {
	actual := parseResultsToInfos([]*dto.Result{})
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_ParseResultsToInfos_Should_Return_Empty_Slice_For_Empty_Data_Imput(t *testing.T) {
	data := []byte{}
	result := dto.NewResult(0, common.TASK_INFO, 0, data)
	actual := parseResultsToInfos([]*dto.Result{result, result, result})
	if len(actual) != 0 {
		t.Error("Expected empty slice but got ", actual)
	}
}

func Test_ParseResultsToInfos_Should_Return_One_Element_When_Called_With_One_In_Slice(t *testing.T) {
	// PREPARE DATA FOR TEST
	expected := &dto.MemoryInfo{1, 1, 1, 1, 1}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := dto.NewResult(0, common.TASK_INFO, int32(expected.Size()), data)
	actual := parseResultsToInfos([]*dto.Result{result})
	for _, elem := range actual {
		AssertEqual(expected, elem, t)
	}
}

func Test_ParseResultsToInfos_Should_Return_All_Elements_From_Single_Input(t *testing.T) {
	// PREPARE DATA FOR TEST
	expected := dto.Dtos{
		&dto.MemoryInfo{1, 1, 1, 1, 1}, &dto.MemoryInfo{1, 1, 1, 1, 1},
		&dto.MemoryInfo{2, 1, 2, 1, 1}, &dto.MemoryInfo{3, 1, 3, 1, 1}}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := dto.NewResult(0, common.TASK_INFO, int32(expected.Size()), data)
	actual := parseResultsToInfos([]*dto.Result{result})
	// ASSERTIONS

	if len(actual) != len(expected) {
		t.Error("Wrong data returned. Expected ", len(expected), " values", " but got ", len(actual))
	}

	for index, elem := range actual {
		AssertEqual(expected[index], elem, t)
	}
}

func Test_ParseResultsToInfos_Should_Return_All_Elements_From_Input(t *testing.T) {
	// PREPARE DATA FOR TEST
	expected := dto.Dtos{
		&dto.MemoryInfo{1, 1, 1, 1, 1}, &dto.MemoryInfo{3, 2, 1, 1, 1},
	}
	data, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred", err)
	}
	result := dto.NewResult(0, common.TASK_INFO, int32(expected.Size()), data)
	actual := parseResultsToInfos([]*dto.Result{result, result, result, result})
	// ASSERTIONS

	expected = append(expected, expected...)
	expected = append(expected, expected...)

	if len(actual) != len(expected) {
		t.Error("Wrong data returned. Expected ", len(expected), " values", " but got ", len(actual))
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
