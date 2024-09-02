package actions

import "tholian-warps/console"
import "tholian-warps/structs"
import "tholian-warps/types"

func Tunnel(folder string, host string, port uint16, protocol types.Protocol) {

	console.Group("actions/Tunnel")

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

	http_proxy.Resolver = &resolver
	https_proxy.Resolver = &resolver

	console.Log("Listening on dns://localhost:8053")
	console.Log("Listening on http://localhost:8080")
	console.Log("Listening on https://localhost:8181")

	go func() {

		err := resolver.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	}()

	go func() {

		err := https_proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	}()

	err := http_proxy.Listen()

	if err != nil {
		console.Error(err.Error())
	}

	console.GroupEnd("actions/Tunnel")

}
