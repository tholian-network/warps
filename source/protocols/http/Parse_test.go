package http

import "bytes"
import "compress/flate"
import "compress/gzip"
import "io"
import "strconv"
import "strings"
import "testing"

func TestParse(t *testing.T) {

	t.Run("Parse(GET)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"GET /path/to/index.html HTTP/1.1",
			"Host: example.com",
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
		}, "\r\n") + "\r\n\r\n")

		request := Parse(bytes)

		if request.Type != "request" {
			t.Errorf("Expected %s to be %s", request.Type, "request")
		}

		if request.URL.String() != "/path/to/index.html" {
			t.Errorf("Expected %s to be %s", request.URL.String(), "/path/to/index.html")
		}

		if request.Method != MethodGet {
			t.Errorf("Expected %s to be %s", request.Method.String(), "GET")
		}

		if request.Status.String() != "" {
			t.Errorf("Expected %s to be empty", request.Status.String())
		}

		if len(request.Headers) == 5 {

			if request.Headers["Accept"] != "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept"], "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
			}

			if request.Headers["Accept-Encoding"] != "bzip2, gzip, deflate, zstd" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept-Encoding"], "bzip2, gzip, deflate, zstd")
			}

			if request.Headers["Content-Encoding"] != "identity" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Encoding"], "identity")
			}

			if request.Headers["Host"] != "example.com" {
				t.Errorf("Expected %s to be %s", request.Headers["Host"], "example.com")
			}

			if request.Headers["User-Agent"] != "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0" {
				t.Errorf("Expected %s to be %s", request.Headers["User-Agent"], "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")
			}

		} else {
			t.Errorf("Expected %d headers to be %d", len(request.Headers), 5)
		}

		if len(request.Payload) != 0 {
			t.Errorf("Expected %d payload bytes to be %d", len(request.Payload), 0)
		}

		if request.Server != nil {
			t.Errorf("Expected %v to be nil", request.Server)
		}

		if request.Certificate != nil {
			t.Errorf("Expected %v to be nil", request.Certificate)
		}

		if request.Encoding != EncodingIdentity {
			t.Errorf("Expected %s to be %s", request.Encoding.String(), "identity")
		}

	})

	t.Run("Parse(GET with parameters)", func(t *testing.T) {

		bytes := []byte(strings.Join([]string{
			"GET /path/to/login.php?username=jean.luc.picard&password=picard#121&email=jl.picard@federation.int HTTP/1.1",
			"Host: example.com",
			"Accept: application/json",
		}, "\r\n") + "\r\n\r\n")

		request := Parse(bytes)

		if request.Type != "request" {
			t.Errorf("Expected %s to be %s", request.Type, "request")
		}

		if request.URL.String() != "/path/to/login.php?username=jean.luc.picard&password=picard#121&email=jl.picard@federation.int" {
			t.Errorf("Expected %s to be %s", request.URL.String(), "/path/to/login.php?username=jean.luc.picard&password=picard#121&email=jl.picard@federation.int")
		}

		if request.Method != MethodGet {
			t.Errorf("Expected %s to be %s", request.Method.String(), "GET")
		}

		if request.Status.String() != "" {
			t.Errorf("Expected %s to be empty", request.Status.String())
		}

		if len(request.Headers) == 4 {

			if request.Headers["Accept"] != "application/json" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept"], "application/json" )
			}

			if request.Headers["Accept-Encoding"] != "bzip2, gzip, deflate, zstd" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept-Encoding"], "bzip2, gzip, deflate, zstd")
			}

			if request.Headers["Content-Encoding"] != "identity" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Encoding"], "identity")
			}

			if request.Headers["Host"] != "example.com" {
				t.Errorf("Expected %s to be %s", request.Headers["Host"], "example.com")
			}

		} else {
			t.Errorf("Expected %d headers to be %d", len(request.Headers), 4)
		}

		if len(request.Payload) != 0 {
			t.Errorf("Expected %d payload bytes to be %d", len(request.Payload), 0)
		}

		if request.Server != nil {
			t.Errorf("Expected %v to be nil", request.Server)
		}

		if request.Certificate != nil {
			t.Errorf("Expected %v to be nil", request.Certificate)
		}

		if request.Encoding != EncodingIdentity {
			t.Errorf("Expected %s to be %s", request.Encoding.String(), "identity")
		}

	})

	t.Run("Parse(POST json)", func(t *testing.T) {

		raw_payload := strings.Join([]string{
			"{",
			"\t\"name\": \"Jean-Luc Picard\",",
			"\t\"email\": \"jl.picard@federation.int\",",
			"\t\"password\": \"picard#121\"",
			"}",
		}, "\n")

		raw_bytes := []byte(strings.Join([]string{
			"POST /path/to/upload.php HTTP/1.1",
			"Host: upload.example.com",
			"Accept: application/json",
			"Content-Type: application/json",
			"Content-Length: " + strconv.Itoa(len(raw_payload)),
			"",
			raw_payload,
		}, "\r\n"))

		request := Parse(raw_bytes)

		if request.Type != "request" {
			t.Errorf("Expected %s to be %s", request.Type, "request")
		}

		if request.URL.String() != "/path/to/upload.php" {
			t.Errorf("Expected %s to be %s", request.URL.String(), "/path/to/upload.php")
		}

		if request.Method != MethodPost {
			t.Errorf("Expected %s to be %s", request.Method.String(), "POST")
		}

		if request.Status.String() != "" {
			t.Errorf("Expected %s to be empty", request.Status.String())
		}

		if len(request.Headers) == 6 {

			if request.Headers["Accept"] != "application/json" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept"], "application/json" )
			}

			if request.Headers["Accept-Encoding"] != "bzip2, gzip, deflate, zstd" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept-Encoding"], "bzip2, gzip, deflate, zstd")
			}

			if request.Headers["Content-Encoding"] != "identity" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Encoding"], "identity")
			}

			if request.Headers["Host"] != "upload.example.com" {
				t.Errorf("Expected %s to be %s", request.Headers["Host"], "upload.example.com")
			}

			if request.Headers["Content-Length"] != strconv.Itoa(len(raw_payload)) {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Length"], strconv.Itoa(len(raw_payload)))
			}

			if request.Headers["Content-Type"] != "application/json" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Type"], "application/json")
			}

			if request.Headers["Host"] != "upload.example.com" {
				t.Errorf("Expected %s to be %s", request.Headers["Host"], "upload.example.com")
			}

		} else {
			t.Errorf("Expected %d headers to be %d", len(request.Headers), 6)
		}

		if len(request.Payload) == len(raw_payload) {

			if bytes.Compare(request.Payload, []byte(raw_payload)) != 0 {
				t.Errorf("Expected payload %v to be %v", request.Payload, []byte(raw_payload))
			}

		} else {
			t.Errorf("Expected %d payload bytes to be %d", len(request.Payload), len(raw_payload))
		}

		if request.Server != nil {
			t.Errorf("Expected %v to be nil", request.Server)
		}

		if request.Certificate != nil {
			t.Errorf("Expected %v to be nil", request.Certificate)
		}

		if request.Encoding != EncodingIdentity {
			t.Errorf("Expected %s to be %s", request.Encoding.String(), "identity")
		}

	})

	t.Run("Parse(POST x-www-form-urlencoded)", func(t *testing.T) {

		raw_payload := []byte("username=jean.luc.picard&password=picard#121&email=jl.picard@federation.int")
		raw_bytes   := []byte(strings.Join([]string{
			"POST /path/to/login.php HTTP/1.1",
			"Host: example.com",
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Content-Type: application/x-www-form-urlencoded",
			"Content-Length: 75",
			"",
		}, "\r\n"))
		raw_bytes = append(raw_bytes, []byte("\r\n")...)
		raw_bytes = append(raw_bytes, raw_payload...)

		request := Parse(raw_bytes)

		if request.Type != "request" {
			t.Errorf("Expected %s to be %s", request.Type, "request")
		}

		if request.URL.String() != "/path/to/login.php" {
			t.Errorf("Expected %s to be %s", request.URL.String(), "/path/to/login.php")
		}

		if request.Method != MethodPost {
			t.Errorf("Expected %s to be %s", request.Method.String(), "POST")
		}

		if request.Status.String() != "" {
			t.Errorf("Expected %s to be empty", request.Status.String())
		}

		if len(request.Headers) == 6 {

			if request.Headers["Accept"] != "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept"], "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
			}

			if request.Headers["Accept-Encoding"] != "bzip2, gzip, deflate, zstd" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept-Encoding"], "bzip2, gzip, deflate, zstd")
			}

			if request.Headers["Content-Encoding"] != "identity" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Encoding"], "identity")
			}

			if request.Headers["Content-Length"] != "75" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Length"], "75")
			}

			if request.Headers["Content-Type"] != "application/x-www-form-urlencoded" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Type"], "application/x-www-form-urlencoded")
			}

			if request.Headers["Host"] != "example.com" {
				t.Errorf("Expected %s to be %s", request.Headers["Host"], "example.com")
			}

		} else {
			t.Errorf("Expected %d headers to be %d", len(request.Headers), 6)
		}

		if len(request.Payload) == len(raw_payload) {

			if bytes.Compare(request.Payload, raw_payload) != 0 {
				t.Errorf("Expected payload %v to be %v", request.Payload, []byte(raw_payload))
			}

		} else {
			t.Errorf("Expected payload %v to be %v", request.Payload, raw_payload)
		}

		if request.Server != nil {
			t.Errorf("Expected %v to be nil", request.Server)
		}

		if request.Certificate != nil {
			t.Errorf("Expected %v to be nil", request.Certificate)
		}

		if request.Encoding != EncodingIdentity {
			t.Errorf("Expected %s to be %s", request.Encoding.String(), "identity")
		}

	})

	// t.Run("Parse(POST bzip2 encoded)", func(t *testing.T) {

		// TODO: gzip encoded Payload test
		// compress/bzip2 doesn't support io.Writer yet

	// })

	t.Run("Parse(POST deflate encoded)", func(t *testing.T) {

		var encoded_payload bytes.Buffer

		raw_payload  := []byte("{\"buffer\": \"This is an example payload!\"}")
		reader       := bytes.NewBuffer(raw_payload)
		writer, err1 := flate.NewWriter(&encoded_payload, flate.BestSpeed)

		if err1 == nil {

			_, err2 := io.Copy(writer, reader)

			writer.Flush()
			writer.Close()

			if err2 != nil {
				t.Errorf("Expected %s to be nil", err2.Error())
			}

		} else {
			t.Errorf("Expected %s to be nil", err1.Error())
		}

		raw_bytes := []byte(strings.Join([]string{
			"POST /path/to/upload.php HTTP/1.1",
			"Host: upload.example.com",
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Content-Encoding: deflate",
			"Content-Type: application/json",
			"Content-Length: " + strconv.Itoa(len(encoded_payload.Bytes())),
			"",
		}, "\r\n"))
		raw_bytes = append(raw_bytes, []byte("\r\n")...)
		raw_bytes = append(raw_bytes, encoded_payload.Bytes()...)

		request := Parse(raw_bytes)

		if request.Type != "request" {
			t.Errorf("Expected %s to be %s", request.Type, "request")
		}

		if request.URL.String() != "/path/to/upload.php" {
			t.Errorf("Expected %s to be %s", request.URL.String(), "/path/to/upload.php")
		}

		if request.Method != MethodPost {
			t.Errorf("Expected %s to be %s", request.Method.String(), "POST")
		}

		if request.Status.String() != "" {
			t.Errorf("Expected %s to be empty", request.Status.String())
		}

		if len(request.Headers) == 6 {

			if request.Headers["Accept"] != "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept"], "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
			}

			if request.Headers["Accept-Encoding"] != "bzip2, gzip, deflate, zstd" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept-Encoding"], "bzip2, gzip, deflate, zstd")
			}

			if request.Headers["Content-Encoding"] != "deflate" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Encoding"], "deflate")
			}

			if request.Headers["Content-Length"] != strconv.Itoa(len(encoded_payload.Bytes())) {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Length"], strconv.Itoa(len(encoded_payload.Bytes())))
			}

			if request.Headers["Content-Type"] != "application/json" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Type"], "application/json")
			}

			if request.Headers["Host"] != "upload.example.com" {
				t.Errorf("Expected %s to be %s", request.Headers["Host"], "upload.example.com")
			}

		} else {
			t.Errorf("Expected %d headers to be %d", len(request.Headers), 6)
		}

		if len(request.Payload) == len(encoded_payload.Bytes()) {

			if bytes.Compare(request.Payload, []byte(encoded_payload.Bytes())) != 0 {
				t.Errorf("Expected payload %v to be %v", request.Payload, []byte(encoded_payload.Bytes()))
			}

		} else {
			t.Errorf("Expected %d payload bytes to be %d", len(request.Payload), len(encoded_payload.Bytes()))
		}

		if request.Server != nil {
			t.Errorf("Expected %v to be nil", request.Server)
		}

		if request.Certificate != nil {
			t.Errorf("Expected %v to be nil", request.Certificate)
		}

		if request.Encoding != EncodingDeflate {
			t.Errorf("Expected %s to be %s", request.Encoding.String(), "deflate")
		}

	})

	t.Run("Parse(POST gzip encoded)", func(t *testing.T) {

		var encoded_payload bytes.Buffer

		raw_payload := []byte("This is an example payload!")
		reader      := bytes.NewBuffer(raw_payload)
		writer      := gzip.NewWriter(&encoded_payload)

		_, err1 := io.Copy(writer, reader)

		writer.Flush()
		writer.Close()

		if err1 != nil {
			t.Errorf("Expected %s to be nil", err1.Error())
		}

		raw_bytes := []byte(strings.Join([]string{
			"POST /path/to/upload.php HTTP/1.1",
			"Host: upload.example.com",
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Content-Encoding: gzip",
			"Content-Type: application/json",
			"Content-Length: " + strconv.Itoa(len(encoded_payload.Bytes())),
			"",
		}, "\r\n"))
		raw_bytes = append(raw_bytes, []byte("\r\n")...)
		raw_bytes = append(raw_bytes, encoded_payload.Bytes()...)

		request := Parse(raw_bytes)

		if request.Type != "request" {
			t.Errorf("Expected %s to be %s", request.Type, "request")
		}

		if request.URL.String() != "/path/to/upload.php" {
			t.Errorf("Expected %s to be %s", request.URL.String(), "/path/to/upload.php")
		}

		if request.Method != MethodPost {
			t.Errorf("Expected %s to be %s", request.Method.String(), "POST")
		}

		if request.Status.String() != "" {
			t.Errorf("Expected %s to be empty", request.Status.String())
		}


		if len(request.Headers) == 6 {

			if request.Headers["Accept"] != "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept"], "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
			}

			if request.Headers["Accept-Encoding"] != "bzip2, gzip, deflate, zstd" {
				t.Errorf("Expected %s to be %s", request.Headers["Accept-Encoding"], "bzip2, gzip, deflate, zstd")
			}

			if request.Headers["Content-Encoding"] != "gzip" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Encoding"], "gzip")
			}

			if request.Headers["Content-Length"] != strconv.Itoa(len(encoded_payload.Bytes())) {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Length"], strconv.Itoa(len(encoded_payload.Bytes())))
			}

			if request.Headers["Content-Type"] != "application/json" {
				t.Errorf("Expected %s to be %s", request.Headers["Content-Type"], "application/json")
			}

			if request.Headers["Host"] != "upload.example.com" {
				t.Errorf("Expected %s to be %s", request.Headers["Host"], "upload.example.com")
			}

		} else {
			t.Errorf("Expected %d headers to be %d", len(request.Headers), 6)
		}

		if len(request.Payload) == len(encoded_payload.Bytes()) {

			if bytes.Compare(request.Payload, []byte(encoded_payload.Bytes())) != 0 {
				t.Errorf("Expected payload %v to be %v", request.Payload, []byte(encoded_payload.Bytes()))
			}

		} else {
			t.Errorf("Expected %d payload bytes to be %d", len(request.Payload), len(encoded_payload.Bytes()))
		}

		if request.Server != nil {
			t.Errorf("Expected %v to be nil", request.Server)
		}

		if request.Certificate != nil {
			t.Errorf("Expected %v to be nil", request.Certificate)
		}

		if request.Encoding != EncodingGzip {
			t.Errorf("Expected %s to be %s", request.Encoding.String(), "gzip")
		}

	})

	// t.Run("Parse(POST zstd encoded)", func(t *testing.T) {

		// TODO: zstd encoded Payload test
		// internal/zstd doesn't support io.Writer yet

	// })

}
