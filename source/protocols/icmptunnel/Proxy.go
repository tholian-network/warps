package icmptunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"
import "tholian-warps/interfaces"
import icmp_tunnel "tholian-warps/protocols/icmptunnel/tunnel"
import "net"
import "strings"

type Proxy struct {
	Host     string                `json:"host"`
	Port     uint16                `json:"port"`
	Key      int32                 `json:"key"`
	Cache    interfaces.ProxyCache `json:"cache"`
	Tunnel   interfaces.Tunnel     `json:"tunnel"`
	Resolver interfaces.Resolver   `json:"resolver"`
	socket   icmp.Socket
	closed   bool
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
	proxy.Key = 0
	proxy.Cache = cache
	proxy.socket = nil
	proxy.closed = false

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

func (proxy *Proxy) SetKey(value int32) {
	proxy.Key = value
}

func (proxy *Proxy) SetSocket(value icmp.Socket) {
	proxy.socket = value
}

func (proxy *Proxy) Destroy() error {

	var err error = nil

	proxy.closed = true

	if proxy.socket != nil {
		err = proxy.socket.Close()
	}

	return err

}

func (proxy *Proxy) Listen() error {

	var err error = nil

	socket := proxy.socket

	if socket == nil {

		ping := icmp.NewPingSocket()
		socket = &ping

	}

	proxy.socket = socket

	err = socket.Listen(proxy.Host)

	if err == nil {

		for {

			if proxy.closed || proxy.socket == nil {
				break
			}

			packet, remote, err2 := socket.ReadPacket()

			if err2 != nil {
				break
			}

			message := icmp.ParseMessage(packet.Payload)

			if message.IsValid() == false {
				continue
			}

			if message.Key != proxy.Key {
				continue
			}

			go func(remote net.Addr, packet icmp.Packet, message icmp.Message) {
				proxy.processPacket(remote, packet, message)
			}(remote, packet, message)

		}

	}

	return err

}

func (proxy *Proxy) processPacket(remote net.Addr, packet icmp.Packet, message icmp.Message) {

	reply := icmp.NewMessage()
	icmp_tunnel.EncodeError(&reply)
	reply.SetIdentifier(message.Identifier)
	reply.SetKey(proxy.Key)

	if message.Type == icmp.MessagePing {

		reply = icmp.NewMessage()
		reply.SetType(icmp.MessagePing)
		reply.SetIdentifier(message.Identifier)
		reply.SetPayload(message.Payload)
		reply.SetKey(proxy.Key)

	} else if icmp_tunnel.IsResolveRequest(&message) {

		query := icmp_tunnel.DecodeResolveRequest(&message)
		response := proxy.ResolvePacket(query)

		if response.Type == "response" {

			reply = icmp.NewMessage()
			icmp_tunnel.EncodeResolveResponse(&reply, response)
			reply.SetIdentifier(message.Identifier)
			reply.SetKey(proxy.Key)

		}

	} else if icmp_tunnel.IsRequest(&message) {

		request := icmp_tunnel.DecodeRequest(&message)
		response := proxy.RequestPacket(request)

		if response.Type == "response" {

			reply = icmp.NewMessage()
			icmp_tunnel.EncodeResponse(&reply, response)
			reply.SetIdentifier(message.Identifier)
			reply.SetKey(proxy.Key)

		}

	}

	proxy.writeMessage(remote, packet, reply)

}

func (proxy *Proxy) writeMessage(remote net.Addr, packet icmp.Packet, message icmp.Message) {

	if proxy.socket == nil {
		return
	}

	reply := icmp.NewPacket()
	reply.SetType("reply")
	reply.SetIdentifier(packet.Identifier)
	reply.SetSequence(packet.Sequence)
	reply.SetPayload(message.Bytes())

	proxy.socket.WritePacket(reply, remote)

}
