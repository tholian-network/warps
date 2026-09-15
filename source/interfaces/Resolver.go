package interfaces

import "tholian-warps/protocols/dns"

type Resolver interface {
	Resolve(string)           dns.Packet
	ResolvePacket(dns.Packet) dns.Packet
}
