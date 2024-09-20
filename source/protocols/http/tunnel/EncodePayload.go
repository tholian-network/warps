package tunnel

import "tholian-endpoint/protocols/http"
import "strconv"

func EncodePayload(packet *http.Packet, payload []byte) bool {

	var result bool = false

	if packet.Type == "request" {

		packet.SetMethod(http.MethodPost)
		packet.SetHeader("Accept", "application/dns-message")
		packet.SetHeader("Content-Type", "application/dns-message")
		packet.SetHeader("Content-Length", strconv.Itoa(len(payload)))
		packet.SetPayload(payload)

		result = true

	} else if packet.Type == "response" {

		packet.SetStatus(http.StatusOK)
		packet.SetHeader("Content-Type", "application/dns-message")
		packet.SetHeader("Content-Length", strconv.Itoa(len(payload)))
		packet.SetPayload(payload)

		result = true

	}

	return result

}
