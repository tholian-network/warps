package dns

import "tholian-endpoint/protocols/http"
import "tholian-warps/protocols/test"
import net_url "net/url"
import "strconv"
import "testing"
import "time"

func TestProxy(t *testing.T) {

	t.Run("Proxy with ProxyCache and Small Payload", func(t *testing.T) {

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

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		go func() {

			err := proxy.Listen()

			if err != nil {
				t.Errorf("Unexpected error '%s'", err.Error())
			}

		}()

		go func() {

			time.Sleep(1 * time.Second)

			err := proxy.Destroy()

			if err != nil {
				t.Errorf("Unexpected error '%s'", err.Error())
			}

		}()

		response := tunnel.RequestPacket(request)
		response.Decode()

		if expected.Status != response.Status {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(expected.Status).String(), http.Status(response.Status).String())
		}

		if string(expected.Payload) != string(response.Payload) {
			t.Errorf("Expected HTTP response payload '%s' but got '%s'", string(expected.Payload), string(response.Payload))
		}

		time.Sleep(1 * time.Second)

	})

	t.Run("Proxy with ProxyCache and Large Payload", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.txt")
		expected := http.NewPacket()
		expected.SetURL(*url)
		expected.SetStatus(http.StatusOK)
		expected.SetHeader("Content-Type", "text/plain")
		expected.SetHeader("X-Proxy", "SpyProxyCache")

		payload := make([]byte, 0)

		for l := 0; l < 100; l++ {
			payload = append(payload, []byte("Hello, line " + strconv.Itoa(l) + "!\n")...)
		}

		expected.SetPayload(payload)

		cache := test.NewSpyProxyCache(true, &expected, false)
		proxy := NewProxy("localhost", 13337, &cache)
		tunnel := NewTunnel("127.0.0.1", 13337)

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		go func() {

			err := proxy.Listen()

			if err != nil {
				t.Errorf("Unexpected error '%s'", err.Error())
			}

		}()

		go func() {

			time.Sleep(1 * time.Second)

			err := proxy.Destroy()

			if err != nil {
				t.Errorf("Unexpected error '%s'", err.Error())
			}

		}()

		response := tunnel.RequestPacket(request)
		response.Decode()

		if expected.Status != response.Status {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(expected.Status).String(), http.Status(response.Status).String())
		}

		if string(expected.Payload) != string(response.Payload) {
			t.Errorf("Expected HTTP different response payload")
		}

		time.Sleep(1 * time.Second)

	})

}
