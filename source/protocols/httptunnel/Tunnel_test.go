package httptunnel

import "tholian-warps/protocols/http"
import net_http "net/http"
import "net/http/httptest"
import net_url "net/url"
import "strconv"
import "testing"

func TestTunnel(t *testing.T) {

	t.Run("Tunnel with HTTP Payload", func(t *testing.T) {

		server := httptest.NewServer(net_http.HandlerFunc(func(w net_http.ResponseWriter, r *net_http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Hello, world!"))
		}))

		defer server.Close()

		url, _ := net_url.Parse(server.URL)
		port, _ := strconv.ParseUint(url.Port(), 10, 16)

		tunnel := NewTunnel(url.Hostname(), uint16(port))

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)
		request.SetHeader("Host", url.Host)

		response := tunnel.RequestPacket(request)
		response.Decode()

		if response.Status != http.StatusOK {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(http.StatusOK).String(), http.Status(response.Status).String())
		}

		if string(response.Payload) != "Hello, world!" {
			t.Errorf("Expected HTTP response payload '%s' but got '%s'", "Hello, world!", string(response.Payload))
		}

	})

}
