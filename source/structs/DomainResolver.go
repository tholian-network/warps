package structs

import "tholian-endpoint/protocols/dns"
import endpoint_types "tholian-endpoint/types"
import "tholian-warps/types"
import "net"
import "strings"

type DomainResolver struct {
	Host     string         `json:"host"`
	Port     uint16         `json:"port"`
	Protocol types.Protocol `json:"protocol"`
	Cache    *DomainCache   `json:"cache"`
	Tunnel   *types.Tunnel  `json:"tunnel"`
}

func NewResolver(host string, port uint16, cache *DomainCache, tunnel *types.Tunnel) DomainResolver {

	var resolver DomainResolver

	if host == "*" || host == "any" || host == "localhost" {
		resolver.Host = "0.0.0.0"
	} else if strings.ToLower(host) == host {
		resolver.Host = host
	} else {
		resolver.Host = "0.0.0.0"
	}

	if port > 0 && port < 65535 {
		resolver.Port = port
	} else {
		resolver.Port = 5353
	}

	if cache != nil {
		resolver.Cache = cache
	} else {

		cache := NewDomainCache("/tmp/tholian-warps/resolver")
		resolver.Cache = &cache

	}

	resolver.Protocol = types.ProtocolDNS

	if tunnel != nil {
		resolver.Tunnel = tunnel
	}

	return resolver

}

func (resolver *DomainResolver) Resolve(domain string) dns.Packet {

	var response dns.Packet

	if endpoint_types.IsDomain(domain) {

		query := dns.NewPacket()
		query.SetType("query")
		query.AddQuestion(dns.NewQuestion(domain, dns.TypeA))
		query.AddQuestion(dns.NewQuestion(domain, dns.TypeAAAA))

		response = resolver.ResolvePacket(query)

	}

	return response

}

func (resolver *DomainResolver) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	if resolver.Cache.Exists(query) {

		response = resolver.Cache.Read(query)

	} else {

		if resolver.Tunnel != nil {

			if resolver.Tunnel.Protocol == types.ProtocolHTTPS {

				// TODO: Tunnel request through HTTPS

			} else if resolver.Tunnel.Protocol == types.ProtocolHTTP {

				// TODO: Tunnel request through HTTP

			} else if resolver.Tunnel.Protocol == types.ProtocolDNS {

				// TODO: Tunnel request through DNS

			} else if resolver.Tunnel.Protocol == types.ProtocolICMP {

				// TODO: Tunnel request through ICMP

			}

		} else {

			response = dns.ResolvePacket(query)

		}

	}


	return response

}

func (resolver *DomainResolver) Listen() error {

	var err error = nil

	if resolver.Protocol == types.ProtocolHTTPS {

		// TODO: Parse HTTP / JSON payload

	} else if resolver.Protocol == types.ProtocolHTTP {

		// TODO: Parse HTTP / JSON payload

	} else if resolver.Protocol == types.ProtocolDNS {

		connection, err1 := net.ListenUDP("udp", &net.UDPAddr{
			Port: int(resolver.Port),
			IP:   net.ParseIP(resolver.Host),
		})

		if err1 == nil {

			defer connection.Close()

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

	}

	return err

}

func (resolver *DomainResolver) SetProtocol(protocol types.Protocol) {

	if protocol == types.ProtocolHTTPS {
		resolver.Protocol = protocol
	} else if protocol == types.ProtocolHTTP {
		resolver.Protocol = protocol
	} else if protocol == types.ProtocolDNS {
		resolver.Protocol = protocol
	}

}

