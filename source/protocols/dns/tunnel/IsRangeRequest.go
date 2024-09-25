package tunnel

import "tholian-endpoint/protocols/dns"
import "strings"

func IsRangeRequest(query *dns.Packet) bool {

	var result bool = false

	if query.Type == "query" {

		for a := 0; a < len(query.Additionals); a++ {

			question := query.Additionals[a]

			if question.Type == dns.TypeURI && strings.Contains(question.Name, ".bytes.") {
				result = true
				break
			}

		}

	}

	return result

}
