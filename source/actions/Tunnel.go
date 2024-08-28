package actions

import "tholian-warps/structs"
import "tholian-warps/types"
import "os"

func Tunnel(host string, port uint16, protocol types.Protocol) {

	pwd, err := os.Getwd()

	if err == nil {

		web_cache := structs.NewWebCache(pwd + "/tholian-warps/proxy")
		domain_cache := structs.NewDomainCache(pwd + "/tholian-warps/resolver")

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

}
