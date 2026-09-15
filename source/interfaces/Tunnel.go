package interfaces

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"

type Tunnel interface {
	ResolvePacket(dns.Packet)  dns.Packet
	RequestPacket(http.Packet) http.Packet
}
