package actions

import "tholian-warps/console"
import "tholian-warps/structs"
import "tholian-warps/types"

func Peer(folder string, host string) {

	console.Group("actions/Peer")

	web_cache := structs.NewWebCache(folder + "/proxy")
	domain_cache := structs.NewDomainCache(folder + "/resolver")

	resolver := structs.NewResolver("localhost", 8053, &domain_cache, nil)
	http_proxy := structs.NewProxy("localhost", 8080, &web_cache, nil, types.ProtocolHTTP)
	https_proxy := structs.NewProxy("localhost", 8443, &web_cache, nil, types.ProtocolHTTPS)

	http_proxy.SetResolver(&resolver)
	https_proxy.SetResolver(&resolver)

	console.Log("Listening on dns://localhost:8053")
	console.Log("Listening on http://localhost:8080")
	console.Log("Listening on https://localhost:8443")

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

	console.GroupEnd("actions/Peer")

}
