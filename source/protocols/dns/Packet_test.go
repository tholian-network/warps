package dns

import "testing"

func TestPacket(t *testing.T) {

	t.Run("NewPacket defaults", func(t *testing.T) {

		packet := NewPacket()

		if packet.Type != "" {
			t.Errorf("Expected type to be empty but got '%s'", packet.Type)
		}

		if len(packet.Questions) != 0 {
			t.Errorf("Expected no questions but got '%d'", len(packet.Questions))
		}

		if len(packet.Answers) != 0 {
			t.Errorf("Expected no answers but got '%d'", len(packet.Answers))
		}

		if packet.Server != nil {
			t.Errorf("Expected server to be nil")
		}

	})

	t.Run("SetType", func(t *testing.T) {

		packet := NewPacket()
		packet.SetType("query")

		if packet.Type != "query" {
			t.Errorf("Expected type '%s' but got '%s'", "query", packet.Type)
		}

		packet.SetType("response")

		if packet.Type != "response" {
			t.Errorf("Expected type '%s' but got '%s'", "response", packet.Type)
		}

		packet.SetType("invalid")

		if packet.Type != "response" {
			t.Errorf("Expected type '%s' but got '%s'", "response", packet.Type)
		}

	})

	t.Run("Query round-trip", func(t *testing.T) {

		query := NewPacket()
		query.SetType("query")
		query.AddQuestion(NewQuestion("example.com", TypeA))
		query.AddQuestion(NewQuestion("example.com", TypeAAAA))

		parsed := Parse(query.Bytes())

		if parsed.Type != "query" {
			t.Errorf("Expected type '%s' but got '%s'", "query", parsed.Type)
		}

		if len(parsed.Questions) != 2 {
			t.Errorf("Expected '%d' questions but got '%d'", 2, len(parsed.Questions))
		}

		if parsed.Questions[0].Name != "example.com" {
			t.Errorf("Expected name '%s' but got '%s'", "example.com", parsed.Questions[0].Name)
		}

		if parsed.Questions[0].Type != TypeA {
			t.Errorf("Expected type '%s' but got '%s'", "A", parsed.Questions[0].Type.String())
		}

		if parsed.Questions[1].Type != TypeAAAA {
			t.Errorf("Expected type '%s' but got '%s'", "AAAA", parsed.Questions[1].Type.String())
		}

	})

	t.Run("Response round-trip", func(t *testing.T) {

		record := NewRecord("example.com", TypeA)
		record.SetIPv4("1.3.3.7")

		response := NewPacket()
		response.SetType("response")
		response.SetResponseCode(ResponseCodeNoError)
		response.AddQuestion(NewQuestion("example.com", TypeA))
		response.AddAnswer(record)

		parsed := Parse(response.Bytes())

		if parsed.Type != "response" {
			t.Errorf("Expected type '%s' but got '%s'", "response", parsed.Type)
		}

		if len(parsed.Answers) != 1 {
			t.Errorf("Expected '%d' answers but got '%d'", 1, len(parsed.Answers))
		}

		if parsed.Answers[0].ToIPv4() != "1.3.3.7" {
			t.Errorf("Expected '%s' but got '%s'", "1.3.3.7", parsed.Answers[0].ToIPv4())
		}

	})

	t.Run("IsBlocked with blocked address", func(t *testing.T) {

		record := NewRecord("example.com", TypeA)
		record.SetIPv4("0.0.0.0")

		response := NewPacket()
		response.SetType("response")
		response.AddAnswer(record)

		if IsBlocked(response) != true {
			t.Errorf("Expected response to be blocked")
		}

	})

	t.Run("IsBlocked with allowed address", func(t *testing.T) {

		record := NewRecord("example.com", TypeA)
		record.SetIPv4("1.3.3.7")

		response := NewPacket()
		response.SetType("response")
		response.AddAnswer(record)

		if IsBlocked(response) != false {
			t.Errorf("Expected response to be allowed")
		}

	})

}
