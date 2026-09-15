package interfaces

import "tholian-warps/protocols/dns"

type ResolverCache interface {
	Exists(dns.Packet) bool
	Read(dns.Packet)   dns.Packet
	Write(dns.Packet)  bool
}
