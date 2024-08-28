package structs

import warps_url "tholian-warps/net/url"
import "os"
import net_url "net/url"
import "sort"
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

func (cache *WebCache) Exists(raw_url string) bool {

	var result bool = false

	url, err0 := net_url.Parse(raw_url)

	if err0 == nil {

		resolved := cache.Resolve(url)

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

func (cache *WebCache) Resolve(url *net_url.URL) string {

	resolved := ""

	tmp := strings.Split(url.Path, "/")

	if len(tmp) > 1 {

		if tmp[0] == "" && strings.Contains(tmp[len(tmp)-1], ".") {
			resolved = url.Host + "/" + url.Path[1:]
		} else if tmp[0] == "" && tmp[len(tmp)-1] != "" {
			resolved = url.Host + "/" + url.Path[1:]
		} else {
			resolved = url.Host + "/" + url.Path[1:] + "index.html"
		}

	} else {

		resolved = url.Host + "/index.html"

	}

	query := url.Query()

	if len(query) > 0 {

		parameters := []string{}

		for key := range query {

			val := query.Get(key)

			if !warps_url.IsXSSParameter(key, val) && !warps_url.IsTrackingParameter(url.Host, key, val) {
				parameters = append(parameters, key)
			}

		}

		sort.Strings(parameters)

		for p := 0; p < len(parameters); p++ {

			key := parameters[p]
			val := query.Get(key)

			if p == 0 {
				resolved += "?" + key + "=" + val
			} else {
				resolved += "&" + key + "=" + val
			}

		}

	}

	return resolved

}
