package tunnel

import "tholian-warps/protocols/http"
import "tholian-warps/protocols/icmp"

func DecodeRequest(message *icmp.Message) http.Packet {
	return http.Parse(message.Payload)
}
