package dnstunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/protocols/test"
import net_url "net/url"
import "strconv"
import "testing"
import "time"

func TestTunnel(t *testing.T) {

	t.Run("Tunnel with ResolverCache and DNS Payload", func(t *testing.T) {

		record := dns.NewRecord("example.com", dns.TypeA)
		record.SetIPv4("1.3.3.7")

		expected := dns.NewPacket()
		expected.SetType("response")
		expected.SetResponseCode(dns.ResponseCodeNoError)
		expected.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))
		expected.AddAnswer(record)

		cache := test.NewSpyResolverCache(true, &expected, false)
		resolver := NewResolver("localhost", 13337, &cache)
		proxy := NewProxy("localhost", 13337, nil)
		proxy.SetResolver(&resolver)

		tunnel := NewTunnel("127.0.0.1", 13337)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		go func() {
			proxy.Listen()
		}()

		go func() {
			time.Sleep(2 * time.Second)
			proxy.Destroy()
		}()

		time.Sleep(100 * time.Millisecond)
		response := tunnel.ResolvePacket(query)

		if len(response.Answers) == 1 && response.Answers[0].Type == dns.TypeA {
			// success
		} else {
			t.Errorf("Expected DNS response to have a Type A record")
		}

		time.Sleep(2 * time.Second)

	})

	t.Run("Tunnel with ProxyCache and Large HTTP Payload", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.txt")
		expected := http.NewPacket()
		expected.SetURL(*url)
		expected.SetStatus(http.StatusOK)
		expected.SetHeader("Content-Type", "text/plain")
		expected.SetHeader("X-Proxy", "SpyProxyCache")

		payload := make([]byte, 0)

		for l := 0; l < 100; l++ {
			payload = append(payload, []byte("Hello, line "+strconv.Itoa(l)+"!\n")...)
		}

		expected.SetPayload(payload)

		cache := test.NewSpyProxyCache(true, &expected, false)
		proxy := NewProxy("localhost", 13337, &cache)
		tunnel := NewTunnel("127.0.0.1", 13337)

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		go func() {
			proxy.Listen()
		}()

		go func() {
			time.Sleep(2 * time.Second)
			proxy.Destroy()
		}()

		time.Sleep(100 * time.Millisecond)
		response := tunnel.RequestPacket(request)
		response.Decode()

		if expected.Status != response.Status {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(expected.Status).String(), http.Status(response.Status).String())
		}

		if string(expected.Payload) != string(response.Payload) {
			t.Errorf("Expected different HTTP response payload")
		}

		time.Sleep(2 * time.Second)

	})

}
