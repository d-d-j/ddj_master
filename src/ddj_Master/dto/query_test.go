package dto

import "testing"

func Test_EncodeDecode_Query(t *testing.T) {

	t.Skip("Method not implemented")

	expected := Query{1, []int32{1}, 2, []int32{0, 1}, 4, []int64{5, 7, 11, 21}, 0, nil}
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	var actual Query
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Query_Size(t *testing.T) {
	actual := (&Query{1, []int32{1}, 2, []int32{0, 1}, 4, []int64{5, 7, 11, 21}, 0, nil}).Size()
	expected := (7*32 + 4*64) / 8
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Query_Size_For_Empty_Query(t *testing.T) {
	actual := (&Query{}).Size()
	expected := (4 * 32) / 8
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Encode_Query_With_Additional_Data(t *testing.T) {

	input := Query{1, []int32{1}, 2, []int32{0, 1}, 2, []int64{5, 7, 11, 21}, 10, HistogramTimeData{5, 10, 15}}
	expected := []byte{
		1, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 7,
		0, 0, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 21, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 5, 0, 0, 0, 0,
		0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0,
	}
	actual, err := input.Encode()

	if err != nil {
		t.Error("Error occurred")
	}
	if input.Size() != len(expected) {
		t.Error("Size doesn't match with data. Expected ", len(expected), " but got ", input.Size())
	}

	for i, _ := range expected {
		if expected[i] != actual[i] {
			t.Error("Expected ", expected, " but got ", actual)
			break
		}
	}
}
