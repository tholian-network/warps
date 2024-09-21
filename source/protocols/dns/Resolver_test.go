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

	t.Run("Resolver with DNS Type A Payload", func(t *testing.T) {

		resolver := NewResolver("localhost", 13337, nil)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		response := resolver.ResolvePacket(query)

		if len(query.Questions) == len(response.Questions) && len(response.Answers) > 0 {

			record := response.Answers[0]

			if record.Type != dns.TypeA {
				t.Errorf("Expected DNS response record type to be '%s' but got '%s'", dns.TypeA.String(), dns.Type(record.Type).String())
			}

		} else {
			t.Errorf("Expected DNS response questions to be '%d' but got '%d'", len(query.Questions), len(response.Questions))
		}

	})

	t.Run("Resolver with DNS Type AAAA Payload", func(t *testing.T) {

		resolver := NewResolver("localhost", 13337, nil)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeAAAA))

		response := resolver.ResolvePacket(query)

		if len(query.Questions) == len(response.Questions) && len(response.Answers) > 0 {

			record := response.Answers[0]

			if record.Type != dns.TypeAAAA {
				t.Errorf("Expected DNS response record type to be '%s' but got '%s'", dns.TypeAAAA.String(), dns.Type(record.Type).String())
			}

		} else {
			t.Errorf("Expected DNS response questions to be '%d' but got '%d'", len(query.Questions), len(response.Questions))
		}

	})

}
