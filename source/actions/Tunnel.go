package actions

import "tholian-endpoint/types"
import "tholian-warps/console"
import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
// import "tholian-warps/protocols/https"
import "tholian-warps/structs"

func Tunnel(folder string, host string, port uint16, protocol types.Protocol) {

	console.Group("actions/Tunnel")

	web_cache := structs.NewWebCache(folder + "/proxy")
	dns_cache := structs.NewDomainCache(folder + "/resolver")

	resolver := dns.NewResolver("localhost", 8053, &dns_cache)

	http_proxy := http.NewProxy("localhost", 8080, &web_cache)
	http_proxy.SetResolver(&resolver)

	// TODO
	// https_proxy := https.NewProxy("localhost", 8443, &web_cache)
	// https_proxy.SetResolver(&resolver)

	if protocol == types.ProtocolDNS {

		// TODO
		// tunnel := dns.NewTunnel(host, port)

		// resolver.SetTunnel(&tunnel)
		// http_proxy.SetTunnel(&tunnel)
		// https_proxy.SetTunnel(&tunnel)

	} else if protocol == types.ProtocolHTTP {

		tunnel := http.NewTunnel(host, port)

		resolver.SetTunnel(&tunnel)
		http_proxy.SetTunnel(&tunnel)
		// https_proxy.SetTunnel(&tunnel)

	} else if protocol == types.ProtocolHTTPS {

		// TODO
		// tunnel := https.NewTunnel(host, port)

		// resolver.SetTunnel(&tunnel)
		// http_proxy.SetTunnel(&tunnel)
		// https_proxy.SetTunnel(&tunnel)

	}

	console.Log("Listening on dns://localhost:8053")
	console.Log("Listening on http://localhost:8080")
	console.Log("Listening on https://localhost:8443")

	go func() {

		err := resolver.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	}()

	// go func() {

	// 	err := https_proxy.Listen()

	// 	if err != nil {
	// 		console.Error(err.Error())
	// 	}

	// }()

	err := http_proxy.Listen()

	if err != nil {
		console.Error(err.Error())
	}

	console.GroupEnd("actions/Tunnel")

}
