package tunnel

import "tholian-endpoint/protocols/dns"
import "strings"

func DecodePayload(response *dns.Packet) []byte {

	var payload []byte

	if response.Type == "response" {

		for a := 0; a < len(response.Answers); a++ {

			record := response.Answers[a]

			if record.Type == dns.TypeTXT && strings.Contains(record.Name, ".bytes.") {
				payload = append(payload, record.Data...)
			}

		}

	}

	return payload

}
