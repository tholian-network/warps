package tunnel

import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"

func IsResponse(message *icmp.Message) bool {

	var result bool = false

	if message.Type == icmp.MessageData && message.Target != "dns-query" {

		response := http.Parse(message.Payload)

		if response.Type == "response" {
			result = true
		}

	}

	return result

}
