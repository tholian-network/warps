package tunnel

import "tholian-warps/protocols/icmp"

func EncodeError(message *icmp.Message) {

	message.SetType(icmp.MessageKick)
	message.SetTarget("")
	message.SetPayload([]byte{})

}
