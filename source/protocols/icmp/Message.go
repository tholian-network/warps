package icmp

type MessageType int32

const (
	MessageData  MessageType = 0
	MessagePing  MessageType = 1
	MessageKick  MessageType = 2
	MessageMagic MessageType = 57005
)

type Message struct {
	Identifier string      `json:"id"`
	Type       MessageType `json:"type"`
	Target     string      `json:"target"`
	Payload    []byte      `json:"payload"`
	Protocol   int32       `json:"protocol"`
	Magic      int32       `json:"magic"`
	Key        int32       `json:"key"`
	Timeout    int32       `json:"timeout"`
}

func NewMessage() Message {

	var message Message

	message.Identifier = ""
	message.Type = MessageData
	message.Target = ""
	message.Payload = make([]byte, 0)
	message.Protocol = 0
	message.Magic = int32(MessageMagic)
	message.Key = 0
	message.Timeout = 0

	return message

}

func (message *Message) SetIdentifier(value string) {
	message.Identifier = value
}

func (message *Message) SetType(value MessageType) {
	message.Type = value
}

func (message *Message) SetTarget(value string) {
	message.Target = value
}

func (message *Message) SetPayload(value []byte) {
	message.Payload = value
}

func (message *Message) SetProtocol(value int32) {
	message.Protocol = value
}

func (message *Message) SetKey(value int32) {
	message.Key = value
}

func (message *Message) SetTimeout(value int32) {
	message.Timeout = value
}

func (message *Message) IsValid() bool {
	return message.Magic == int32(MessageMagic)
}

func (message *Message) Bytes() []byte {

	var buffer []byte

	if message.Identifier != "" {
		buffer = appendTag(buffer, 1, 2)
		buffer = appendVarint(buffer, uint64(len(message.Identifier)))
		buffer = append(buffer, message.Identifier...)
	}

	if message.Type != 0 {
		buffer = appendTag(buffer, 2, 0)
		buffer = appendVarint(buffer, uint64(message.Type))
	}

	if message.Target != "" {
		buffer = appendTag(buffer, 3, 2)
		buffer = appendVarint(buffer, uint64(len(message.Target)))
		buffer = append(buffer, message.Target...)
	}

	if len(message.Payload) > 0 {
		buffer = appendTag(buffer, 4, 2)
		buffer = appendVarint(buffer, uint64(len(message.Payload)))
		buffer = append(buffer, message.Payload...)
	}

	if message.Protocol != 0 {
		buffer = appendTag(buffer, 5, 0)
		buffer = appendVarint(buffer, uint64(toZigzag32(message.Protocol)))
	}

	if message.Magic != 0 {
		buffer = appendTag(buffer, 6, 0)
		buffer = appendVarint(buffer, uint64(toZigzag32(message.Magic)))
	}

	if message.Key != 0 {
		buffer = appendTag(buffer, 7, 0)
		buffer = appendVarint(buffer, uint64(toZigzag32(message.Key)))
	}

	if message.Timeout != 0 {
		buffer = appendTag(buffer, 8, 0)
		buffer = appendVarint(buffer, uint64(message.Timeout))
	}

	return buffer

}

func ParseMessage(buffer []byte) Message {

	var message Message
	message.Payload = make([]byte, 0)

	for len(buffer) > 0 {

		tag, tag_length := readVarint(buffer)

		if tag_length <= 0 {
			break
		}

		buffer = buffer[tag_length:]

		field := int(tag >> 3)
		wire := int(tag & 7)

		if wire == 0 {

			value, value_length := readVarint(buffer)

			if value_length <= 0 {
				break
			}

			buffer = buffer[value_length:]

			if field == 2 {
				message.Type = MessageType(int32(value))
			} else if field == 5 {
				message.Protocol = fromZigzag32(uint32(value))
			} else if field == 6 {
				message.Magic = fromZigzag32(uint32(value))
			} else if field == 7 {
				message.Key = fromZigzag32(uint32(value))
			} else if field == 8 {
				message.Timeout = int32(value)
			}

		} else if wire == 2 {

			length, length_length := readVarint(buffer)

			if length_length <= 0 {
				break
			}

			buffer = buffer[length_length:]

			if uint64(len(buffer)) < length {
				break
			}

			value := buffer[0:int(length)]
			buffer = buffer[int(length):]

			if field == 1 {
				message.Identifier = string(value)
			} else if field == 3 {
				message.Target = string(value)
			} else if field == 4 {
				message.Payload = value
			}

		} else {
			break
		}

	}

	return message

}

func appendTag(buffer []byte, field int, wire int) []byte {
	return appendVarint(buffer, uint64(field<<3|wire))
}

func appendVarint(buffer []byte, value uint64) []byte {

	for value >= 0x80 {
		buffer = append(buffer, byte(value)|0x80)
		value >>= 7
	}

	return append(buffer, byte(value))

}

func readVarint(buffer []byte) (uint64, int) {

	var value uint64

	for i := 0; i < len(buffer) && i < 10; i++ {

		value |= uint64(buffer[i]&0x7f) << (7 * i)

		if buffer[i]&0x80 == 0 {
			return value, i + 1
		}

	}

	return 0, 0

}

func toZigzag32(value int32) uint32 {
	return uint32(value<<1) ^ uint32(value>>31)
}

func fromZigzag32(value uint32) int32 {
	return int32(value>>1) ^ -int32(value&1)
}
