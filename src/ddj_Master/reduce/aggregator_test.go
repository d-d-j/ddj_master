package reduce

import (
	"testing"
)

func Test_NonAggregation_Should_Return_Nil_When_Input_Is_Nil(t *testing.T) {
	var nonAggr NonAggregation
	actual := nonAggr.Aggregate(nil)
	if actual != nil {
		t.Error("Expected nil but got ", actual)
	}
}
