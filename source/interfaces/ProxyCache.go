package interfaces

import "tholian-warps/protocols/http"

type ProxyCache interface {
	Exists(http.Packet) bool
	Read(http.Packet)   http.Packet
	Write(http.Packet)  bool
}
