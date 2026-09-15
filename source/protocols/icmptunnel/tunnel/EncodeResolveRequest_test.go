package tunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/icmp"
import "testing"

func TestEncodeResolveRequest(t *testing.T) {

	t.Run("DNS Query round-trip", func(t *testing.T) {

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		message := icmp.NewMessage()
		EncodeResolveRequest(&message, query)

		if IsResolveRequest(&message) != true {
			t.Errorf("Expected message to be a resolve request")
		}

		decoded := DecodeResolveRequest(&message)

		if decoded.Type != "query" {
			t.Errorf("Expected type '%s' but got '%s'", "query", decoded.Type)
		}

		if len(decoded.Questions) != 1 {
			t.Errorf("Expected '%d' questions but got '%d'", 1, len(decoded.Questions))
		}

		if decoded.Questions[0].Name != "example.com" {
			t.Errorf("Expected name '%s' but got '%s'", "example.com", decoded.Questions[0].Name)
		}

	})

	t.Run("DNS Response round-trip", func(t *testing.T) {

		record := dns.NewRecord("example.com", dns.TypeA)
		record.SetIPv4("1.3.3.7")

		response := dns.NewPacket()
		response.SetType("response")
		response.SetResponseCode(dns.ResponseCodeNoError)
		response.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		response.AddAnswer(record)

		message := icmp.NewMessage()
		EncodeResolveResponse(&message, response)

		if IsResolveResponse(&message) != true {
			t.Errorf("Expected message to be a resolve response")
		}

		decoded := DecodeResolveResponse(&message)

		if decoded.Type != "response" {
			t.Errorf("Expected type '%s' but got '%s'", "response", decoded.Type)
		}

		if len(decoded.Answers) != 1 {
			t.Errorf("Expected '%d' answers but got '%d'", 1, len(decoded.Answers))
		}

		if decoded.Answers[0].Type != dns.TypeA {
			t.Errorf("Expected answer type '%s' but got '%s'", dns.TypeA.String(), dns.Type(decoded.Answers[0].Type).String())
		}

	})

}
