package tunnel

import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"

func EncodeRequest(message *icmp.Message, request http.Packet) {

	message.SetType(icmp.MessageData)

	if request.URL != nil {
		message.SetTarget(request.URL.String())
	}

	message.SetPayload(request.Bytes())

}
