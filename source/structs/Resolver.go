package structs

import "tholian-warps/types"
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

func (resolver *DomainResolver) Listen() bool {

	var result bool = false

	// TODO: DNS listener
	// TODO: Route DNS requests through tunnel if Tunnel != nil

	return result

}

func (resolver *DomainResolver) SetProtocol(protocol types.Protocol) {

	if protocol == types.ProtocolHTTPS {
		resolver.Protocol = protocol
	} else if protocol == types.ProtocolDNS {
		resolver.Protocol = protocol
	}

}

