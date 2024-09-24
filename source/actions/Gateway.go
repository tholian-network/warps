package actions

import "tholian-endpoint/types"
import "tholian-warps/certificates"
import "tholian-warps/console"
import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/protocols/https"
import "tholian-warps/protocols/socks"
import "tholian-warps/structs"
import "tholian-warps/utils/arguments"

func Gateway(folder string, listen *arguments.Config) {

	console.Group("actions/Gateway")

	if listen.Protocol == types.ProtocolANY {

		web_cache := structs.NewProxyCache(folder + "/proxy")
		dns_cache := structs.NewResolverCache(folder + "/resolver")

		resolver := dns.NewResolver("127.0.0.1", 53535, &dns_cache)
		dns_proxy := dns.NewProxy(listen.Host, 1053, &web_cache)
		dns_proxy.SetResolver(&resolver)

		http_proxy := http.NewProxy(listen.Host, 1080, &web_cache)
		http_proxy.SetResolver(&resolver)

		https_proxy := https.NewProxy(listen.Host, 1443, &web_cache)
		https_proxy.SetCertificate(certificates.Proxy)
		https_proxy.SetResolver(&resolver)

		socks_proxy := socks.NewProxy(listen.Host, 1090, &web_cache)
		socks_proxy.SetResolver(&resolver)

		console.Log("Listening on dns://" + listen.Host + ":1053")
		console.Log("Listening on http://" + listen.Host + ":1080")
		console.Log("Listening on https://" + listen.Host + ":1443")
		console.Log("Listening on socks://" + listen.Host + ":1090")

		go func() {

			err := resolver.Listen()

			if err != nil {
				console.Error(err.Error())
			}

		}()

		go func() {

			err := dns_proxy.Listen()

			if err != nil {
				console.Error(err.Error())
			}

		}()

		go func() {

			err := http_proxy.Listen()

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

		err := socks_proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	} else if listen.Protocol == types.ProtocolDNS {

		web_cache := structs.NewProxyCache(folder + "/proxy")
		dns_cache := structs.NewResolverCache(folder + "/resolver")

		resolver := dns.NewResolver("127.0.0.1", 53535, &dns_cache)
		proxy := dns.NewProxy(listen.Host, listen.Port, &web_cache)
		proxy.SetResolver(&resolver)

		console.Log("Listening on " + listen.String())

		err := proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	} else if listen.Protocol == types.ProtocolHTTP {

		web_cache := structs.NewProxyCache(folder + "/proxy")
		dns_cache := structs.NewResolverCache(folder + "/resolver")

		resolver := dns.NewResolver("127.0.0.1", 53535, &dns_cache)
		proxy := http.NewProxy(listen.Host, listen.Port, &web_cache)
		proxy.SetResolver(&resolver)

		console.Log("Listening on " + listen.String())

		err := proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	} else if listen.Protocol == types.ProtocolHTTPS {

		web_cache := structs.NewProxyCache(folder + "/proxy")
		dns_cache := structs.NewResolverCache(folder + "/resolver")

		resolver := dns.NewResolver("127.0.0.1", 53535, &dns_cache)
		proxy := https.NewProxy(listen.Host, listen.Port, &web_cache)
		proxy.SetCertificate(certificates.Proxy)
		proxy.SetResolver(&resolver)

		console.Log("Listening on " + listen.String())

		err := proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	} else if listen.Protocol == types.ProtocolSOCKS {

		web_cache := structs.NewProxyCache(folder + "/proxy")
		dns_cache := structs.NewResolverCache(folder + "/resolver")

		resolver := dns.NewResolver("127.0.0.1", 53535, &dns_cache)
		proxy := socks.NewProxy(listen.Host, listen.Port, &web_cache)
		proxy.SetResolver(&resolver)

		console.Log("Listening on " + listen.String())

		err := proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	}

	console.GroupEnd("actions/Gateway")

}

