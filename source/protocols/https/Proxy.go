package https

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-warps/console"
import "tholian-warps/interfaces"
import http_tunnel "tholian-warps/protocols/http/tunnel"
import utils_net "tholian-warps/utils/net"
import "crypto/tls"
import "net"
import "strconv"
import "strings"

type Proxy struct {
	Host        string                `json:"host"`
	Port        uint16                `json:"port"`
	Cache       interfaces.ProxyCache `json:"cache"`
	Tunnel      interfaces.Tunnel     `json:"tunnel"`
	Resolver    interfaces.Resolver   `json:"resolver"`
	certificate *tls.Certificate
	listener    net.Listener
}

func NewProxy(host string, port uint16, cache interfaces.ProxyCache) Proxy {

	var proxy Proxy

	if host == "*" || host == "any" || host == "localhost" || host == "127.0.0.1" {
		proxy.Host = "0.0.0.0"
	} else if strings.ToLower(host) == host {
		proxy.Host = host
	} else {
		proxy.Host = "0.0.0.0"
	}

	proxy.Port = port
	proxy.Cache = cache

	return proxy

}

func (proxy *Proxy) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	if proxy.Resolver != nil {
		response = proxy.Resolver.ResolvePacket(query)
	} else if proxy.Tunnel != nil {
		response = proxy.Tunnel.ResolvePacket(query)
	} else {
		response = dns.ResolvePacket(query)
	}

	return response

}

func (proxy *Proxy) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	if http_tunnel.IsResolveRequest(&request) {

		dns_query := dns.Parse(http_tunnel.DecodePayload(&request))

		if dns_query.Type == "query" {

			dns_response := proxy.ResolvePacket(dns_query)

			if dns_response.Type == "response" {

				response = http.NewPacket()
				response.SetURL(*request.URL)

				http_tunnel.EncodePayload(&response, dns_response.Bytes())

			} else {

				response = http.NewPacket()
				response.SetURL(*request.URL)

				http_tunnel.EncodeError(&dns_query, &response, http.StatusNotFound)

			}

		} else {

			response = http.NewPacket()
			response.SetURL(*request.URL)

			http_tunnel.EncodeError(&dns_query, &response, http.StatusNotFound)

		}

	} else {

		if proxy.Cache != nil && proxy.Cache.Exists(request) {

			response = proxy.Cache.Read(request)

		} else if proxy.Tunnel != nil {

			response = proxy.Tunnel.RequestPacket(request)

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

				tmp := http.RequestPacket(request)

				if tmp.Type == "response" {

					response = tmp

				} else {

					response = http.NewPacket()
					response.SetURL(*request.URL)
					response.SetStatus(http.StatusNotFound)
					response.SetPayload([]byte{})

				}

			} else {

				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusRequestTimeout)
				response.SetPayload([]byte{})

			}

		}

	}

	return response

}

func (proxy *Proxy) SetCertificate(value tls.Certificate) {
	proxy.certificate = &value
}

func (proxy *Proxy) SetResolver(value interfaces.Resolver) {
	proxy.Resolver = value
}

func (proxy *Proxy) SetTunnel(value interfaces.Tunnel) {
	proxy.Tunnel = value
}

func (proxy *Proxy) Listen() error {

	var err error = nil

	host := "0.0.0.0:" + strconv.FormatUint(uint64(proxy.Port), 10)

	listener, err1 := tls.Listen("tcp", host, &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{*proxy.certificate},
		MaxVersion:         tls.VersionTLS12,
	})

	if err1 == nil {

		proxy.listener = listener

		for {

			connection, err2 := listener.Accept()

			if err2 == nil {

				buffer := utils_net.ReadConnection(connection)

				if len(buffer) > 0 {

					packet := http.Parse(buffer)

					if packet.Method == http.MethodConnect {

						connection.Write([]byte("HTTP/1.1 401 Unauthorized\r\n\r\n"))
						connection.Close()

					} else if packet.Method.String() != "" {

						go func(connection net.Conn, packet http.Packet) {

							response := proxy.RequestPacket(packet)

							if response.Type == "response" {

								if proxy.Cache != nil {
									proxy.Cache.Write(response)
								}

								connection.Write(response.Bytes())
								connection.Close()

							} else {

								response := http.NewPacket()
								response.SetStatus(http.StatusInternalServerError)

								connection.Write(response.Bytes())
								connection.Close()

							}

						}(connection, packet)

					} else {
						defer connection.Close()
					}

				} else {
					defer connection.Close()
				}

			} else {
				console.Error(err2.Error())
			}

		}

	} else {
		err = err1
	}

	return err

}
