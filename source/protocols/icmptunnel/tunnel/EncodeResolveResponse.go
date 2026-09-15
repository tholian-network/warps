package tunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/icmp"

func EncodeResolveResponse(message *icmp.Message, response dns.Packet) {

	message.SetType(icmp.MessageData)
	message.SetTarget("dns-query")
	message.SetPayload(response.Bytes())

}
