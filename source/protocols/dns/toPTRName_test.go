package dns

import "testing"

func TestToPTRName(t *testing.T) {

	t.Run("IPv4", func(t *testing.T) {

		expected := "7.3.3.1.in-addr.arpa"
		name := toPTRName("1.3.3.7")

		if name != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

	t.Run("Domain", func(t *testing.T) {

		expected := "example.com"
		name := toPTRName("example.com")

		if name != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

}
