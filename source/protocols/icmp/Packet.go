package icmp

type Packet struct {
	Type       string `json:"type"`
	Code       uint8  `json:"code"`
	Identifier uint16 `json:"identifier"`
	Sequence   uint16 `json:"sequence"`
	Payload    []byte `json:"payload"`
}

func NewPacket() Packet {

	var packet Packet

	packet.Type = "request"
	packet.Code = 0
	packet.Identifier = 0
	packet.Sequence = 0
	packet.Payload = make([]byte, 0)

	return packet

}

func (packet *Packet) SetType(value string) {

	if value == "request" {
		packet.Type = "request"
	} else if value == "reply" {
		packet.Type = "reply"
	}

}

func (packet *Packet) SetCode(value uint8) {
	packet.Code = value
}

func (packet *Packet) SetIdentifier(value uint16) {
	packet.Identifier = value
}

func (packet *Packet) SetSequence(value uint16) {
	packet.Sequence = value
}

func (packet *Packet) SetPayload(value []byte) {
	packet.Payload = value
}

func (packet *Packet) Bytes() []byte {

	bytes := make([]byte, 8)

	if packet.Type == "request" {
		bytes[0] = 8
	} else if packet.Type == "reply" {
		bytes[0] = 0
	}

	bytes[1] = packet.Code

	bytes[4] = byte(packet.Identifier >> 8)
	bytes[5] = byte(packet.Identifier & 0xff)
	bytes[6] = byte(packet.Sequence >> 8)
	bytes[7] = byte(packet.Sequence & 0xff)

	if len(packet.Payload) > 0 {
		bytes = append(bytes, packet.Payload...)
	}

	checksum := toChecksum(bytes)

	bytes[2] = byte(checksum >> 8)
	bytes[3] = byte(checksum & 0xff)

	return bytes

}

func toChecksum(data []byte) uint16 {

	var sum uint32

	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}

	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}

	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	return ^uint16(sum)

}
