package tunnel

import "tholian-endpoint/protocols/dns"
import "strings"

func EncodePayload(response *dns.Packet, payload []byte) bool {

	var result bool = false

	if response.Type == "response" {

		range_domain := ""

		for a := 0; a < len(response.Answers); a++ {

			record := response.Answers[a]

			if record.Type == dns.TypeURI && strings.HasPrefix(record.Name, "bytes.") {

				tmp := strings.Split(record.Name, ".")

				// bytes.0-123/124.example.com
				if len(tmp) > 3 && tmp[0] == "bytes" && strings.Contains(tmp[1], "-") {
					range_domain = record.Name
					break
				}

			}

		}

		if range_domain != "" {

			payload_record := dns.NewRecord(range_domain, dns.TypeTXT)

			data := make([]byte, 0)
			data = append(data, payload...)

			payload_record.Data = data

			response.AddAnswer(payload_record)

		}

	}

	return result

}
