package tunnel

import "tholian-warps/protocols/icmp"

func IsResolveResponse(message *icmp.Message) bool {

	var result bool = false

	if message.Type == icmp.MessageData && message.Target == "dns-query" {
		result = true
	}

	return result

}
