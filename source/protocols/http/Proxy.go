package http

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-warps/console"
import "tholian-warps/interfaces"
import "tholian-warps/structs"
import utils_net "tholian-warps/utils/net"
import "encoding/base64"
import "net"
import "strings"

type Proxy struct {
	Host     string              `json:"host"`
	Port     uint16              `json:"port"`
	Cache    *structs.WebCache   `json:"cache"`
	Tunnel   interfaces.Tunnel   `json:"tunnel"`
	Resolver interfaces.Resolver `json:"resolver"`
}

func NewProxy(host string, port uint16, cache *structs.WebCache) Proxy {

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

	if proxy.Tunnel != nil {
		response = proxy.Tunnel.ResolvePacket(query)
	} else if proxy.Resolver != nil {
		response = proxy.Resolver.ResolvePacket(query)
	} else {
		response = dns.ResolvePacket(query)
	}

	return response

}

func (proxy *Proxy) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	if proxy.Tunnel != nil {

		response = proxy.Tunnel.RequestPacket(request)

	} else {

		if request.GetHeader("Content-Type") == "application/dns-message" {

			dns_request_payload := []byte{}

			if request.Method == http.MethodGet {

				tmp1 := request.URL.Query()

				if tmp1.Get("dns") != "" {

					tmp2, err2 := base64.URLEncoding.DecodeString(tmp1.Get("dns"))

					if err2 == nil {
						dns_request_payload = tmp2
					}

				}

			} else if request.Method == http.MethodPost {

				request.Decode()
				dns_request_payload = request.Payload

			}

			dns_request := dns.Parse(dns_request_payload)

			if dns_request.Type == "query" {

				var dns_response dns.Packet
				var resolved bool = false

				if proxy.Resolver != nil {

					tmp := proxy.Resolver.ResolvePacket(dns_request)

					if tmp.Type == "response" {
						dns_response = tmp
						resolved = true
					} else {
						dns_response = EmptyDNSResponse(dns_request)
					}

				} else {

					tmp := dns.ResolvePacket(dns_request)

					if tmp.Type == "response" {
						dns_response = tmp
						resolved = true
					} else {
						dns_response = EmptyDNSResponse(dns_request)
					}

				}

				response = http.NewPacket()
				response.SetURL(*request.URL)

				if resolved == true {
					response.SetStatus(http.StatusOK)
				} else {
					response.SetStatus(http.StatusNotFound)
				}

				response.SetHeader("Content-Type", "application/dns-message")
				response.SetPayload(dns_response.Bytes())

			} else {

				dns_response := EmptyDNSResponse(dns_request)

				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusNotFound)
				response.SetHeader("Content-Type", "application/dns-message")
				response.SetPayload(dns_response.Bytes())

			}

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
					response = EmptyHTTPResponse(request)
				}

			} else {
				response = EmptyHTTPResponse(request)
			}

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

	listener, err1 := net.ListenTCP("tcp", &net.TCPAddr{
		Port: int(proxy.Port),
		IP:   net.ParseIP(proxy.Host),
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

						// TODO: Remove this
						connection.Write([]byte("HTTP/1.1 401 Unauthorized\r\n\r\n"))
						connection.Close()

						// TODO: This is wrong, because CONNECT should directly connect and not delegate via TLS connection

						// connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

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

	return err

}
