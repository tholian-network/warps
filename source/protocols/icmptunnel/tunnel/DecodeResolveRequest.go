package tunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/icmp"

func DecodeResolveRequest(message *icmp.Message) dns.Packet {
	return dns.Parse(message.Payload)
}
