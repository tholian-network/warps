package tunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/icmp"

func EncodeResolveRequest(message *icmp.Message, query dns.Packet) {

	message.SetType(icmp.MessageData)
	message.SetTarget("dns-query")
	message.SetPayload(query.Bytes())

}
