package icmptunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/test"
import "testing"

func TestResolver(t *testing.T) {

	t.Run("Resolver with ResolverCache", func(t *testing.T) {

		response := dns.NewPacket()
		response.SetType("response")
		response.SetResponseCode(dns.ResponseCodeNoError)
		response.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		cache := test.NewSpyResolverCache(true, &response, true)
		resolver := NewResolver("localhost", 13337, &cache)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		resolver.ResolvePacket(query)

		if cache.ResolvedExists != "A:example.com" {
			t.Errorf("Expected SpyResolverCache to Exists '%s' but got '%s'", "A:example.com", cache.ResolvedExists)
		}

		if cache.ResolvedRead != "A:example.com" {
			t.Errorf("Expected SpyResolverCache to Read '%s' but got '%s'", "A:example.com", cache.ResolvedRead)
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

		if tunnel.Resolved != "A:example.com" {
			t.Errorf("Expected SpyTunnel to resolve '%s' but got '%s'", "A:example.com", tunnel.Resolved)
		}

	})

}
