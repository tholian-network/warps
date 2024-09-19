package interfaces

import "tholian-endpoint/protocols/dns"

type Resolver interface {
	Resolve(string)           dns.Packet
	ResolvePacket(dns.Packet) dns.Packet
}
