package http

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-endpoint/types"
import http_tunnel "tholian-warps/protocols/http/tunnel"
import net_url "net/url"
import "strconv"

type Tunnel struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func NewTunnel(host string, port uint16) Tunnel {

	var tunnel Tunnel

	tunnel.Host = host
	tunnel.Port = port

	return tunnel

}

func (tunnel *Tunnel) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	url, err := net_url.Parse("http://" + tunnel.Host + ":" + strconv.FormatUint(uint64(tunnel.Port), 10) + "/dns-query")

	if err == nil {

		tunnel_request := http.NewPacket()
		tunnel_request.SetURL(*url)

		http_tunnel.EncodePayload(&tunnel_request, query.Bytes())

		tunnel_response := http.RequestPacket(tunnel_request)

		if http_tunnel.IsResolveResponse(&tunnel_response) {

			response = dns.Parse(http_tunnel.DecodePayload(&tunnel_response))

		} else {

			response = dns.NewPacket()
			response.SetType("response")
			response.SetIdentifier(query.Identifier)
			response.SetResponseCode(dns.ResponseCodeNonExistDomain)

			for q := 0; q < len(query.Questions); q++ {
				response.AddQuestion(query.Questions[q])
			}

		}

	} else {

		response = dns.NewPacket()
		response.SetType("response")
		response.SetIdentifier(query.Identifier)
		response.SetResponseCode(dns.ResponseCodeNonExistDomain)

		for q := 0; q < len(query.Questions); q++ {
			response.AddQuestion(query.Questions[q])
		}

	}

	return response

}

func (tunnel *Tunnel) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	tunnel_request := request
	tunnel_request.SetServer(types.Server{
		Domain:    tunnel.Host,
		Addresses: []string{tunnel.Host},
		Port:      tunnel.Port,
		Protocol:  types.ProtocolHTTP,
		Schema:    "DEFAULT",
	})

	tunnel_response := http.RequestPacket(tunnel_request)

	if tunnel_response.Type == "response" {

		response = tunnel_response
		response.Decode()

	} else {

		response = http.NewPacket()
		response.SetURL(*request.URL)
		response.SetStatus(http.StatusNotFound)
		response.SetPayload([]byte{})

	}

	return response

}
