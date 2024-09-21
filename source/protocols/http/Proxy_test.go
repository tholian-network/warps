package http

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "testing"

type SpyProxyCache struct {
	read    bool
	written bool
}

func (cache *SpyProxyCache) Exists(_ http.Packet) bool {
	return true
}

func (cache *SpyProxyCache) Read(_ http.Packet) http.Packet {
	var response http.Packet
	cache.read = true
	return response
}

func (cache *SpyProxyCache) Write(_ http.Packet) bool {
	cache.written = true
	return true
}

func TestProxy(t *testing.T) {

	t.Run("Proxy with ProxyCache", func(t *testing.T) {

		spycache := SpyProxyCache{}
		proxy := NewProxy("localhost", 13337, &spycache)

		proxy.RequestPacket(http.Packet{})

		if spycache.read == false {
			t.Errorf("Expected cache to read HTTP Packet")
		}

	})

	t.Run("Proxy with Resolver", func(t *testing.T) {

		spyresolver := SpyResolver{}
		proxy := NewProxy("localhost", 13337, nil)
		proxy.SetResolver(&spyresolver)

		proxy.ResolvePacket(dns.Packet{})

		if spyresolver.resolved == false {
			t.Errorf("Expected resolver to resolve DNS Packet")
		}

	})

	t.Run("Proxy with Tunnel", func(t *testing.T) {

		spytunnel := SpyTunnel{}
		proxy := NewProxy("localhost", 13337, nil)
		proxy.SetTunnel(&spytunnel)

		proxy.ResolvePacket(dns.Packet{})

		if spytunnel.resolved == false {
			t.Errorf("Expected tunnel to resolve DNS Packet")
		}

		proxy.RequestPacket(http.Packet{})

		if spytunnel.requested == false {
			t.Errorf("Expected tunnel to request HTTP Packet")
		}

	})

}
