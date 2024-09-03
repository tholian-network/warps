package structs

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-warps/types"
import "net"
import "strings"

type Proxy struct {
	Host     string          `json:"host"`
	Port     uint16          `json:"port"`
	Protocol types.Protocol  `json:"protocol"`
	Cache    *WebCache       `json:"cache"`
	Tunnel   *types.Tunnel   `json:"tunnel"`
	Resolver *DomainResolver `json:"resolver"`
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

func (proxy *Proxy) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	if proxy.Tunnel != nil {

		if proxy.Tunnel.Protocol == types.ProtocolDNS {

			// TODO: Send request via DNS packet

		} else if proxy.Tunnel.Protocol == types.ProtocolHTTP {

			// TODO: Send HTTP request to other proxy

		} else if proxy.Tunnel.Protocol == types.ProtocolHTTPS {

			// TODO: Send HTTPS request to other proxy

		}

		// TODO: create request for tunnel protocol
		// TODO: do request, wait for response
		// TODO: after response, respond with HTTP response

	} else {

		if proxy.Resolver != nil {

			request.SetResolveMethod(func(domain string) dns.Packet {
				return proxy.Resolver.Resolve(domain)
			})
			request.Resolve()

		} else {
			request.Resolve()
		}

		if request.Server != nil {

			data := http.Request(request)

			if data.Type == "response" {
				response = data
			}

		} else {

			// TODO: Respond with a new packet response
			// TODO: Status 500 Internal Server Error?

		}

	}

	return response

}

func (proxy *Proxy) Listen() error {

	var err error = nil

	if proxy.Protocol == types.ProtocolDNS {

		// TODO: DNS listener

	} else if proxy.Protocol == types.ProtocolHTTP {

		host := "0.0.0.0"

		if proxy.Host != "localhost" && proxy.Host != "127.0.0.1" && proxy.Host != "0.0.0.0" {
			host = proxy.Host
		}

		listener, err1 := net.ListenTCP("tcp", &net.TCPAddr{
			Port: int(proxy.Port),
			IP:   net.ParseIP(host),
		})

		if err1 == nil {

			defer listener.Close()

			for {

				connection, err2 := listener.Accept()

				if err2 == nil {

					buffer := make([]byte, 2048)
					length, err3 := connection.Read(buffer)

					if err3 == nil {

						packet := http.Parse(buffer[0:length])
						buffer = make([]byte, 2048)

						if packet.Method.String() != "" {

							response := proxy.RequestPacket(packet)

							if response.Type == "response" {

								proxy.Cache.Write(response)
								connection.Write(response.Bytes())

								connection.Close()

							} else {

								response := http.NewPacket()
								response.SetStatus(http.StatusInternalServerError)

								connection.Write(response.Bytes())
								connection.Close()

							}

						}

					}

				}

			}

		} else {
			err = err1
		}

	} else if proxy.Protocol == types.ProtocolHTTPS {

		// TODO: Use ServeHTTPS

	}

	return err

}
