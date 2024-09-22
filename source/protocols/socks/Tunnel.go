package socks

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
// import socks_tunnel "tholian-warps/protocols/socks/tunnel"

type Tunnel struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func NewTunnel(host string, port uint16) Tunnel {

	var tunnel Tunnel

	tunnel.Host = host
	tunnel.Port = port

	return tunnel

}

func (tunnel *Tunnel) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	// TODO: Implement SOCKS Protocol

	return response

}

func (tunnel *Tunnel) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	// TODO: Implement SOCKS Protocol

	return response

}
