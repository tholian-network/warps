package dns

import "testing"

func TestRecord(t *testing.T) {

	t.Run("NewRecord with Type A", func(t *testing.T) {

		record := NewRecord("example.com", TypeA)

		if record.Name != "example.com" {
			t.Errorf("Expected name '%s' but got '%s'", "example.com", record.Name)
		}

		if record.Type != TypeA {
			t.Errorf("Expected type '%s' but got '%s'", "A", record.Type.String())
		}

		if record.ToIPv4() != "0.0.0.0" {
			t.Errorf("Expected '%s' but got '%s'", "0.0.0.0", record.ToIPv4())
		}

	})

	t.Run("NewRecord with Type AAAA", func(t *testing.T) {

		record := NewRecord("example.com", TypeAAAA)

		if record.ToIPv6() != "[0000:0000:0000:0000:0000:0000:0000:0000]" {
			t.Errorf("Expected '%s' but got '%s'", "[0000:0000:0000:0000:0000:0000:0000:0000]", record.ToIPv6())
		}

	})

	t.Run("NewRecord with Type URI", func(t *testing.T) {

		record := NewRecord("example.com", TypeURI)

		if record.ToURL() != "http://localhost/index.html" {
			t.Errorf("Expected '%s' but got '%s'", "http://localhost/index.html", record.ToURL())
		}

	})

	t.Run("SetIPv4 and ToIPv4", func(t *testing.T) {

		record := NewRecord("example.com", TypeA)
		record.SetIPv4("1.3.3.7")

		if record.ToIPv4() != "1.3.3.7" {
			t.Errorf("Expected '%s' but got '%s'", "1.3.3.7", record.ToIPv4())
		}

	})

	t.Run("SetIPv6 and ToIPv6", func(t *testing.T) {

		record := NewRecord("example.com", TypeAAAA)
		record.SetIPv6("fe80::1337")

		if record.ToIPv6() != "[fe80:0000:0000:0000:0000:0000:0000:1337]" {
			t.Errorf("Expected '%s' but got '%s'", "[fe80:0000:0000:0000:0000:0000:0000:1337]", record.ToIPv6())
		}

	})

	t.Run("SetURL and ToURL", func(t *testing.T) {

		record := NewRecord("example.com", TypeURI)
		record.SetURL("http://example.com/index.html")

		if record.ToURL() != "http://example.com/index.html" {
			t.Errorf("Expected '%s' but got '%s'", "http://example.com/index.html", record.ToURL())
		}

	})

	t.Run("SetDomain and ToDomain", func(t *testing.T) {

		record := NewRecord("example.com", TypeCNAME)
		record.SetDomain("target.example.com")

		if record.ToDomain() != "target.example.com" {
			t.Errorf("Expected '%s' but got '%s'", "target.example.com", record.ToDomain())
		}

	})

	t.Run("SetPort on Type SRV", func(t *testing.T) {

		record := NewRecord("example.com", TypeSRV)
		record.SetPort(8080)

		if record.ToPort() != 8080 {
			t.Errorf("Expected '%d' but got '%d'", 8080, record.ToPort())
		}

	})

	t.Run("SetTTL", func(t *testing.T) {

		record := NewRecord("example.com", TypeA)
		record.SetTTL(3600)

		if record.TTL != 3600 {
			t.Errorf("Expected '%d' but got '%d'", 3600, record.TTL)
		}

	})

}
