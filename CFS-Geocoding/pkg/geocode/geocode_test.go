package geocode

import "testing"

func Test_NormalizeAddress_LeaveXInWord(t *testing.T) {
	expected := "3600 BlaXne"
	actual := NormalizeAddress("36XX BlaXne")
	if actual != expected {
		t.Errorf("Expected %v but received %v", expected, actual)
	}

}
