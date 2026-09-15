package dns

import "testing"

func TestType(t *testing.T) {

	t.Run("String(A)", func(t *testing.T) {
		if TypeA.String() != "A" {
			t.Errorf("Expected '%s' but got '%s'", "A", TypeA.String())
		}
	})

	t.Run("String(NS)", func(t *testing.T) {
		if TypeNS.String() != "NS" {
			t.Errorf("Expected '%s' but got '%s'", "NS", TypeNS.String())
		}
	})

	t.Run("String(CNAME)", func(t *testing.T) {
		if TypeCNAME.String() != "CNAME" {
			t.Errorf("Expected '%s' but got '%s'", "CNAME", TypeCNAME.String())
		}
	})

	t.Run("String(SOA)", func(t *testing.T) {
		if TypeSOA.String() != "SOA" {
			t.Errorf("Expected '%s' but got '%s'", "SOA", TypeSOA.String())
		}
	})

	t.Run("String(PTR)", func(t *testing.T) {
		if TypePTR.String() != "PTR" {
			t.Errorf("Expected '%s' but got '%s'", "PTR", TypePTR.String())
		}
	})

	t.Run("String(MX)", func(t *testing.T) {
		if TypeMX.String() != "MX" {
			t.Errorf("Expected '%s' but got '%s'", "MX", TypeMX.String())
		}
	})

	t.Run("String(TXT)", func(t *testing.T) {
		if TypeTXT.String() != "TXT" {
			t.Errorf("Expected '%s' but got '%s'", "TXT", TypeTXT.String())
		}
	})

	t.Run("String(AAAA)", func(t *testing.T) {
		if TypeAAAA.String() != "AAAA" {
			t.Errorf("Expected '%s' but got '%s'", "AAAA", TypeAAAA.String())
		}
	})

	t.Run("String(SRV)", func(t *testing.T) {
		if TypeSRV.String() != "SRV" {
			t.Errorf("Expected '%s' but got '%s'", "SRV", TypeSRV.String())
		}
	})

	t.Run("String(URI)", func(t *testing.T) {
		if TypeURI.String() != "URI" {
			t.Errorf("Expected '%s' but got '%s'", "URI", TypeURI.String())
		}
	})

	t.Run("String(unknown)", func(t *testing.T) {
		if Type(100).String() != "" {
			t.Errorf("Expected '%s' but got '%s'", "", Type(100).String())
		}
	})

}
