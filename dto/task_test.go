package dto

import (
	"github.com/d-d-j/ddj_master/common"
	"testing"
)

func Test_NewTask_Should_Fill_AggregationType_When_Query_Is_Provided_As_Data(t *testing.T) {
	expected := common.AGGREGATION_ADD
	query := &Query{1, []int32{1}, 2, []int32{0, 1}, 4, []int64{5, 7, 11, 21}, expected, nil}
	request := RestRequest{common.TASK_SELECT, query, nil}
	actual := NewTask(1, request, nil).AggregationType

	if actual != expected {
		t.Error("Expected ", expected, " but got ", query)
	}
}

func Test_NewTask_Should_Fill_AggregationType_With_Uninitialized_When_Query_Is_NOT_Provided_As_Data(t *testing.T) {
	expected := int32(common.CONST_UNINITIALIZED)
	data := NewElement(0, 1, 2, 0.3)
	request := RestRequest{common.TASK_INSERT, data, nil}
	actual := NewTask(1, request, nil).AggregationType

	if actual != expected {
		t.Error("Expected ", expected, " but got ", actual)
	}
}

func Test_NewTask_Should_Panic_When_TaskType_Is_Select_But_Data_Is_NOT_Query(t *testing.T) {
	data := NewElement(0, 1, 2, 0.3)
	request := RestRequest{common.TASK_SELECT, data, nil}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()

	NewTask(1, request, nil)

}
