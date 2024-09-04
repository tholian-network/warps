package structs

import "tholian-endpoint/protocols/http"
import utils_url "tholian-warps/utils/net/url"
import "os"
import net_url "net/url"
import "sort"
import "strings"

func resolveWebCacheFile(url *net_url.URL) string {

	result := ""
	tmp := strings.Split(url.Path, "/")

	if len(tmp) > 1 {

		if tmp[0] == "" && strings.Contains(tmp[len(tmp)-1], ".") {
			result = url.Host + "/" + url.Path[1:]
		} else if tmp[0] == "" && strings.TrimSpace(tmp[len(tmp)-1]) != "" {
			result = url.Host + "/" + url.Path[1:]
		} else {
			result = url.Host + "/" + url.Path[1:len(url.Path)-1] + "/index.html"
		}

	} else {

		result = url.Host + "/index.html"

	}

	query := url.Query()

	if len(query) > 0 {

		parameters := []string{}

		for key := range query {

			val := query.Get(key)

			if !utils_url.IsXSSParameter(key, val) && !utils_url.IsTrackingParameter(url.Host, key, val) {
				parameters = append(parameters, key)
			}

		}

		sort.Strings(parameters)

		for p := 0; p < len(parameters); p++ {

			key := parameters[p]
			val := query.Get(key)

			if p == 0 {
				result += "?" + key + "=" + val
			} else {
				result += "&" + key + "=" + val
			}

		}

	}

	return result

}

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

		resolved := resolveWebCacheFile(request.URL)

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
		// TODO: Transfer-Encoding
		// TODO: Content-Encoding
		// TODO: Store as plaintext!
	}

	return response

}

func (cache *WebCache) Write(response http.Packet) bool {

	var result bool = false

	if response.Type == "response" && response.URL != nil {
		// TODO: Transfer-Encoding
		// TODO: Content-Encoding
		// TODO: Store as plaintext!
	}

	return result

}

