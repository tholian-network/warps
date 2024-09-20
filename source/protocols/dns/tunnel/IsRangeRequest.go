package tunnel

import "tholian-endpoint/protocols/dns"
import "strings"

func IsRangeRequest(query *dns.Packet) bool {

	var result bool = false

	if query.Type == "query" {

		for a := 0; a < len(query.Answers); a++ {

			record := query.Answers[a]

			if record.Type == dns.TypeURI && strings.HasPrefix(record.Name, "bytes.") {

				tmp := strings.Split(record.Name, ".")

				// bytes.0-123.124.example.com
				if len(tmp) > 3 && tmp[0] == "bytes" && strings.Contains(tmp[1], "-") {
					result = true
					break
				}

			}

		}

	}

	return result

}
