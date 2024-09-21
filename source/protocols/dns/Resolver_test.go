package dns

import "tholian-endpoint/protocols/dns"
import "tholian-warps/protocols/test"
import "testing"

func TestResolver(t *testing.T) {

	t.Run("Resolver with ResolverCache", func(t *testing.T) {

		record := dns.NewRecord("example.com", dns.TypeA)
		record.SetIPv4("1.3.3.7")
		response := dns.NewPacket()
		response.SetType("response")
		response.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		cache := test.NewSpyResolverCache(true, &response, true)
		resolver := NewResolver("localhost", 13337, &cache)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		resolver.ResolvePacket(query)

		if cache.ResolvedRead == "A:example.com" {
			t.Errorf("Expected cache to read HTTP Packet")
		}

	})

	t.Run("Resolver with Tunnel", func(t *testing.T) {

		tunnel := test.NewSpyTunnel(false)
		resolver := NewResolver("localhost", 13337, nil)
		resolver.SetTunnel(&tunnel)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		resolver.ResolvePacket(query)

		if tunnel.Resolved == "A:example.com" {
			t.Errorf("Expected tunnel to resolve DNS Packet")
		}

	})

}
