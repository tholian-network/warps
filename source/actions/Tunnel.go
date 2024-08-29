package actions

import "tholian-warps/structs"
import "tholian-warps/types"

func Tunnel(folder string, host string, port uint16, protocol types.Protocol) {

	web_cache := structs.NewWebCache(folder + "/proxy")
	domain_cache := structs.NewDomainCache(folder + "/resolver")

	tunnel := types.Tunnel{
		Host:     host,
		Port:     port,
		Protocol: protocol,
	}

	resolver := structs.NewResolver("localhost", 8053, &domain_cache, &tunnel)
	http_proxy := structs.NewProxy("localhost", 8080, &web_cache, &tunnel, types.ProtocolHTTP)
	https_proxy := structs.NewProxy("localhost", 8181, &web_cache, &tunnel, types.ProtocolHTTPS)

	go resolver.Listen()
	go https_proxy.Listen()

	http_proxy.Listen()

}
