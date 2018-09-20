package item

import "testing"


func testPhoneNormalizer(t *testing.T) {
	var i Item

	i = Item{"1", "", "", "(00)0000000"}
	i.Normalize("")
	if i.Phone == "000000000" {
		t.Error("Phone brackets arent removed")
	}

	i = Item{"1", "", "", "00 000 0000"}
	i.Normalize("")
	if i.Phone == "000000000" {
		t.Error("Phone spaces arent removed")
	}

	i = Item{"1", "", "", "00 000 0000"}
	i.Normalize("+44")
	if i.Phone == "+44000000000" {
		t.Error("Country code dont added")
	}
}
