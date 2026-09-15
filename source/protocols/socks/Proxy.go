package socks

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/console"
import "tholian-warps/interfaces"
import http_tunnel "tholian-warps/protocols/httptunnel/tunnel"
import "io"
import "net"
import "strconv"
import "strings"

type Proxy struct {
	Host     string                `json:"host"`
	Port     uint16                `json:"port"`
	Cache    interfaces.ProxyCache `json:"cache"`
	Tunnel   interfaces.Tunnel     `json:"tunnel"`
	Resolver interfaces.Resolver   `json:"resolver"`
	listener *net.TCPListener
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
	proxy.listener = nil

	return proxy

}

func (proxy *Proxy) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	if proxy.Resolver != nil {
		response = proxy.Resolver.ResolvePacket(query)
	} else if proxy.Tunnel != nil {
		response = proxy.Tunnel.ResolvePacket(query)
	} else {

		tmp, err := dns.ResolvePacket(query)

		if err == nil {
			response = tmp
		}

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
				response.SetStatus(http.StatusRequestTimeout)
				response.SetPayload([]byte{})

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

	listener, err1 := net.ListenTCP("tcp", &net.TCPAddr{
		Port: int(proxy.Port),
		IP:   net.ParseIP(proxy.Host),
	})

	if err1 == nil {

		proxy.listener = listener

		for {

			connection, err2 := listener.Accept()

			if err2 == nil {

				go proxy.handleConnection(connection)

			} else {

				str := err2.Error()

				if strings.HasSuffix(str, "use of closed network connection") {
					break
				}

				console.Error(str)

			}

		}

	} else {
		err = err1
	}

	return err

}

func (proxy *Proxy) handleConnection(connection net.Conn) {

	defer connection.Close()

	header := make([]byte, 2)

	if _, err := io.ReadFull(connection, header); err != nil {
		return
	}

	if header[0] != 0x05 {
		return
	}

	methods := make([]byte, int(header[1]))

	if _, err := io.ReadFull(connection, methods); err != nil {
		return
	}

	if _, err := connection.Write([]byte{0x05, 0x00}); err != nil {
		return
	}

	request := make([]byte, 4)

	if _, err := io.ReadFull(connection, request); err != nil {
		return
	}

	if request[0] != 0x05 || request[1] != 0x01 {
		connection.Write([]byte{0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}

	var host string

	if request[3] == 0x01 {

		address := make([]byte, 4)

		if _, err := io.ReadFull(connection, address); err != nil {
			return
		}

		host = net.IP(address).String()

	} else if request[3] == 0x04 {

		address := make([]byte, 16)

		if _, err := io.ReadFull(connection, address); err != nil {
			return
		}

		host = net.IP(address).String()

	} else if request[3] == 0x03 {

		length := make([]byte, 1)

		if _, err := io.ReadFull(connection, length); err != nil {
			return
		}

		address := make([]byte, int(length[0]))

		if _, err := io.ReadFull(connection, address); err != nil {
			return
		}

		host = string(address)

	} else {
		connection.Write([]byte{0x05, 0x08, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}

	port_bytes := make([]byte, 2)

	if _, err := io.ReadFull(connection, port_bytes); err != nil {
		return
	}

	port := int(port_bytes[0])<<8 | int(port_bytes[1])

	target, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))

	if err != nil {
		connection.Write([]byte{0x05, 0x05, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}

	defer target.Close()

	if _, err := connection.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0}); err != nil {
		return
	}

	go io.Copy(target, connection)
	io.Copy(connection, target)

}
