package dns

import "tholian-endpoint/protocols/dns"
import "strings"

func toDNSPayload(response dns.Packet) []byte {

	var payload []byte

	if response.Type == "response" {

		for a := 0; a < len(response.Answers); a++ {

			record := response.Answers[a]

			if record.Type == dns.TypeTXT && strings.HasPrefix(record.Name, "bytes.") {

				tmp := strings.Split(record.Name, ".")

				// bytes.0-123.124.example.com
				if len(tmp) > 3 && tmp[0] == "bytes" && strings.Contains(tmp[1], "-") {
					payload = append(payload, record.Data)
				}

			}

		}

	}

	return payload

}
