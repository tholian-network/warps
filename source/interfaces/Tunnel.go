package interfaces

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"

type Tunnel interface {
	ResolvePacket(dns.Packet)  dns.Packet
	RequestPacket(http.Packet) http.Packet
}