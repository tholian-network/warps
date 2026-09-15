package dns

import "testing"

func TestClass(t *testing.T) {

	t.Run("String(Internet)", func(t *testing.T) {
		if ClassInternet.String() != "Internet" {
			t.Errorf("Expected '%s' but got '%s'", "Internet", ClassInternet.String())
		}
	})

	t.Run("String(Chaosnet)", func(t *testing.T) {
		if ClassChaosnet.String() != "Chaosnet" {
			t.Errorf("Expected '%s' but got '%s'", "Chaosnet", ClassChaosnet.String())
		}
	})

	t.Run("String(Hesoidnet)", func(t *testing.T) {
		if ClassHesoidnet.String() != "Hesoidnet" {
			t.Errorf("Expected '%s' but got '%s'", "Hesoidnet", ClassHesoidnet.String())
		}
	})

	t.Run("String(None)", func(t *testing.T) {
		if ClassNone.String() != "" {
			t.Errorf("Expected '%s' but got '%s'", "", ClassNone.String())
		}
	})

	t.Run("String(unknown)", func(t *testing.T) {
		if Class(100).String() != "" {
			t.Errorf("Expected '%s' but got '%s'", "", Class(100).String())
		}
	})

}
