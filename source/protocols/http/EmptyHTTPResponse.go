package http

import "tholian-endpoint/protocols/http"
import "strings"

func EmptyHTTPResponse(request http.Packet) http.Packet {

	var response http.Packet

	// response.SetType("response")
	response.SetURL(*request.URL)

	content_type := request.GetHeader("Accept")

	if strings.Contains(content_type, ";") {
		// TODO: Support multiple accept headers
		response.SetHeader("Content-Type", "text/plain")
	} else if content_type != "" {
		response.SetHeader("Content-Type", content_type)
	}

	response.SetStatus(http.StatusNotFound)
	response.SetPayload([]byte{})

	return response

}
