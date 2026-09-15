package dns

import "testing"

func TestQuestion(t *testing.T) {

	t.Run("NewQuestion with Type A", func(t *testing.T) {

		question := NewQuestion("example.com", TypeA)

		if question.Name != "example.com" {
			t.Errorf("Expected name '%s' but got '%s'", "example.com", question.Name)
		}

		if question.Type != TypeA {
			t.Errorf("Expected type '%s' but got '%s'", "A", question.Type.String())
		}

		if question.Class != ClassInternet {
			t.Errorf("Expected class '%s' but got '%s'", "Internet", question.Class.String())
		}

	})

	t.Run("NewQuestion with Type PTR", func(t *testing.T) {

		question := NewQuestion("1.3.3.7", TypePTR)

		if question.Name != "7.3.3.1.in-addr.arpa" {
			t.Errorf("Expected name '%s' but got '%s'", "7.3.3.1.in-addr.arpa", question.Name)
		}

		if question.Type != TypePTR {
			t.Errorf("Expected type '%s' but got '%s'", "PTR", question.Type.String())
		}

	})

	t.Run("SetName with long label", func(t *testing.T) {

		question := NewQuestion("", TypeA)

		label := ""

		for l := 0; l < 64; l++ {
			label += "a"
		}

		question.SetName(label + ".example.com")

		if question.Name != "" {
			t.Errorf("Expected name to be empty but got '%s'", question.Name)
		}

	})

	t.Run("Bytes", func(t *testing.T) {

		question := NewQuestion("example.com", TypeA)
		bytes := question.Bytes()

		if len(bytes) > 0 {
			// 7 example 3 com 0 = 1+7+1+3+1 = 13 bytes for labels
			// plus 2 bytes type + 2 bytes class = 17 bytes
			if len(bytes) != 17 {
				t.Errorf("Expected '%d' bytes but got '%d'", 17, len(bytes))
			}
		} else {
			t.Errorf("Expected question bytes to be non-empty")
		}

	})

}
