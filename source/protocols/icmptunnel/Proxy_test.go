package icmptunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/protocols/test"
import net_url "net/url"
import "testing"
import "time"

func TestProxy(t *testing.T) {

	t.Run("Proxy with ProxyCache and HTTP Payload", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.html")
		expected := http.NewPacket()
		expected.SetURL(*url)
		expected.SetStatus(http.StatusOK)
		expected.SetHeader("Content-Type", "text/html")
		expected.SetHeader("X-Proxy", "SpyProxyCache")
		expected.SetPayload([]byte("Hello, world!"))

		cache := test.NewSpyProxyCache(true, &expected, false)
		proxy := NewProxy("localhost", 13337, &cache)
		tunnel := NewTunnel("127.0.0.1", 13337)

		socket_a, socket_b := test.NewLoopbackSocket()

		proxy.SetSocket(socket_b)
		tunnel.SetSocket(socket_a)

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		go func() {
			err1 := proxy.Listen()
			if err1 != nil {
				t.Errorf("Unexpected error '%s'", err1.Error())
			}
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
			t.Errorf("Expected HTTP response payload '%s' but got '%s'", string(expected.Payload), string(response.Payload))
		}

		time.Sleep(2 * time.Second)

	})

	t.Run("Proxy with ResolverCache and DNS Payload", func(t *testing.T) {

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

		socket_a, socket_b := test.NewLoopbackSocket()

		proxy.SetSocket(socket_b)
		tunnel.SetSocket(socket_a)

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion("example.com", dns.TypeA))

		go func() {
			err1 := proxy.Listen()
			if err1 != nil {
				t.Errorf("Unexpected error '%s'", err1.Error())
			}
		}()

		go func() {
			time.Sleep(2 * time.Second)
			proxy.Destroy()
		}()

		time.Sleep(100 * time.Millisecond)
		response := tunnel.ResolvePacket(query)

		if len(response.Questions) > 0 && len(response.Answers) > 0 {

			if response.Answers[0].Type != dns.TypeA {
				t.Errorf("Expected DNS response record type to be '%s' but got '%s'", dns.TypeA.String(), dns.Type(response.Answers[0].Type).String())
			}

		} else {
			t.Errorf("Expected DNS response to have questions and answers")
		}

		time.Sleep(2 * time.Second)

	})

	t.Run("Proxy Chain with ProxyCache and HTTP Payload", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.html")
		expected := http.NewPacket()
		expected.SetURL(*url)
		expected.SetStatus(http.StatusOK)
		expected.SetHeader("X-Proxy", "SpyProxyCache")
		expected.SetPayload([]byte("Hello, world!"))

		cache := test.NewSpyProxyCache(true, &expected, false)
		gateway := NewProxy("localhost", 13337, &cache)
		forward := NewProxy("localhost", 13337, nil)
		client := NewTunnel("127.0.0.1", 13337)
		hop := NewTunnel("127.0.0.1", 13337)

		socket_client, socket_forward := test.NewLoopbackSocket()
		socket_hop, socket_gateway := test.NewLoopbackSocket()

		client.SetSocket(socket_client)
		forward.SetSocket(socket_forward)
		hop.SetSocket(socket_hop)
		gateway.SetSocket(socket_gateway)

		forward.SetTunnel(&hop)

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		go func() {
			forward.Listen()
		}()

		go func() {
			gateway.Listen()
		}()

		go func() {
			time.Sleep(2 * time.Second)
			forward.Destroy()
			gateway.Destroy()
		}()

		time.Sleep(100 * time.Millisecond)
		response := client.RequestPacket(request)
		response.Decode()

		if expected.Status != response.Status {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(expected.Status).String(), http.Status(response.Status).String())
		}

		if string(expected.Payload) != string(response.Payload) {
			t.Errorf("Expected HTTP response payload '%s' but got '%s'", string(expected.Payload), string(response.Payload))
		}

		time.Sleep(2 * time.Second)

	})

}
