package structs

import "tholian-endpoint/protocols/http"
import utils_url "tholian-warps/utils/net/url"
import utils_http "tholian-warps/utils/protocols/http"
import "encoding/json"
import "os"
import "strings"

type WebCache struct {
	Folder string `json:"folder"`
}

func NewWebCache(folder string) WebCache {

	var cache WebCache

	if strings.HasSuffix(folder, "/") {
		folder = folder[0:len(folder)-1]
	}

	stat, err1 := os.Stat(folder)

	if err1 == nil && stat.IsDir() {

		cache.Folder = folder

	} else {

		err2 := os.MkdirAll(folder, 0750)

		if err2 == nil {
			cache.Folder = folder
		}

	}

	return cache

}

func (cache *WebCache) Exists(request http.Packet) bool {

	var result bool = false

	if request.Type == "request" && request.URL != nil {

		resolved := utils_url.ResolveWebCache(request.URL)

		if resolved != "" {

			stat1, err1 := os.Stat(cache.Folder + "/headers/" + resolved)
			stat2, err2 := os.Stat(cache.Folder + "/payload/" + resolved)

			if err1 == nil && err2 == nil {

				if !stat1.IsDir() && !stat2.IsDir() {
					result = true
				}

			}

		}

	}

	return result

}

func (cache *WebCache) Read(request http.Packet) http.Packet {

	var response http.Packet

	if request.Type == "request" && request.URL != nil {

		resolved := utils_url.ResolveWebCache(request.URL)

		if resolved != "" {

			response = http.NewPacket()
			response.SetURL(*request.URL)
			response.SetStatus(http.StatusOK)

			buffer1, err1 := os.ReadFile(cache.Folder + "/headers/" + resolved)

			if err1 == nil {

				headers := make(map[string]string)
				err12 := json.Unmarshal(buffer1, &headers)

				if err12 == nil {

					for key, val := range headers {
						response.SetHeader(key, val)
					}

				}

			}

			buffer2, err2 := os.ReadFile(cache.Folder + "/payload/" + resolved)

			if err2 == nil {
				response.SetPayload(buffer2)
			}

		}

	}

	return response

}

func (cache *WebCache) Write(response http.Packet) bool {

	var result bool = false

	if response.Type == "response" && response.URL != nil {

		resolved := utils_url.ResolveWebCache(response.URL)

		if resolved != "" {

			response.Decode()

			headers := make(map[string]string)
			payload := response.Payload

			for key, val := range response.Headers {

				if !utils_http.IsFilteredHeader(key) {
					headers[key] = val
				}

			}

			buffer_headers, err0 := json.MarshalIndent(headers, "", "\t")

			if err0 == nil {

				err1 := os.WriteFile(cache.Folder + "/headers/" + resolved, buffer_headers, 0666)
				err2 := os.WriteFile(cache.Folder + "/payload/" + resolved, payload, 0666)

				if err1 == nil && err2 == nil {
					result = true
				}

			}

		}

	}

	return result

}

