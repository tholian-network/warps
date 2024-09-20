package http

func IsFilteredHeader(header string) bool {

	var result bool = false

	headers := []string{

		// Request Headers
		"Accept",
		"Accept-Encoding",
		"Accept-Ranges",
		"Cookie",
		"If-Match",
		"If-Modified-Since",
		"If-Range",
		"If-Unmodified-Since",
		"Upgrade",

		// Response Headers
		"Cache-Control",
		"Content-Range",
		"Link",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Range",
		"Refresh",
		"Set-Cookie",
		"Strict-Transport-Security",
		"Upgrade-Insecure-Requests",
		"Via",

		// Request + Response Headers
		"Connection",
		"Content-Encoding",
		"Content-Length",
		"Content-Type",
		"Keep-Alive",
		"Transfer-Encoding",

	}

	for h := 0; h < len(headers); h++ {

		if headers[h] == header {
			result = true
			break
		}

	}

	return result

}
