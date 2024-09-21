package dns

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-endpoint/types"
import dns_tunnel "tholian-warps/protocols/dns/tunnel"
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

	tunnel_query := query
	tunnel_query.SetServer(types.Server{
		Domain:    tunnel.Host,
		Addresses: []string{tunnel.Host},
		Port:      tunnel.Port,
		Protocol:  types.ProtocolDNS,
		Schema:    "DEFAULT",
	})

	tunnel_response := dns.ResolvePacket(tunnel_query)

	if tunnel_response.Type == "response" {

		response = tunnel_response

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

	tunnel_request := dns.NewPacket()
	tunnel_request.SetType("query")

	// Range: bytes=0-
	dns_tunnel.EncodeContentRange(&tunnel_request, request.URL, 0, -1, -1)

	tunnel_request.SetServer(types.Server{
		Domain:    tunnel.Host,
		Addresses: []string{tunnel.Host},
		Port:      tunnel.Port,
		Protocol:  types.ProtocolDNS,
		Schema:    "DEFAULT",
	})

	first_response := dns.ResolvePacket(tunnel_request)

	if first_response.Type == "response" {

		if first_response.Codes.Response == dns.ResponseCodeNoError {

			headers := dns_tunnel.DecodeHeaders(&first_response)
			payload := make([]byte, 0)
			payload = append(payload, dns_tunnel.DecodePayload(&first_response)...)

			_, payload_from, payload_to, payload_size := dns_tunnel.DecodeContentRange(&first_response)

			if payload_from == 0 && payload_to != 0 && payload_size > 1024 {

				// TODO: HTTP response SetHeader("Content-Range", "bytes 0-123/1234") header

				// Some network routes only support 1232 bytes DNS packet size
				// The default payload frame size is 1024 bytes, but in case
				// the network route supports bigger frame sizes
				frame_size := payload_to + 1
				range_error := false

				for len(payload) < payload_size {

					frame_from := len(payload)
					frame_to := len(payload) + frame_size

					frame_request := dns.NewPacket()
					frame_request.SetType("query")

					// Range: bytes=<from>-<to>
					dns_tunnel.EncodeContentRange(&tunnel_request, request.URL, frame_from, frame_to, -1)

					frame_request.SetServer(types.Server{
						Domain:    tunnel.Host,
						Addresses: []string{tunnel.Host},
						Port:      tunnel.Port,
						Protocol:  types.ProtocolDNS,
						Schema:    "DEFAULT",
					})

					frame_response := dns.ResolvePacket(frame_request)
					_, frame_response_from, frame_response_to, frame_response_size := dns_tunnel.DecodeContentRange(&frame_response)

					if frame_from == frame_response_from && frame_to == frame_response_to && payload_size == frame_response_size {
						payload = append(payload, dns_tunnel.DecodePayload(&frame_response)...)
					} else {
						range_error = true
						break
					}

				}

				if range_error == false {

					response = http.NewPacket()
					response.SetURL(*request.URL)
					response.SetStatus(http.StatusOK)

					for key, val := range headers {
						response.SetHeader(key, val)
					}

					response.SetHeader("Content-Range", "bytes " + strconv.Itoa(0) + "-" + strconv.Itoa(len(payload) - 1) + "/" + strconv.Itoa(len(payload)))
					response.SetPayload(payload)

				} else {

					response = http.NewPacket()
					response.SetURL(*request.URL)
					response.SetStatus(http.StatusRangeNotSatisfiable)
					response.SetPayload([]byte{})

				}

			} else if payload_from == 0 && payload_to != 0 && payload_size > 0 {

				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusOK)

				for key, val := range headers {
					response.SetHeader(key, val)
				}

				response.SetHeader("Content-Range", "bytes " + strconv.Itoa(payload_from) + "-" + strconv.Itoa(payload_to) + "/" + strconv.Itoa(payload_size))
				response.SetPayload(payload)

			} else {

				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusRangeNotSatisfiable)
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
