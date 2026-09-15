package tunnel

import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"
import net_url "net/url"
import "testing"

func TestEncodeRequest(t *testing.T) {

	t.Run("HTTP Request round-trip", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.html")
		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)

		message := icmp.NewMessage()
		EncodeRequest(&message, request)

		if IsRequest(&message) != true {
			t.Errorf("Expected message to be a request")
		}

		if message.Target != url.String() {
			t.Errorf("Expected target '%s' but got '%s'", url.String(), message.Target)
		}

		decoded := DecodeRequest(&message)

		if decoded.Method != http.MethodGet {
			t.Errorf("Expected method '%s' but got '%s'", http.MethodGet.String(), decoded.Method.String())
		}

		if decoded.URL.String() != url.String() {
			t.Errorf("Expected url '%s' but got '%s'", url.String(), decoded.URL.String())
		}

	})

	t.Run("HTTP Response round-trip", func(t *testing.T) {

		url, _ := net_url.Parse("http://example.com/index.html")
		response := http.NewPacket()
		response.SetURL(*url)
		response.SetStatus(http.StatusOK)
		response.SetPayload([]byte("Hello, world!"))

		message := icmp.NewMessage()
		EncodeResponse(&message, response)

		if IsResponse(&message) != true {
			t.Errorf("Expected message to be a response")
		}

		decoded := DecodeResponse(&message)

		if decoded.Status != http.StatusOK {
			t.Errorf("Expected status '%s' but got '%s'", http.Status(http.StatusOK).String(), http.Status(decoded.Status).String())
		}

		if string(decoded.Payload) != "Hello, world!" {
			t.Errorf("Expected payload '%s' but got '%s'", "Hello, world!", string(decoded.Payload))
		}

	})

}
