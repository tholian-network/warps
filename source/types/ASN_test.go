package types

import "testing"

func TestASN(t *testing.T) {

	t.Run("IsASN()", func(t *testing.T) {

		got1 := IsASN("AS1234")
		got2 := IsASN("1234")
		got3 := IsASN("ASABC")

		if got1 != true {
			t.Errorf("Expected %t to be true", got1)
		}

		if got2 != false {
			t.Errorf("Expected %t to be false", got2)
		}

		if got3 != false {
			t.Errorf("Expected %t to be false", got3)
		}

	})

	t.Run("ParseASN()", func(t *testing.T) {

		asn1 := ParseASN("AS1337")
		asn2 := ParseASN("1337")
		asn3 := ParseASN("AS")
		asn4 := ParseASN("ABC")

		if asn1 == nil {
			t.Errorf("Expected %s to be valid", "AS1337")
		} else if asn1.String() != "AS1337" {
			t.Errorf("Expected %s to be %s", asn1.String(), "AS1337")
		}

		if asn2 == nil {
			t.Errorf("Expected %s to be valid", "AS1337")
		} else if asn2.String() != "AS1337" {
			t.Errorf("Expected %s to be %s", asn2.String(), "AS1337")
		}

		if asn3 != nil {
			t.Errorf("Expected %s to be invalid", "AS")
		}

		if asn4 != nil {
			t.Errorf("Expected %s to be invalid", "ABC")
		}

	})

	t.Run("ToASN()", func(t *testing.T) {

		asn1 := ToASN("AS1337")
		asn2 := ToASN("1337")
		asn3 := ToASN("AS")
		asn4 := ToASN("ABC")

		if asn1.String() != "AS1337" {
			t.Errorf("Expected %s to be %s", asn1.String(), "AS1337")
		}

		if asn2.String() != "AS1337" {
			t.Errorf("Expected %s to be %s", asn2.String(), "AS1337")
		}

		if asn3.String() != "AS0" {
			t.Errorf("Expected %s to be %s", asn3.String(), "AS0")
		}

		if asn4.String() != "AS0" {
			t.Errorf("Expected %s to be %s", asn4.String(), "AS0")
		}

	})

}
