package tunnel

import net_url "net/url"
import "testing"

func TestToRecordName(t *testing.T) {

	t.Run("Domain", func(t *testing.T) {

		url, _ := net_url.Parse("http://sub.domain.example.com/index.html")
		expected := "sub.domain.example.com"
		name := ToRecordName(url)

		if expected != name {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

	t.Run("Domain and Port", func(t *testing.T) {

		url, _ := net_url.Parse("http://sub.domain.example.com:1337/index.html")
		expected := "sub.domain.example.com"
		name := ToRecordName(url)

		if expected != name {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

	t.Run("IPv4", func(t *testing.T) {

		url, _ := net_url.Parse("http://1.3.3.7/index.html")
		expected := "7.3.3.1.in-addr.arpa"
		name := ToRecordName(url)

		if expected != name {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

	t.Run("IPv4 and Port", func(t *testing.T) {

		url, _ := net_url.Parse("http://1.3.3.7:1337/index.html")
		expected := "7.3.3.1.in-addr.arpa"
		name := ToRecordName(url)

		if expected != name {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

	t.Run("IPv6", func(t *testing.T) {

		url, _ := net_url.Parse("http://[fe80::1337]/index.html")
		expected := "7.3.3.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.e.f.ip6.arpa"
		name := ToRecordName(url)

		if expected != name {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

	t.Run("IPv6 and Port", func(t *testing.T) {

		url, _ := net_url.Parse("http://[fe80::1337]:1337/index.html")
		expected := "7.3.3.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.e.f.ip6.arpa"
		name := ToRecordName(url)

		if expected != name {
			t.Errorf("Expected '%s' but got '%s'", expected, name)
		}

	})

}
