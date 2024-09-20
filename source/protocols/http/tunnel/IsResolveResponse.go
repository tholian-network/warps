package tunnel

import "tholian-endpoint/protocols/http"

func IsResolveResponse(response *http.Packet) bool {

	var result bool = false

	if response.Type == "response" {

		content_type := response.GetHeader("Content-Type")

		if content_type == "application/dns-message" {
			result = true
		}

	}

	return result

}
