package dns

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-endpoint/types"
import utils_url "tholian-warps/utils/net/url"
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

	// TODO: Relay DNS Packets via tunnel

	return response

}

func (tunnel *Tunnel) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	domain := utils_url.ToHost(request.URL.String())

	// Range: bytes=0-
	range_domain := "bytes.0-." + domain
	range_record := dns.NewRecord(range_domain, dns.TypeURI)
	range_record.SetURL(request.URL.String())

	tunnel_request := dns.NewPacket()
	tunnel_request.SetType("query")
	tunnel_request.AddQuestion(dns.NewQuestion(domain, dns.TypeA)) // Obfuscation Question
	tunnel_request.AddQuestion(dns.NewQuestion(range_domain, dns.TypeURI))
	tunnel_request.AddAnswer(range_record)
	tunnel_request.SetServer(types.Server{
		Domain:    tunnel.Host,
		Addresses: []string{tunnel.Host},
		Port:      tunnel.Port,
		Protocol:  types.ProtocolDNS,
		Schema:    "DEFAULT",
	})

	first_response := dns.ResolvePacket(tunnel_request)

	if first_response.Type == "response" {

		// TODO: Content-Type and other headers necessary?
		// headers := toDNSHeaders(first_response)

		payload := make([]byte, 0)
		payload = append(payload, toDNSPayload(first_response)...)

		payload_from, payload_to, payload_size := toDNSContentRange(first_response)

		if payload_from == 0 && payload_to != 0 && payload_size > 0 {

			// Some network routes only support 1232 bytes DNS packet size
			// The default payload frame size is 1024 bytes, but in case
			// the network route supports bigger frame sizes
			frame_size := payload_to + 1

			for len(payload) < payload_size {

				frame_from := len(payload)
				frame_to := len(payload) + frame_size

				// Range Request
				frame_range_domain := "bytes." + strconv.Itoa(frame_from) + "-" + strconv.Itoa(frame_to) + "." + domain
				frame_range_record := dns.NewRecord(frame_range_domain, dns.TypeURI)
				frame_range_record.SetURL(request.URL.String())

				frame_request := dns.NewPacket()
				frame_request.SetType("query")
				frame_request.AddQuestion(dns.NewQuestion(domain, dns.TypeA)) // Obfuscation Question
				frame_request.AddQuestion(dns.NewQuestion(frame_range_domain, dns.TypeURI))
				frame_request.AddAnswer(frame_range_record)
				frame_request.SetServer(types.Server{
					Domain:    tunnel.Host,
					Addresses: []string{tunnel.Host},
					Port:      tunnel.Port,
					Protocol:  types.ProtocolDNS,
					Schema:    "DEFAULT",
				})

				frame_response := dns.ResolvePacket(frame_request)
				frame_response_from, frame_response_to, frame_response_size := toDNSContentRange(frame_response)

				if frame_from == frame_response_from && frame_to == frame_response_to && payload_size == frame_response_size {
					payload = append(payload, toDNSPayload(frame_response)...)
				} else {
					// TODO: Respond with Content-Range error
					break
				}

			}

			// TODO: Are more headers necessary?
			// TODO: Content-Type?
			response := http.NewPacket()
			response.SetURL(*request.URL)
			response.SetPayload(payload)

		} else {
			// TODO: Respond with Content-Range error
		}

	}

	return response

}
