package structs

import "tholian-endpoint/protocols/http"
import net_url "net/url"
import "os"
import "strconv"
import "testing"

func TestWebCache(t *testing.T) {

	t.Run("Exists", func(t *testing.T) {

		url, _ := net_url.Parse("http://localhost:8080/folder/to/file.html")
		tmp, _ := os.MkdirTemp("", "tholian-warps-webcache-*")
		webcache := NewWebCache(tmp)

		request := http.NewPacket()
		request.SetURL(*url)
		request.SetMethod(http.MethodGet)

		payload := []byte("Hello, World!")
		response := http.NewPacket()
		response.SetURL(*url)
        response.SetStatus(http.StatusOK)
		response.SetHeader("Content-Type", "text/html")
		response.SetHeader("Content-Length", strconv.Itoa(len(payload)))
		response.SetPayload(payload)

		result1 := webcache.Exists(request)

		if result1 != false {
			t.Errorf("Expected '%t' but got '%t'", false, result1)
		}

		result2 := webcache.Write(response)

		if result2 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result2)
		}

		result3 := webcache.Exists(request)

		if result3 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result3)
		}

		os.RemoveAll(tmp)

	})

	t.Run("Read", func(t *testing.T) {

		url, _ := net_url.Parse("http://localhost:8080/folder/to/file.html")
		tmp, _ := os.MkdirTemp("", "tholian-warps-webcache-*")
		webcache := NewWebCache(tmp)

		request := http.NewPacket()
		request.SetURL(*url)
		request.SetMethod(http.MethodGet)

		payload := []byte("Hello, World!")
		response := http.NewPacket()
		response.SetURL(*url)
        response.SetStatus(http.StatusOK)
		response.SetHeader("Content-Type", "text/html")
		response.SetHeader("Content-Length", strconv.Itoa(len(payload)))
		response.SetPayload(payload)

		result1 := webcache.Read(request)

		if result1.Status != http.StatusNotFound {
			t.Errorf("Expected '%s' but got '%s'", http.Status(http.StatusNotFound).String(), result1.Status.String())
		}

		result2 := webcache.Write(response)

		if result2 != true {
			t.Errorf("Expected '%t' but got '%t'", true, result2)
		}

		result3 := webcache.Read(request)

		if result3.Status != http.StatusOK {
			t.Errorf("Expected '%s' but got '%s'", http.Status(http.StatusOK).String(), result1.Status.String())
		}

		os.RemoveAll(tmp)

	})

}

