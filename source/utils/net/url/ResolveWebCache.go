package url

import net_url "net/url"
import "sort"
import "strings"

func ResolveWebCache(url *net_url.URL) string {

	path := url.Path
	result := ""

	if path == "/" {
		result = url.Host + "/index.html"
	} else if strings.HasPrefix(path, "/") && strings.HasSuffix(path, "/") {
		result = url.Host + path[0:len(path)-1] + "/index.html"
	} else if strings.HasPrefix(path, "/") {
		result = url.Host + path
	}

	query := url.Query()

	if len(query) > 0 {

		parameters := []string{}

		for key := range query {

			val := query.Get(key)

			if !IsXSSParameter(key, val) && !IsTrackingParameter(url.Host, key, val) {
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

