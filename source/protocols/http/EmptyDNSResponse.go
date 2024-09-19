package http

import "tholian-endpoint/protocols/dns"

func EmptyDNSResponse(query dns.Packet) dns.Packet {

	var response dns.Packet

	response.SetType("response")
	response.SetIdentifier(query.Identifier)
	response.Codes.Response = dns.ResponseCodeNonExistDomain

	for q := 0; q < len(query.Questions); q++ {

		question := query.Questions[q]

		if question.Type == dns.TypeA {
			response.AddQuestion(question)
			response.AddAnswer(dns.NewRecord(question.Name, dns.TypeA))
		} else if question.Type == dns.TypeAAAA {
			response.AddQuestion(question)
			response.AddAnswer(dns.NewRecord(question.Name, dns.TypeAAAA))
		}

	}

	return response

}
