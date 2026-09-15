package tunnel

import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"

func IsRequest(message *icmp.Message) bool {

	var result bool = false

	if message.Type == icmp.MessageData && message.Target != "dns-query" {

		request := http.Parse(message.Payload)

		if request.Type == "request" {
			result = true
		}

	}

	return result

}
