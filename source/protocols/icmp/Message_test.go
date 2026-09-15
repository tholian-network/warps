package icmp

import "testing"

func TestMessage(t *testing.T) {

	t.Run("Encode and Decode", func(t *testing.T) {

		message := NewMessage()
		message.SetIdentifier("request-id")
		message.SetType(MessageData)
		message.SetTarget("http://example.com/")
		message.SetPayload([]byte("Hello, world!"))
		message.SetKey(12345)
		message.SetProtocol(-1)
		message.SetTimeout(60)

		decoded := ParseMessage(message.Bytes())

		if decoded.Identifier != "request-id" {
			t.Errorf("Expected identifier '%s' but got '%s'", "request-id", decoded.Identifier)
		}

		if decoded.Type != MessageData {
			t.Errorf("Expected type '%d' but got '%d'", MessageData, decoded.Type)
		}

		if decoded.Target != "http://example.com/" {
			t.Errorf("Expected target '%s' but got '%s'", "http://example.com/", decoded.Target)
		}

		if string(decoded.Payload) != "Hello, world!" {
			t.Errorf("Expected payload '%s' but got '%s'", "Hello, world!", string(decoded.Payload))
		}

		if decoded.Key != 12345 {
			t.Errorf("Expected key '%d' but got '%d'", 12345, decoded.Key)
		}

		if decoded.Protocol != -1 {
			t.Errorf("Expected protocol '%d' but got '%d'", -1, decoded.Protocol)
		}

		if decoded.Timeout != 60 {
			t.Errorf("Expected timeout '%d' but got '%d'", 60, decoded.Timeout)
		}

		if decoded.IsValid() != true {
			t.Errorf("Expected message to be valid")
		}

	})

	t.Run("Default Magic", func(t *testing.T) {

		message := NewMessage()

		if message.Magic != int32(MessageMagic) {
			t.Errorf("Expected magic '%d' but got '%d'", MessageMagic, message.Magic)
		}

		if message.IsValid() != true {
			t.Errorf("Expected message to be valid")
		}

	})

	t.Run("Invalid Magic", func(t *testing.T) {

		message := NewMessage()
		message.Magic = 0

		if message.IsValid() != false {
			t.Errorf("Expected message to be invalid")
		}

	})

}
