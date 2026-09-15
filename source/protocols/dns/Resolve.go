package dns

import "tholian-warps/types"

func Resolve(subject string) (Packet, error) {

	var response Packet
	var err error = nil

	if types.IsIPv4(subject) {

		request := NewPacket()
		request.SetType("query")
		request.AddQuestion(NewQuestion(subject, TypePTR))

		response, err = ResolvePacket(request)

	} else if types.IsIPv6(subject) {

		request := NewPacket()
		request.SetType("query")
		request.AddQuestion(NewQuestion(subject, TypePTR))

		response, err = ResolvePacket(request)

	} else if types.IsDomain(subject) {

		request := NewPacket()
		request.SetType("query")
		request.AddQuestion(NewQuestion(subject, TypeA))
		request.AddQuestion(NewQuestion(subject, TypeAAAA))

		response, err = ResolvePacket(request)

	}

	return response, err

}
