package structs

import "tholian-warps/types"
import "fmt"
import "log"
import "net/http"
import "strconv"
import "strings"

type Proxy struct {
	Host     string         `json:"host"`
	Port     uint16         `json:"port"`
	Protocol types.Protocol `json:"protocol"`
	Cache    *WebCache      `json:"cache"`
	Tunnel   *types.Tunnel  `json:"tunnel"`
}

func NewProxy(host string, port uint16, cache *WebCache, tunnel *types.Tunnel, protocol types.Protocol) Proxy {

	var proxy Proxy

	if host == "*" || host == "any" || host == "localhost" {
		proxy.Host = "localhost"
	} else if strings.ToLower(host) == host {
		proxy.Host = host
	} else {
		proxy.Host = "localhost"
	}

	if port > 0 && port < 65535 {

		proxy.Port = port

	} else {

		if protocol == types.ProtocolHTTP {
			proxy.Port = 8080
		} else if protocol == types.ProtocolHTTPS {
			proxy.Port = 8181
		} else if protocol == types.ProtocolDNS {
			proxy.Port = 8053
		}

	}

	if cache != nil {
		proxy.Cache = cache
	} else {

		cache := NewWebCache("/tmp/tholian-warps/proxy")
		proxy.Cache = &cache

	}

	proxy.Protocol = protocol

	return proxy

}

// func (proxy *Proxy) ServeDNS(response dns.ResponseWriter, request *dns.Request) {
// }

func (proxy *Proxy) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	if proxy.Tunnel != nil {

		if proxy.Tunnel.Protocol == types.ProtocolDNS {

			// TODO: Send request via DNS packet

		} else if proxy.Tunnel.Protocol == types.ProtocolHTTP {

			// TODO: Send HTTP request to other proxy

		} else if proxy.Tunnel.Protocol == types.ProtocolHTTPS {

			// TODO: Send HTTPS request to other proxy

		}

		fmt.Println("TODO: Proxy " + request.URL.String() + " through " + proxy.Tunnel.Host + ":" + strconv.FormatUint(uint64(proxy.Tunnel.Port), 10))
		fmt.Println(request.URL)

		// TODO: create request for tunnel protocol
		// TODO: do request, wait for response
		// TODO: after response, respond with HTTP response

	} else {

		// TODO: Send HTTP request
		// TODO: Send HTTPS request

	}


	response.WriteHeader(http.StatusNotFound)

}

func (proxy *Proxy) Listen() bool {

	var result bool = false

	if proxy.Protocol == types.ProtocolDNS {

		// TODO: DNS listener

	} else if proxy.Protocol == types.ProtocolHTTP {

		host := ""
		port := strconv.FormatUint(uint64(proxy.Port), 10)

		if proxy.Host != "localhost" && proxy.Host != "127.0.0.1" && proxy.Host != "0.0.0.0" {
			host = proxy.Host
		}

		err := http.ListenAndServe(host + ":" + port, proxy)

		if err != nil {
			log.Fatal("Port " + port + " already in use!")
		} else {
			result = true
		}

	} else if proxy.Protocol == types.ProtocolHTTPS {

		// TODO: Use ServeHTTPS

	}


	return result

}
