package http

import "bytes"
import "net/http"
import "net/http/httptest"
import "strconv"
import "testing"
import "fmt"

func TestRequest(t *testing.T) {

	t.Run("Request(http://ipv4:port)", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.Header().Set("Content-Type", "text/plain")
			response.Header().Set("X-Status", "OK")
			response.WriteHeader(http.StatusOK)
			response.Write([]byte(http.StatusText(http.StatusOK)))
		}))

		response := Request(server.URL + "/path/to/index.html")

		if response.Type != "response" {
			t.Errorf("Expected %s to be %s", response.Type, "response")
		}

		if response.URL != nil {
			t.Errorf("Expected %s to be nil", response.URL.String())
		}

		if response.Status != StatusOK {
			t.Errorf("Expected %s to be %s", response.Status.String(), "OK")
		}

		if len(response.Headers) == 7 {

			if response.Headers["Accept-Encoding"] != "bzip2, gzip, deflate, zstd" {
				t.Errorf("Expected %s to be %s", response.Headers["Accept-Encoding"], "bzip2, gzip, deflate, zstd")
			}

			if response.Headers["Connection"] != "close" {
				t.Errorf("Expected %s to be %s", response.Headers["Connection"], "close")
			}

			if response.Headers["Content-Encoding"] != "identity" {
				t.Errorf("Expected %s to be %s", response.Headers["Content-Encoding"], "identity")
			}

			if response.Headers["Content-Length"] != strconv.Itoa(len("OK")) {
				t.Errorf("Expected %s to be %s", response.Headers["Content-Length"], strconv.Itoa(len("OK")))
			}

			if response.Headers["Content-Type"] != "text/plain" {
				t.Errorf("Expected %s to be %s", response.Headers["Content-Type"], "text/plain")
			}

			if response.Headers["X-Status"] != "OK" {
				t.Errorf("Expected %s to be %s", response.Headers["X-Status"], "OK")
			}

		} else {
			t.Errorf("Expected %d headers to be %d", len(response.Headers), 7)
		}

		if len(response.Payload) == 2 {

			if bytes.Compare(response.Payload, []byte("OK")) != 0 {
				t.Errorf("Expected payload %v to be %v", response.Payload, []byte("OK"))
			}

		} else {
			t.Errorf("Expected %d payload bytes to be %d", len(response.Payload), 2)
		}

		if response.Server != nil {
			t.Errorf("Expected %v to be nil", response.Server)
		}

		if response.Certificate != nil {
			t.Errorf("Expected %v to be nil", response.Certificate)
		}

		if response.Encoding != EncodingIdentity {
			t.Errorf("Expected %s to be %s", response.Encoding.String(), "identity")
		}

		server.Close()

	})

	t.Run("Request(https://ipv4:port)", func(t *testing.T) {

		server := httptest.NewTLSServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.Header().Set("Content-Type", "text/plain")
			response.Header().Set("X-Status", "OK")
			response.WriteHeader(http.StatusOK)
			response.Write([]byte(http.StatusText(http.StatusOK)))
		}))

		response := Request(server.URL + "/path/to/index.html")

		fmt.Println(response)

		server.Close()

		// TODO
	})

	t.Run("Request(http://domain:port)", func(t *testing.T) {
		// TODO
	})

	t.Run("Request(https://domain:port)", func(t *testing.T) {
		// TODO
	})


}
