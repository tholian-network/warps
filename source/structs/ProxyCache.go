package structs

import "tholian-endpoint/protocols/http"
import utils_url "tholian-warps/utils/net/url"
import utils_http "tholian-warps/utils/protocols/http"
import "encoding/json"
import "os"
import "path"
import "strings"

type ProxyCache struct {
	Folder string `json:"folder"`
}

func NewProxyCache(folder string) ProxyCache {

	var cache ProxyCache

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

func (cache *ProxyCache) Exists(request http.Packet) bool {

	var result bool = false

	if request.Type == "request" && request.URL != nil {

		resolved := utils_url.ResolveCache(request.URL)

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

func (cache *ProxyCache) Read(request http.Packet) http.Packet {

	var response http.Packet

	if request.Type == "request" && request.URL != nil {

		resolved := utils_url.ResolveCache(request.URL)

		if resolved != "" {

			buffer1, err1 := os.ReadFile(cache.Folder + "/headers/" + resolved)
			buffer2, err2 := os.ReadFile(cache.Folder + "/payload/" + resolved)

			if err1 == nil && err2 == nil {

				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusOK)

				headers := make(map[string]string)
				err12 := json.Unmarshal(buffer1, &headers)

				if err12 == nil {

					for key, val := range headers {
						response.SetHeader(key, val)
					}

				}

				response.SetPayload(buffer2)

			} else {

				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusNotFound)

			}

		} else {

			response = http.NewPacket()
			response.SetURL(*request.URL)
			response.SetStatus(http.StatusNotFound)

		}

	} else {

		response = http.NewPacket()
		response.SetURL(*request.URL)
		response.SetStatus(http.StatusInternalServerError)

	}

	return response

}

func (cache *ProxyCache) Write(response http.Packet) bool {

	var result bool = false

	if response.Type == "response" && response.URL != nil {

		resolved := utils_url.ResolveCache(response.URL)

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

				folder := path.Dir(resolved)
				_, err1 := os.Stat(cache.Folder + "/headers/" + folder)
				_, err2 := os.Stat(cache.Folder + "/payload/" + folder)

				if err1 != nil {

					err12 := os.MkdirAll(cache.Folder + "/headers/" + folder, 0755)

					if err12 == nil {
						_, err1 = os.Stat(cache.Folder + "/headers/" + folder)
					}

				}

				if err2 != nil {

					err22 := os.MkdirAll(cache.Folder + "/payload/" + folder, 0755)

					if err22 == nil {
						_, err2 = os.Stat(cache.Folder + "/payload/" + folder)
					}

				}

				if err1 == nil && err2 == nil {

					err12 := os.WriteFile(cache.Folder + "/headers/" + resolved, buffer_headers, 0666)
					err22 := os.WriteFile(cache.Folder + "/payload/" + resolved, payload, 0666)

					if err12 == nil && err22 == nil {
						result = true
					}

				}

			}

		}

	}

	return result

}

