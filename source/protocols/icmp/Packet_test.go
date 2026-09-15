package icmp

import "testing"

func TestPacket(t *testing.T) {

	t.Run("Bytes with request", func(t *testing.T) {

		packet := NewPacket()
		packet.SetType("request")
		packet.SetIdentifier(1337)
		packet.SetSequence(42)
		packet.SetPayload([]byte("Hello, world!"))

		bytes := packet.Bytes()

		if len(bytes) != 8+len(packet.Payload) {
			t.Errorf("Expected ICMP packet to have '%d' bytes but got '%d'", 8+len(packet.Payload), len(bytes))
		}

		if bytes[0] != 8 {
			t.Errorf("Expected ICMP type to be '%d' but got '%d'", 8, bytes[0])
		}

		if bytes[4] != 0x05 || bytes[5] != 0x39 {
			t.Errorf("Expected ICMP identifier to be '0x0539' but got '%02x%02x'", bytes[4], bytes[5])
		}

		if bytes[6] != 0x00 || bytes[7] != 0x2a {
			t.Errorf("Expected ICMP sequence to be '0x002a' but got '%02x%02x'", bytes[6], bytes[7])
		}

	})

	t.Run("Bytes with reply", func(t *testing.T) {

		packet := NewPacket()
		packet.SetType("reply")
		packet.SetIdentifier(1337)
		packet.SetSequence(42)
		packet.SetPayload([]byte("Hello, world!"))

		bytes := packet.Bytes()

		if bytes[0] != 0 {
			t.Errorf("Expected ICMP type to be '%d' but got '%d'", 0, bytes[0])
		}

	})

	t.Run("toChecksum", func(t *testing.T) {

		data := []byte{8, 0, 0, 0, 0x05, 0x39, 0x00, 0x2a}
		checksum := toChecksum(data)

		if checksum == 0 {
			t.Errorf("Expected ICMP checksum to be non-zero")
		}

	})

}

func TestParse(t *testing.T) {

	t.Run("Parse request", func(t *testing.T) {

		packet := NewPacket()
		packet.SetType("request")
		packet.SetIdentifier(1337)
		packet.SetSequence(42)
		packet.SetPayload([]byte("Hello, world!"))

		parsed := Parse(packet.Bytes())

		if parsed.Type != "request" {
			t.Errorf("Expected ICMP type to be '%s' but got '%s'", "request", parsed.Type)
		}

		if parsed.Identifier != 1337 {
			t.Errorf("Expected ICMP identifier to be '%d' but got '%d'", 1337, parsed.Identifier)
		}

		if parsed.Sequence != 42 {
			t.Errorf("Expected ICMP sequence to be '%d' but got '%d'", 42, parsed.Sequence)
		}

		if string(parsed.Payload) != "Hello, world!" {
			t.Errorf("Expected ICMP payload to be '%s' but got '%s'", "Hello, world!", string(parsed.Payload))
		}

	})

	t.Run("Parse reply", func(t *testing.T) {

		packet := NewPacket()
		packet.SetType("reply")
		packet.SetIdentifier(1337)
		packet.SetSequence(42)

		parsed := Parse(packet.Bytes())

		if parsed.Type != "reply" {
			t.Errorf("Expected ICMP type to be '%s' but got '%s'", "reply", parsed.Type)
		}

	})

}
