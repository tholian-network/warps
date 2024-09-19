package dns

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/types"
import "tholian-warps/interfaces"
import "net"
import "strings"

type Resolver struct {
	Host   string                   `json:"host"`
	Port   uint16                   `json:"port"`
	Cache  interfaces.ResolverCache `json:"cache"`
	Tunnel interfaces.Tunnel        `json:"tunnel"`
}

func NewResolver(host string, port uint16, cache interfaces.ResolverCache) Resolver {

	var resolver Resolver

	if host == "*" || host == "any" || host == "localhost" {
		resolver.Host = "localhost"
	} else if strings.ToLower(host) == host {
		resolver.Host = host
	} else {
		resolver.Host = "localhost"
	}

	if port > 0 && port < 65535 {
		resolver.Port = port
	} else {
		resolver.Port = 5353
	}

	if cache != nil {
		resolver.Cache = cache
	}

	return resolver

}

func (resolver *Resolver) Resolve(domain string) dns.Packet {

	var response dns.Packet

	if types.IsDomain(domain) {

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion(domain, dns.TypeA))
		query.AddQuestion(dns.NewQuestion(domain, dns.TypeAAAA))

		response = resolver.ResolvePacket(query)

	}

	return response

}

func (resolver *Resolver) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	if resolver.Cache != nil && resolver.Cache.Exists(query) {
		response = resolver.Cache.Read(query)
	} else if resolver.Tunnel != nil {
		response = resolver.Tunnel.ResolvePacket(query)
	} else {
		response = dns.ResolvePacket(query)
	}

	return response

}

func (resolver *Resolver) Listen() error {

	var err error = nil

	connection, err1 := net.ListenUDP("udp", &net.UDPAddr{
		Port: int(resolver.Port),
		IP:   net.ParseIP(resolver.Host),
	})

	if err1 == nil {

		for {

			buffer := make([]byte, 512)
			length, remote, err := connection.ReadFromUDP(buffer)

			if err == nil {

				packet := dns.Parse(buffer[0:length])
				buffer = make([]byte, 512)

				if packet.Type == "query" && len(packet.Questions) > 0 {

					response := resolver.ResolvePacket(packet)

					if response.Type == "response" {
						resolver.Cache.Write(response)
						connection.WriteTo(response.Bytes(), remote)
					}

				}

			}

		}

	} else {
		err = err1
	}

	return err

}

func (resolver *Resolver) SetTunnel(value interfaces.Tunnel) {
	resolver.Tunnel = value
}
