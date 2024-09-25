package dns

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-endpoint/types"
import "tholian-warps/console"
import dns_tunnel "tholian-warps/protocols/dns/tunnel"
import "strconv"

type Tunnel struct {
	Host  string `json:"host"`
	Port  uint16 `json:"port"`
	debug bool
}

func NewTunnel(host string, port uint16) Tunnel {

	var tunnel Tunnel

	tunnel.Host = host
	tunnel.Port = port
	tunnel.debug = false

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

	tunnel_response, err := dns.ResolvePacket(tunnel_query)

	if err == nil && tunnel_response.Type == "response" {

		response = tunnel_response

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

	if tunnel.debug {
		console.Group("dns/Tunnel/RequestPacket")
	}

	tunnel_request := dns.NewPacket()
	tunnel_request.SetType("query")

	// Domains: 0-.bytes.domain.tld and 0-.headers.domain.tld
	dns_tunnel.EncodeFirstRequest(&tunnel_request, request.URL)

	tunnel_request.SetServer(types.Server{
		Domain:    tunnel.Host,
		Addresses: []string{tunnel.Host},
		Port:      tunnel.Port,
		Protocol:  types.ProtocolDNS,
		Schema:    "DEFAULT",
	})

	if tunnel.debug {
		console.Warn("Request:")
		console.Inspect(tunnel_request.Server)
	}

	first_response, first_response_err := dns.ResolvePacket(tunnel_request)

	if first_response_err == nil && first_response.Type == "response" {

		if first_response.Codes.Response == dns.ResponseCodeNoError {

			headers := dns_tunnel.DecodeHeaders(&first_response)
			payload := make([]byte, 0)
			payload = append(payload, dns_tunnel.DecodePayload(&first_response)...)

			payload_from, payload_to, payload_size := dns_tunnel.DecodeContentRange(&first_response)

			if tunnel.debug {
				console.Log("First Response Content-Range: " + strconv.Itoa(payload_from) + "-" + strconv.Itoa(payload_to) + "/" + strconv.Itoa(payload_size))
			}

			if payload_from == 0 && payload_to > 0 && payload_size > 512 {

				// Some network routes only support 1232 bytes DNS packet size
				// The default payload frame size is 1024 bytes, but in case
				// the network route supports bigger frame sizes
				frame_size := payload_to + 1
				range_error := false

				for len(payload) < payload_size {

					frame_from := len(payload)
					frame_to := len(payload) + frame_size

					if frame_to > payload_size - 1 {
						frame_to = payload_size - 1
					}

					frame_request := dns.NewPacket()
					frame_request.SetType("query")

					// Range: bytes=<from>-<to>
					dns_tunnel.EncodeFrameRequest(&frame_request, request.URL, frame_from, frame_to)

					frame_request.SetServer(types.Server{
						Domain:    tunnel.Host,
						Addresses: []string{tunnel.Host},
						Port:      tunnel.Port,
						Protocol:  types.ProtocolDNS,
						Schema:    "DEFAULT",
					})

					if tunnel.debug {
						console.Log("Frame Request Range: " + strconv.Itoa(frame_from) + "-" + strconv.Itoa(frame_to))
					}

					frame_response, frame_response_err := dns.ResolvePacket(frame_request)

					if frame_response_err == nil && frame_response.Type == "response" {

						frame_response_from, frame_response_to, frame_response_size := dns_tunnel.DecodeContentRange(&frame_response)

						if tunnel.debug {
							console.Log("Frame Response Content-Range: " + strconv.Itoa(frame_response_from) + "-" + strconv.Itoa(frame_response_to) + "/" + strconv.Itoa(frame_response_size))
						}

						if frame_from == frame_response_from && frame_to == frame_response_to && payload_size == frame_response_size {
							payload = append(payload, dns_tunnel.DecodePayload(&frame_response)...)
						} else {
							range_error = true
							break
						}

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

					if tunnel.debug {
						console.Info("Response Content-Range: 0-" + strconv.Itoa(len(payload) - 1) + "/" + strconv.Itoa(len(payload)))
						console.Inspect(response)
					}

					response.SetPayload(payload)

				} else {

					response = http.NewPacket()
					response.SetURL(*request.URL)
					response.SetStatus(http.StatusRangeNotSatisfiable)
					response.SetPayload([]byte{})

				}

			} else if payload_from == 0 && payload_to > 0 && payload_size > 0 {

				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusOK)

				for key, val := range headers {
					response.SetHeader(key, val)
				}

				response.SetHeader("Content-Range", "bytes " + strconv.Itoa(payload_from) + "-" + strconv.Itoa(payload_to) + "/" + strconv.Itoa(payload_size))
				response.SetPayload(payload)

				if tunnel.debug {
					console.Info("Response Content-Range: 0-" + strconv.Itoa(payload_to) + "/" + strconv.Itoa(payload_size))
					console.Inspect(response)
				}

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

	if tunnel.debug {
		console.GroupEnd("dns/Tunnel/RequestPacket")
	}

	return response

}

func (tunnel *Tunnel) SetDebug(value bool) {
	tunnel.debug = value
}
