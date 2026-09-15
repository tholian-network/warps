package icmptunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/interfaces"
import "tholian-warps/types"

type Resolver struct {
	Host   string                   `json:"host"`
	Port   uint16                   `json:"port"`
	Cache  interfaces.ResolverCache `json:"cache"`
	Tunnel interfaces.Tunnel        `json:"tunnel"`
}

func NewResolver(host string, port uint16, cache interfaces.ResolverCache) Resolver {

	var resolver Resolver

	resolver.Host = host
	resolver.Port = port
	resolver.Cache = cache
	resolver.Tunnel = nil

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

		tmp, err := dns.ResolvePacket(query)

		if err == nil && tmp.Type == "response" {

			response = tmp

		} else {

			response = dns.NewPacket()
			response.SetType("response")
			response.SetIdentifier(query.Identifier)
			response.SetResponseCode(dns.ResponseCodeNonExistDomain)
			response.Flags.RecursionAvailable = true

			for q := 0; q < len(query.Questions); q++ {
				response.AddQuestion(query.Questions[q])
			}

		}

	}

	return response

}

func (resolver *Resolver) SetTunnel(value interfaces.Tunnel) {
	resolver.Tunnel = value
}
