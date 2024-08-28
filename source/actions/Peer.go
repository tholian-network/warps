package actions

import "tholian-warps/structs"
import "tholian-warps/types"
import "os"

func Peer(host string) {

	pwd, err := os.Getwd()

	if err == nil {

		web_cache := structs.NewWebCache(pwd + "/tholian-warps/proxy")
		domain_cache := structs.NewDomainCache(pwd + "/tholian-warps/resolver")

		resolver := structs.NewResolver("localhost", 8053, &domain_cache, nil)
		http_proxy := structs.NewProxy("localhost", 8080, &web_cache, nil, types.ProtocolHTTP)
		https_proxy := structs.NewProxy("localhost", 8181, &web_cache, nil, types.ProtocolHTTPS)

		go resolver.Listen()
		go https_proxy.Listen()

		http_proxy.Listen()

	}

}
