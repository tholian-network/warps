package interfaces

import "tholian-endpoint/protocols/dns"

type ResolverCache interface {
	Exists(dns.Packet) bool
	Read(dns.Packet)   dns.Packet
	Write(dns.Packet)  bool
}
