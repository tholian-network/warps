package https

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/protocols/test"
import net_url "net/url"
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
