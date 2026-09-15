package tunnel

import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"

func EncodeResponse(message *icmp.Message, response http.Packet) {

	message.SetType(icmp.MessageData)

	if response.URL != nil {
		message.SetTarget(response.URL.String())
	}

	message.SetPayload(response.Bytes())

}
