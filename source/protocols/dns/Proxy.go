package dns

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-warps/interfaces"
import dns_tunnel "tholian-warps/protocols/dns/tunnel"
import "net"
import net_url "net/url"
import "strings"

type Proxy struct {
	Host     string                `json:"host"`
	Port     uint16                `json:"port"`
	Cache    interfaces.ProxyCache `json:"cache"`
	Tunnel   interfaces.Tunnel     `json:"tunnel"`
	Resolver interfaces.Resolver   `json:"resolver"`
}

func NewProxy(host string, port uint16, cache interfaces.ProxyCache) Proxy {

	var proxy Proxy

	if host == "*" || host == "any" || host == "localhost" {
		proxy.Host = "0.0.0.0"
	} else if host == "127.0.0.1" {
		proxy.Host = "0.0.0.0"
	} else if strings.ToLower(host) == host {
		proxy.Host = host
	} else {
		proxy.Host = "0.0.0.0"
	}

	if port > 0 && port < 65535 {
		proxy.Port = port
	} else {
		proxy.Port = 8080
	}

	if cache != nil {
		proxy.Cache = cache
	}

	return proxy

}

func (proxy *Proxy) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	if dns_tunnel.IsRangeRequest(&query) {

		request_url, request_from, request_to, _ := dns_tunnel.DecodeContentRange(&query)

		if request_url != "" {

			if request_from == 0 && request_to == -1 {

				url, _ := net_url.Parse(request_url)

				http_request := http.NewPacket()
				http_request.SetURL(*url)
				http_request.SetMethod(http.MethodGet)
				http_request.SetHeader("Range", "bytes=0-")

				http_response := proxy.RequestPacket(http_request)

				if http_response.Type == "response" {

					http_response.Decode()

					response = dns.NewPacket()
					response.SetType("response")
					response.SetIdentifier(query.Identifier)
					response.SetResponseCode(dns.ResponseCodeNoError)

					dns_tunnel.EncodeContentRange(&response, request_url, 0, 1023, len(http_response.Payload))
					dns_tunnel.EncodeHeaders(&response, http_response.Headers)
					dns_tunnel.EncodePayload(&response, http_response.Payload[0:1024])

				} else {

					response = dns.NewPacket()
					response.SetType("response")
					response.SetIdentifier(query.Identifier)
					response.SetResponseCode(dns.ResponseCodeNonExistDomain)

					dns_tunnel.EncodeContentRange(&response, request_url, 0, 0, 0)
					dns_tunnel.EncodePayload(&response, []byte{})

				}

			} else if request_from > 0 && request_to > request_from {

				url, _ := net_url.Parse(request_url)

				http_request := http.NewPacket()
				http_request.SetURL(*url)
				http_request.SetMethod(http.MethodGet)
				http_request.SetHeader("Range", "bytes=0-")

				http_response := proxy.RequestPacket(http_request)

				if http_response.Type == "response" {

					http_response.Decode()

					response = dns.NewPacket()
					response.SetType("response")
					response.SetIdentifier(query.Identifier)
					response.SetResponseCode(dns.ResponseCodeNoError)

					dns_tunnel.EncodeContentRange(&response, request_url, request_from, request_to, len(http_response.Payload))
					dns_tunnel.EncodePayload(&response, http_response.Payload[request_from:request_to+1])

				} else {

					response = dns.NewPacket()
					response.SetType("response")
					response.SetIdentifier(query.Identifier)
					response.SetResponseCode(dns.ResponseCodeNonExistDomain)

					dns_tunnel.EncodeContentRange(&response, request_url, 0, 0, 0)
					dns_tunnel.EncodePayload(&response, []byte{})

				}

			} else {

				response = dns.NewPacket()
				response.SetType("response")
				response.SetIdentifier(query.Identifier)
				response.SetResponseCode(dns.ResponseCodeNonExistDomain)

				dns_tunnel.EncodeContentRange(&response, request_url, 0, 0, 0)
				dns_tunnel.EncodePayload(&response, []byte{})

			}

		} else {

			response = dns.NewPacket()
			response.SetType("response")
			response.SetIdentifier(query.Identifier)
			response.SetResponseCode(dns.ResponseCodeNonExistDomain)

			dns_tunnel.EncodeContentRange(&response, request_url, 0, 0, 0)
			dns_tunnel.EncodePayload(&response, []byte{})

		}

	} else {

		if proxy.Tunnel != nil {
			response = proxy.Tunnel.ResolvePacket(query)
		} else if proxy.Resolver != nil {
			response = proxy.Resolver.ResolvePacket(query)
		} else {
			response = dns.ResolvePacket(query)
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

func (proxy *Proxy) Listen() error {

	var err error = nil

	listener, err1 := net.ListenUDP("udp", &net.UDPAddr{
		Port: int(proxy.Port),
		IP:   net.ParseIP(proxy.Host),
	})

	if err1 == nil {

		for {

			buffer := make([]byte, 1232)
			length, remote, err := listener.ReadFromUDP(buffer)

			if err == nil {

				packet := dns.Parse(buffer[0:length])
				buffer = make([]byte, 1232)

				if packet.Type == "query" && len(packet.Questions) > 0 {

					response := proxy.ResolvePacket(packet)

					if response.Type == "response" {
						listener.WriteTo(response.Bytes(), remote)
					}

				}

			}

		}

	} else {
		err = err1
	}

	return err

}

