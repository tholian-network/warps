package interfaces

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"

type Proxy interface {
	ResolvePacket(dns.Packet)  dns.Packet
	RequestPacket(http.Packet) http.Packet
	Destroy()                  error
	Listen()                   error
}
