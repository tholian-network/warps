package icmp

func Parse(buffer []byte) Packet {

	packet := NewPacket()

	if len(buffer) >= 8 {

		if buffer[0] == 8 {
			packet.SetType("request")
		} else {
			packet.SetType("reply")
		}

		packet.SetCode(buffer[1])
		packet.SetIdentifier(uint16(buffer[4])<<8 | uint16(buffer[5]))
		packet.SetSequence(uint16(buffer[6])<<8 | uint16(buffer[7]))

		if len(buffer) > 8 {
			packet.SetPayload(buffer[8:])
		}

	}

	return packet

}
