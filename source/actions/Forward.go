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

func Forward(folder string, listen *arguments.Config, tunnel *arguments.Config) {

	console.Group("actions/Forward")

	if listen.Protocol == types.ProtocolDNS {

		web_cache := structs.NewProxyCache(folder + "/proxy")
		dns_cache := structs.NewResolverCache(folder + "/resolver")

		resolver := dns.NewResolver("127.0.0.1", 53535, &dns_cache)
		proxy := dns.NewProxy(listen.Host, listen.Port, &web_cache)
		proxy.SetResolver(&resolver)

		if tunnel.Protocol == types.ProtocolDNS {

			tmp := dns.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTP {

			tmp := http.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTPS {

			tmp := https.NewTunnel(tunnel.Host, tunnel.Port)
			tmp.SetCertificate(certificates.Proxy)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolSOCKS {

			tmp := socks.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		}

		console.Log("Tunneling to " + tunnel.String())
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

		if tunnel.Protocol == types.ProtocolDNS {

			tmp := dns.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTP {

			tmp := http.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTPS {

			tmp := https.NewTunnel(tunnel.Host, tunnel.Port)
			tmp.SetCertificate(certificates.Proxy)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolSOCKS {

			tmp := socks.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		}

		console.Log("Tunneling to " + tunnel.String())
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

		if tunnel.Protocol == types.ProtocolDNS {

			tmp := dns.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTP {

			tmp := http.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTPS {

			tmp := https.NewTunnel(tunnel.Host, tunnel.Port)
			tmp.SetCertificate(certificates.Proxy)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolSOCKS {

			tmp := socks.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		}

		console.Log("Tunneling to " + tunnel.String())
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

		if tunnel.Protocol == types.ProtocolDNS {

			tmp := dns.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTP {

			tmp := http.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolHTTPS {

			tmp := https.NewTunnel(tunnel.Host, tunnel.Port)
			tmp.SetCertificate(certificates.Proxy)
			proxy.SetTunnel(&tmp)

		} else if tunnel.Protocol == types.ProtocolSOCKS {

			tmp := socks.NewTunnel(tunnel.Host, tunnel.Port)
			proxy.SetTunnel(&tmp)

		}

		console.Log("Tunneling to " + tunnel.String())
		console.Log("Listening on " + listen.String())

		err := proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	}

	console.GroupEnd("actions/Forward")

}

