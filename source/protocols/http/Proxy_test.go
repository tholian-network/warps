package http

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-warps/protocols/test"
import net_url "net/url"
import "strings"
import "testing"

func TestProxy(t *testing.T) {

	t.Run("Proxy with ProxyCache", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.html")
		expected := http.NewPacket()
		expected.SetURL(*url)
		expected.SetStatus(http.StatusOK)
		expected.SetHeader("Content-Type", "text/html")
		expected.SetHeader("X-Proxy", "SpyProxyCache")
		expected.SetPayload([]byte("Hello, world!"))

		cache := test.NewSpyProxyCache(true, &expected, false)
		proxy := NewProxy("localhost", 13337, &cache)

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		response := proxy.RequestPacket(request)
		response.Decode()

		if cache.RequestedRead != url.String() {
			t.Errorf("Expected ProxyCache to Read '%s' but got '%s'", url.String(), cache.RequestedRead)
		}

		if response.Status != http.StatusOK {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(http.StatusOK).String(), http.Status(response.Status).String())
		}

		if string(expected.Payload) != string(response.Payload) {
			t.Errorf("Expected different HTTP response payload")
		}

	})

	t.Run("Proxy with Internet DNS Payload", func(t *testing.T) {

		proxy := NewProxy("localhost", 13337, nil)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		response := proxy.ResolvePacket(query)

		if response.Codes.Response != dns.ResponseCodeNoError {
			t.Errorf("Expected DNS response code '%s' but got '%s'", dns.ResponseCode(dns.ResponseCodeNoError).String(), dns.ResponseCode(response.Codes.Response).String())
		}

	})

	t.Run("Proxy with Internet HTTP Payload", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.html")
		cache := test.NewSpyProxyCache(false, nil, false)
		proxy := NewProxy("localhost", 13337, &cache)

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		response := proxy.RequestPacket(request)
		response.Decode()

		if cache.RequestedExists != "http://example.com/index.html" {
			t.Errorf("Expected SpyCache to read URL '%s' but got '%s'", "http://example.com/index.html", cache.RequestedExists)
		}

		if response.Status != http.StatusOK {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(http.StatusOK).String(), http.Status(response.Status).String())
		}

		html := string(response.Payload)

		if !strings.Contains(html, "<title>Example Domain</title>") {
			t.Errorf("Expected different HTTP response payload")
		}

	})

	t.Run("Proxy with Resolver", func(t *testing.T) {

		resolver := test.NewSpyResolver(false)
		proxy := NewProxy("localhost", 13337, nil)
		proxy.SetResolver(&resolver)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		proxy.ResolvePacket(query)

		if resolver.Resolved != "A:example.com" {
			t.Errorf("Expected Resolver to resolve '%s' but got '%s'", "A:example.com", resolver.Resolved)
		}

	})

	t.Run("Proxy with Tunnel", func(t *testing.T) {

		tunnel := test.NewSpyTunnel(false)
		proxy := NewProxy("localhost", 13337, nil)
		proxy.SetTunnel(&tunnel)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		proxy.ResolvePacket(query)

		if tunnel.Resolved != "A:example.com" {
			t.Errorf("Expected Tunnel to resolve '%s' but got '%s'", "A:example.com", tunnel.Resolved)
		}

		url, _ := net_url.Parse("http://example.com/index.html")
		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		proxy.RequestPacket(request)

		if tunnel.Requested != url.String() {
			t.Errorf("Expected Tunnel to request '%s' but got '%s'", url.String(), tunnel.Requested)
		}

	})

}
