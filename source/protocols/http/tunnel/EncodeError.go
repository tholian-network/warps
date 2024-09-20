package tunnel

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/protocols/http"

func EncodeError(query *dns.Packet, response *http.Packet, status http.Status) {

	dns_response := dns.NewPacket()
	dns_response.SetType("response")
	dns_response.SetIdentifier(query.Identifier)
	dns_response.SetResponseCode(dns.ResponseCodeNonExistDomain)

	for q := 0; q < len(query.Questions); q++ {
		dns_response.AddQuestion(query.Questions[q])
	}

	response.SetStatus(status)
	response.SetHeader("Content-Type", "application/dns-message")
	response.SetPayload(dns_response.Bytes())

}
