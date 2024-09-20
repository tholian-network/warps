package tunnel

import "tholian-endpoint/protocols/http"

func IsResolveRequest(request *http.Packet) bool {

	var result bool = false

	if request.Type == "request" {

		content_type := request.GetHeader("Accept")

		if content_type == "application/dns-message" {
			result = true
		}

	}

	return result

}
