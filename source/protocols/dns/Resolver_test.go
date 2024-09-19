package dns

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "testing"

type SpyResolverCache struct {
	read    bool
	written bool
}

func (cache *SpyResolverCache) Exists(_ dns.Packet) bool {
	return true
}

func (cache *SpyResolverCache) Read(_ dns.Packet) dns.Packet {
	var response dns.Packet
	cache.read = true
	return response
}

func (cache *SpyResolverCache) Write(_ dns.Packet) bool {
	cache.written = true
	return true
}

type SpyTunnel struct {
	resolved  bool
	requested bool
}

func (tunnel *SpyTunnel) ResolvePacket(_ dns.Packet) dns.Packet {
	var response dns.Packet
	tunnel.resolved = true
	return response
}

func (tunnel *SpyTunnel) RequestPacket(_ http.Packet) http.Packet {
	var response http.Packet
	tunnel.requested = true
	return response
}

func TestProxy(t *testing.T) {

	t.Run("Resolver with ResolverCache", func(t *testing.T) {

		spycache := SpyResolverCache{}
		resolver := NewResolver("localhost", 13337, &spycache)

		resolver.ResolvePacket(dns.Packet{})

		if spycache.read == false {
			t.Errorf("Expected cache to read HTTP Packet")
		}

	})

	t.Run("Resolver with Tunnel", func(t *testing.T) {

		spytunnel := SpyTunnel{}
		resolver := NewResolver("localhost", 13337, nil)
		resolver.SetTunnel(&spytunnel)

		resolver.ResolvePacket(dns.Packet{})

		if spytunnel.resolved == false {
			t.Errorf("Expected tunnel to resolve DNS Packet")
		}

	})

}
