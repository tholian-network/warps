package icmptunnel

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"
import "tholian-warps/console"
import icmp_tunnel "tholian-warps/protocols/icmptunnel/tunnel"
import "math/rand"
import "net"
import "strconv"
import "sync"
import "time"

type Tunnel struct {
	Host       string     `json:"host"`
	Port       uint16     `json:"port"`
	Key        int32      `json:"key"`
	Identifier uint16     `json:"identifier"`
	socket     icmp.Socket
	sequence   uint16
	debug      bool
	mutex      *sync.Mutex
}

func NewTunnel(host string, port uint16) Tunnel {

	var tunnel Tunnel

	tunnel.Host = host
	tunnel.Port = port
	tunnel.Key = 0
	tunnel.Identifier = uint16(rand.Uint64())
	tunnel.socket = nil
	tunnel.sequence = 0
	tunnel.debug = false
	tunnel.mutex = &sync.Mutex{}

	return tunnel

}

func (tunnel *Tunnel) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	tunnel.mutex.Lock()
	defer tunnel.mutex.Unlock()

	identifier := strconv.FormatUint(uint64(rand.Uint64()), 16)

	message := icmp.NewMessage()
	icmp_tunnel.EncodeResolveRequest(&message, query)
	message.SetIdentifier(identifier)
	message.SetKey(tunnel.Key)

	tunnel.sendMessage(&message)

	reply := tunnel.readReply(identifier)

	if reply.Type == icmp.MessageData && icmp_tunnel.IsResolveResponse(&reply) {

		response = icmp_tunnel.DecodeResolveResponse(&reply)

	} else {

		response = dns.NewPacket()
		response.SetType("response")
		response.SetIdentifier(query.Identifier)
		response.SetResponseCode(dns.ResponseCodeNonExistDomain)
		response.Flags.RecursionAvailable = true

		for q := 0; q < len(query.Questions); q++ {
			response.AddQuestion(query.Questions[q])
		}

	}

	return response

}

func (tunnel *Tunnel) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	tunnel.mutex.Lock()
	defer tunnel.mutex.Unlock()

	if tunnel.debug {
		console.Group("icmp/Tunnel/RequestPacket")
	}

	identifier := strconv.FormatUint(uint64(rand.Uint64()), 16)

	message := icmp.NewMessage()
	icmp_tunnel.EncodeRequest(&message, request)
	message.SetIdentifier(identifier)
	message.SetKey(tunnel.Key)

	tunnel.sendMessage(&message)

	reply := tunnel.readReply(identifier)

	if reply.Type == icmp.MessageData && icmp_tunnel.IsResponse(&reply) {

		response = icmp_tunnel.DecodeResponse(&reply)
		response.Decode()

	} else {

		response = http.NewPacket()
		response.SetURL(*request.URL)
		response.SetStatus(http.StatusNotFound)
		response.SetPayload([]byte{})

	}

	if tunnel.debug {
		console.GroupEnd("icmp/Tunnel/RequestPacket")
	}

	return response

}

func (tunnel *Tunnel) SetDebug(value bool) {
	tunnel.debug = value
}

func (tunnel *Tunnel) SetKey(value int32) {
	tunnel.Key = value
}

func (tunnel *Tunnel) SetSocket(value icmp.Socket) {
	tunnel.socket = value
}

func (tunnel *Tunnel) ensureSocket() icmp.Socket {

	if tunnel.socket == nil {

		ping := icmp.NewPingSocket()
		tunnel.socket = &ping
		tunnel.socket.Listen("0.0.0.0")

	}

	return tunnel.socket

}

func (tunnel *Tunnel) resolveServer() *net.IPAddr {

	var result *net.IPAddr = nil

	ip := net.ParseIP(tunnel.Host)

	if ip != nil && ip.To4() != nil {
		result = &net.IPAddr{IP: ip.To4()}
	} else {

		addr, err := net.ResolveIPAddr("ip", tunnel.Host)

		if err == nil {
			result = addr
		}

	}

	return result

}

func (tunnel *Tunnel) sendMessage(message *icmp.Message) {

	socket := tunnel.ensureSocket()
	address := tunnel.resolveServer()

	if address == nil {
		return
	}

	packet := icmp.NewPacket()
	packet.SetType("request")
	packet.SetIdentifier(tunnel.Identifier)

	tunnel.sequence++
	packet.SetSequence(tunnel.sequence)

	packet.SetPayload(message.Bytes())

	socket.WritePacket(packet, address)

}

func (tunnel *Tunnel) readReply(identifier string) icmp.Message {

	var result icmp.Message

	socket := tunnel.ensureSocket()

	channel := make(chan icmp.Message, 1)

	go func() {

		for {

			packet, _, err := socket.ReadPacket()

			if err != nil {
				channel <- icmp.NewMessage()
				return
			}

			message := icmp.ParseMessage(packet.Payload)

			if message.Identifier == identifier {
				channel <- message
				return
			}

		}

	}()

	select {
	case message := <-channel:
		result = message
	case <-time.After(3 * time.Second):
		result = icmp.NewMessage()
	}

	return result

}
