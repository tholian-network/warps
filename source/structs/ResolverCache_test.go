package structs

import "tholian-endpoint/protocols/dns"
import "bytes"
import "os"
import "testing"

func TestResolverCache(t *testing.T) {

	t.Run("Exists with single question", func(t *testing.T) {

		tmp, _ := os.MkdirTemp("", "tholian-warps-resolvercache-*")
		cache := NewResolverCache(tmp)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		result1 := cache.Exists(query)

		if result1 != false {
			t.Errorf("Expected '%t' but got '%t'", false, result1)
		}

		record := dns.NewRecord("example.com", dns.TypeA)
		record.SetIPv4("1.3.3.7")

		response := dns.NewPacket()
		response.SetType("response")
		response.SetResponseCode(dns.ResponseCodeNoError)
		response.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		response.AddAnswer(record)

		result2 := cache.Write(response)

		if result2 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result2)
		}

		result3 := cache.Exists(query)

		if result3 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result3)
		}

		os.RemoveAll(tmp)

	})

	t.Run("Exists with multiple questions", func(t *testing.T) {

		tmp, _ := os.MkdirTemp("", "tholian-warps-resolvercache-*")
		cache := NewResolverCache(tmp)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeAAAA))

		result1 := cache.Exists(query)

		if result1 != false {
			t.Errorf("Expected '%t' but got '%t'", false, result1)
		}

		record1 := dns.NewRecord("example.com", dns.TypeA)
		record1.SetIPv4("1.3.3.7")
		response1 := dns.NewPacket()
		response1.SetType("response")
		response1.SetResponseCode(dns.ResponseCodeNoError)
		response1.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		response1.AddAnswer(record1)

		result2 := cache.Write(response1)

		if result2 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result2)
		}

		result3 := cache.Exists(query)

		if result3 != false {
			t.Errorf("Expected '%t' but got '%t'", false, result3)
		}

		record21 := dns.NewRecord("example.com", dns.TypeA)
		record21.SetIPv4("1.3.3.7")
		record22 := dns.NewRecord("example.com", dns.TypeAAAA)
		record22.SetIPv6("fe80::1337")
		response2 := dns.NewPacket()
		response2.SetType("response")
		response2.SetResponseCode(dns.ResponseCodeNoError)
		response2.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		response2.AddQuestion(dns.NewQuestion("example.com", dns.TypeAAAA))
		response2.AddAnswer(record21)
		response2.AddAnswer(record22)

		result4 := cache.Write(response2)

		if result4 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result4)
		}

		result5 := cache.Exists(query)

		if result5 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result5)
		}

		os.RemoveAll(tmp)

	})

	t.Run("Read/Write with single question", func(t *testing.T) {

		tmp, _ := os.MkdirTemp("", "tholian-warps-resolvercache-*")
		cache := NewResolverCache(tmp)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		result1 := cache.Read(query)

		if result1.Codes.Response != dns.ResponseCodeNonExistDomain {
			t.Errorf("Expected '%s' but got '%s'", dns.ResponseCode(dns.ResponseCodeNonExistDomain).String(), result1.Codes.Response.String())
		}

		if len(result1.Answers) != 0 {
			t.Errorf("Expected DNS Packet to have no Records")
		}

		record1 := dns.NewRecord("example.com", dns.TypeA)
		record1.SetIPv4("1.3.3.7")
		response1 := dns.NewPacket()
		response1.SetType("response")
		response1.SetResponseCode(dns.ResponseCodeNoError)
		response1.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		response1.AddAnswer(record1)

		result2 := cache.Write(response1)

		if result2 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result2)
		}

		result3 := cache.Read(query)

		if result3.Codes.Response != dns.ResponseCodeNoError {
			t.Errorf("Expected '%s' but got '%s'", dns.ResponseCode(dns.ResponseCodeNoError).String(), result3.Codes.Response.String())
		}

		if len(result3.Answers) == 1 {

			if result3.Answers[0].Type != dns.TypeA {
				t.Errorf("Expected DNS Packet to have dns.TypeA Record")
			}

			if bytes.Equal(result3.Answers[0].Data, record1.Data) != true {
				t.Errorf("Expected DNS Packet to have 1.3.3.7 as Record Data")
			}


		} else {
			t.Errorf("Expected DNS Packet to have Records")
		}

		os.RemoveAll(tmp)

	})

	t.Run("Read/Write with multiple questions", func(t *testing.T) {

		tmp, _ := os.MkdirTemp("", "tholian-warps-resolvercache-*")
		cache := NewResolverCache(tmp)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeAAAA))

		result1 := cache.Read(query)

		if result1.Codes.Response != dns.ResponseCodeNonExistDomain {
			t.Errorf("Expected '%s' but got '%s'", dns.ResponseCode(dns.ResponseCodeNonExistDomain).String(), result1.Codes.Response.String())
		}

		if len(result1.Answers) != 0 {
			t.Errorf("Expected DNS Packet to have no Records")
		}

		record1 := dns.NewRecord("example.com", dns.TypeA)
		record1.SetIPv4("1.3.3.7")
		response1 := dns.NewPacket()
		response1.SetType("response")
		response1.SetResponseCode(dns.ResponseCodeNoError)
		response1.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		response1.AddAnswer(record1)

		result2 := cache.Write(response1)

		if result2 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result2)
		}

		result3 := cache.Read(query)

		if result3.Codes.Response != dns.ResponseCodeNoError {
			t.Errorf("Expected '%s' but got '%s'", dns.ResponseCode(dns.ResponseCodeNoError).String(), result3.Codes.Response.String())
		}

		if len(result3.Answers) != 1 {
			t.Errorf("Expected DNS Packet to have 1 Record")
		}

		record2 := dns.NewRecord("example.com", dns.TypeAAAA)
		record2.SetIPv6("fe80::1337")
		response2 := dns.NewPacket()
		response2.SetType("response")
		response2.SetResponseCode(dns.ResponseCodeNoError)
		response2.AddQuestion(dns.NewQuestion("example.com", dns.TypeAAAA))
		response2.AddAnswer(record2)

		result4 := cache.Write(response2)

		if result4 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result4)
		}

		result5 := cache.Read(query)

		if result5.Codes.Response != dns.ResponseCodeNoError {
			t.Errorf("Expected '%s' but got '%s'", dns.ResponseCode(dns.ResponseCodeNoError).String(), result5.Codes.Response.String())
		}

		if len(result5.Answers) != 2 {
			t.Errorf("Expected DNS Packet to have 2 Records")
		}

	})

}
