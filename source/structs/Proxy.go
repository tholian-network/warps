package structs

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-warps/console"
import "tholian-warps/types"
import utils_net "tholian-warps/utils/net"
import "crypto/tls"
// import "io"
import "net"
import "strconv"
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

		console.Inspect(request)

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
			console.Inspect(data)

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

			for {

				connection, err2 := listener.Accept()

				if err2 == nil {

					buffer := utils_net.ReadConnection(connection)

					if len(buffer) > 0 {

						packet := http.Parse(buffer)

						if packet.Method == http.MethodConnect {

							console.Warn("CONNECT")
							console.Inspect(packet)

							// connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
							connection.Write([]byte("HTTP/1.1 401 Unauthorized\r\n\r\n"))

							// TODO: Remove this
							connection.Close()

							// TODO: This is wrong, because CONNECT should directly connect and not delegate via TLS connection

							// proxied_connection, err1 := net.DialTCP("tcp", nil, &net.TCPAddr{
							// 	IP:   net.ParseIP("127.0.0.1"),
							// 	Port: int(8443),
							// })

							// if err1 == nil {

							// 	channel := make(chan error, 1)

							// 	io_copy := func(read net.Conn, write net.Conn) {
							// 		_, err := io.Copy(read, write)
							// 		channel <- err
							// 	}

							// 	go io_copy(connection, proxied_connection)
							// 	go io_copy(proxied_connection, connection)

							// 	err1 := <-channel
							// 	err2 := <-channel

							// 	if err1 != nil || err2 != nil {
							// 		break
							// 	}

							// 	defer connection.Close()
							// 	defer proxied_connection.Close()

							// }

						} else if packet.Method.String() != "" {

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

					defer connection.Close()

				}

			}

		} else {
			err = err1
		}

	} else if proxy.Protocol == types.ProtocolHTTPS {

		host := "0.0.0.0:8443"

		if proxy.Host != "localhost" && proxy.Host != "127.0.0.1" && proxy.Host != "0.0.0.0" {
			host = proxy.Host + ":" + strconv.Itoa(int(proxy.Port))
		}

		listener, err1 := tls.Listen("tcp", host, &tls.Config{
			Certificates: []tls.Certificate{*Certificate},
			MaxVersion:   tls.VersionTLS12,
		})

		if err1 == nil {

			for {

				connection, err2 := listener.Accept()

				if err2 == nil {

					buffer := utils_net.ReadConnection(connection)

					if len(buffer) > 0 {

						packet := http.Parse(buffer)

						if packet.Type == "request" && packet.Method.String() != "" {

							if packet.URL.Scheme == "" && packet.URL.Host == "" {

								hostname, _ := packet.Headers["Host"]

								if hostname != "" {
									packet.URL.Host = hostname
									packet.URL.Scheme = "https"
								}

							}

							response := proxy.RequestPacket(packet)

							if response.Type == "response" {

								proxy.Cache.Write(response)
								connection.Write(response.Bytes())

							} else {

								response := http.NewPacket()
								response.SetStatus(http.StatusInternalServerError)

								connection.Write(response.Bytes())

							}

						}

					}

					defer connection.Close()

				} else {
					defer connection.Close()
				}

			}

		} else {
			err = err1
		}

	}

	return err

}

func (proxy *Proxy) SetResolver(resolver *DomainResolver) {
	proxy.Resolver = resolver
}
