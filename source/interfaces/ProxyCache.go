package interfaces

import "tholian-endpoint/protocols/http"

type ProxyCache interface {
	Exists(http.Packet) bool
	Read(http.Packet)   http.Packet
	Write(http.Packet)  bool
}
