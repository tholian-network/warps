package http

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"
import "tholian-endpoint/types"
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

		payload := query.Bytes()

		tunnel_request := http.NewPacket()
		tunnel_request.SetMethod(http.MethodPost)
		tunnel_request.SetURL(*url)
		tunnel_request.SetHeader("Accept", "application/dns-message")
		tunnel_request.SetHeader("Content-Type", "application/dns-message")
		tunnel_request.SetHeader("Content-Length", strconv.Itoa(len(payload)))
		tunnel_request.SetPayload(payload)

		tunnel_response := http.RequestPacket(tunnel_request)

		if tunnel_response.Type == "response" && tunnel_response.Status == http.StatusOK && tunnel_response.GetHeader("Content-Type") == "application/dns-message" {

			tunnel_response.Decode()

			tmp := dns.Parse(tunnel_response.Payload)

			if tmp.Type == "response" {
				response = tmp
			} else {
				response = EmptyDNSResponse(query)
			}

		} else {
			response = EmptyDNSResponse(query)
		}

	} else {
		response = EmptyDNSResponse(query)
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
		response = EmptyHTTPResponse(request)
	}

	return response

}
