package http

import "tholian-warps/types"

func RequestPacket(request Packet) Packet {

	var response Packet

	if request.Server == nil {
		request.Resolve()
	}

	if request.Server != nil {

		var tmp Packet

		if request.Server.Protocol == types.ProtocolHTTP {
			tmp = requestTCP(request.Server.RandomizeAddress(), request.Server.Port, request)
		} else if request.Server.Protocol == types.ProtocolHTTPS {
			tmp = requestTLS(request.Server.Domain, request.Server.RandomizeAddress(), request.Server.Port, request)
		}

		if tmp.Type == "response" {
			response = tmp
		}

	}

	return response

}
