package test

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"

type SpyTunnel struct {
	Resolved     string `json:"resolved"`
	Requested    string `json:"requested"`
	WasResolved  bool   `json:"was_resolved"`
	WasRequested bool   `json:"was_requested"`
	isOnline     bool
}

func NewSpyTunnel(isOnline bool) SpyTunnel {

	var tunnel SpyTunnel

	tunnel.isOnline = isOnline

	return tunnel

}

func (tunnel *SpyTunnel) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	if query.Type == "query" && len(query.Questions) > 0 {

		resolved := ""

		for q := 0; q < len(query.Questions); q++ {

			resolved += query.Questions[q].Type.String() + ":" + query.Questions[q].Name

			if q <= len(query.Questions) - 1 {
				resolved += ","
			}

		}

		if tunnel.isOnline == true {
			response = dns.ResolvePacket(query)
		}

		tunnel.WasResolved = true
		tunnel.Resolved = resolved

	}

	return response

}

func (tunnel *SpyTunnel) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	if request.Type == "request" && request.URL != nil {

		requested := request.URL.String()

		if tunnel.isOnline == true {
			response = http.RequestPacket(request)
		}

		tunnel.WasRequested = true
		tunnel.Requested = requested

	}

	return response

}

