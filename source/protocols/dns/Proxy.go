package dns

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-warps/console"
import "tholian-warps/interfaces"
import dns_tunnel "tholian-warps/protocols/dns/tunnel"
import utils_http "tholian-warps/utils/protocols/http"
import "net"
import "strconv"
import "strings"

type Proxy struct {
	Host     string                `json:"host"`
	Port     uint16                `json:"port"`
	Cache    interfaces.ProxyCache `json:"cache"`
	Tunnel   interfaces.Tunnel     `json:"tunnel"`
	Resolver interfaces.Resolver   `json:"resolver"`
	listener *net.UDPConn
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

	if dns_tunnel.IsRangeRequest(&query) {

		request_url := dns_tunnel.DecodeURL(&query)
		request_from, request_to, _ := dns_tunnel.DecodeContentRange(&query)

		if request_url != nil {

			if request_from == 0 && request_to == -1 {

				http_request := http.NewPacket()
				http_request.SetURL(*request_url)
				http_request.SetMethod(http.MethodGet)
				http_request.SetHeader("Range", "bytes=0-")

				http_response := proxy.RequestPacket(http_request)

				if http_response.Type == "response" {

					http_response.Decode()

					response = dns.NewPacket()
					response.SetType("response")
					response.SetIdentifier(query.Identifier)
					response.SetResponseCode(dns.ResponseCodeNoError)
					response.Flags.RecursionAvailable = true

					if len(http_response.Payload) >= 512 {
						dns_tunnel.EncodeFirstResponse(&response, request_url, http_response.Headers, http_response.Payload)
					} else {
						dns_tunnel.EncodeFirstResponse(&response, request_url, http_response.Headers, http_response.Payload)
					}

				} else {

					response := dns.NewPacket()
					response.SetType("response")
					response.SetIdentifier(query.Identifier)
					dns_tunnel.EncodeErrorResponse(&response, request_url, request_from, request_to)

				}

			} else if request_from > 0 && request_to > request_from {

				http_request := http.NewPacket()
				http_request.SetURL(*request_url)
				http_request.SetMethod(http.MethodGet)
				http_request.SetHeader("Range", "bytes=" + strconv.Itoa(request_from) + "-" + strconv.Itoa(request_to))

				http_response := proxy.RequestPacket(http_request)

				if http_response.Type == "response" {

					http_response.Decode()

					if http_response.Status == http.StatusPartialContent {

						partial_from, partial_to, partial_size := utils_http.ParseContentRange(http_response.GetHeader("Content-Range"))

						if partial_from == request_from && partial_to == request_to && partial_size >= request_to+1 {

							response = dns.NewPacket()
							response.SetType("response")
							response.SetIdentifier(query.Identifier)

							dns_tunnel.EncodeFrameResponse(&response, request_url, http_response.Headers, http_response.Payload, partial_from, partial_to, partial_size)

						} else {

							response = dns.NewPacket()
							response.SetType("response")
							response.SetIdentifier(query.Identifier)

							dns_tunnel.EncodeErrorResponse(&response, request_url, request_from, request_to)

						}

					} else if http_response.Status == http.StatusOK {

						if len(http_response.Payload) >= request_to+1 {

							response = dns.NewPacket()
							response.SetType("response")
							response.SetIdentifier(query.Identifier)

							dns_tunnel.EncodeFrameResponse(&response, request_url, http_response.Headers, http_response.Payload[request_from:request_to+1], request_from, request_to, len(http_response.Payload))

						} else {

							response = dns.NewPacket()
							response.SetType("response")
							response.SetIdentifier(query.Identifier)

							dns_tunnel.EncodeErrorResponse(&response, request_url, request_from, request_to)

						}

					}

				} else {

					response = dns.NewPacket()
					response.SetType("response")
					response.SetIdentifier(query.Identifier)

					dns_tunnel.EncodeErrorResponse(&response, request_url, request_from, request_to)

				}

			} else {

				response = dns.NewPacket()
				response.SetType("response")
				response.SetIdentifier(query.Identifier)

				dns_tunnel.EncodeErrorResponse(&response, request_url, request_from, request_to)

			}

		} else {

			response = dns.NewPacket()
			response.SetType("response")
			response.SetIdentifier(query.Identifier)

		}

	} else {

		if proxy.Tunnel != nil {
			response = proxy.Tunnel.ResolvePacket(query)
		} else if proxy.Resolver != nil {
			response = proxy.Resolver.ResolvePacket(query)
		} else {

			tmp, err := dns.ResolvePacket(query)

			if err == nil {
				response = tmp
			}

		}

	}

	return response

}

func (proxy *Proxy) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	if proxy.Cache != nil && proxy.Cache.Exists(request) {

		response = proxy.Cache.Read(request)

	} else if proxy.Tunnel != nil {

		response = proxy.Tunnel.RequestPacket(request)

	} else {

		if proxy.Resolver != nil {

			request.SetResolveMethod(func(domain string) (dns.Packet, error) {
				return proxy.Resolver.Resolve(domain), nil
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
			response.SetStatus(http.StatusNotFound)
			response.SetPayload([]byte{})

		}

	}

	return response

}

func (proxy *Proxy) SetResolver(value interfaces.Resolver) {
	proxy.Resolver = value
}

func (proxy *Proxy) SetTunnel(value interfaces.Tunnel) {
	proxy.Tunnel = value
}

func (proxy *Proxy) Destroy() error {

	var err error = nil

	if proxy.listener != nil {
		err = proxy.listener.Close()
		proxy.listener = nil
	}

	return err

}

func (proxy *Proxy) Listen() error {

	var err error = nil

	listener, err1 := net.ListenUDP("udp", &net.UDPAddr{
		Port: int(proxy.Port),
		IP:   net.ParseIP(proxy.Host),
	})

	if err1 == nil {

		proxy.listener = listener

		for {

			if proxy.listener != nil {

				buffer := make([]byte, 1232)
				length, remote, err2 := proxy.listener.ReadFromUDP(buffer)

				if err2 == nil {

					packet := dns.Parse(buffer[0:length])
					buffer = make([]byte, 1232)

					if packet.Type == "query" && len(packet.Questions) > 0 {

						go func(remote net.Addr, packet dns.Packet) {

							response := proxy.ResolvePacket(packet)

							if response.Type == "response" {

								if proxy.listener != nil {
									proxy.listener.WriteTo(response.Bytes(), remote)
								}

							} else {

								console.Error("Cannot resolve packet")
								console.Inspect(packet)

							}

						}(remote, packet)

					}

				} else {

					str := err2.Error()

					if strings.HasSuffix(str, "use of closed network connection") {
						break
					}

				}

			} else {
				break
			}

		}

	} else {
		err = err1
	}

	return err

}

