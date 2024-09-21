package url

import net_url "net/url"
import "testing"

func TestToHostAndPort(t *testing.T) {

	t.Run("Host", func(t *testing.T) {

		url, _ := net_url.Parse("http://sub.domain.example.com/index.html")
		expected := "sub.domain.example.com:80"
		host := ToHostAndPort(url)

		if host != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, host)
		}

	})

	t.Run("Host and Port", func(t *testing.T) {

		url, _ := net_url.Parse("http://sub.domain.example.com:1337/index.html")
		expected := "sub.domain.example.com:1337"
		host := ToHostAndPort(url)

		if host != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, host)
		}

	})

	t.Run("IPv4", func(t *testing.T) {

		url, _ := net_url.Parse("https://1.33.33.7/index.html")
		expected := "1.33.33.7:443"
		host := ToHostAndPort(url)

		if host != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, host)
		}

	})

	t.Run("IPv4 and Port", func(t *testing.T) {

		url, _ := net_url.Parse("https://1.33.33.7:1337/index.html")
		expected := "1.33.33.7:1337"
		host := ToHostAndPort(url)

		if host != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, host)
		}

	})

	t.Run("IPv6", func(t *testing.T) {

		url, _ := net_url.Parse("http://[fe80::1337]/index.html")
		expected := "[fe80:0000:0000:0000:0000:0000:0000:1337]:80"
		host := ToHostAndPort(url)

		if host != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, host)
		}

	})

	t.Run("IPv6 and Port", func(t *testing.T) {

		url, _ := net_url.Parse("http://[fe80::1337]:1337/index.html")
		expected := "[fe80:0000:0000:0000:0000:0000:0000:1337]:1337"
		host := ToHostAndPort(url)

		if host != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, host)
		}

	})

}
