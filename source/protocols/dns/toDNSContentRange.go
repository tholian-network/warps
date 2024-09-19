package dns

import "tholian-endpoint/protocols/dns"
import "strconv"
import "strings"

func toDNSContentRange(response dns.Packet) (int, int, int) {

	var result_from int = 0
	var result_to   int = 0
	var result_size int = 0

	if response.Type == "response" {

		for a := 0; a < len(response.Answers); a++ {

			record := response.Answers[a]

			if record.Type == dns.TypeURI && strings.HasPrefix(record.Name, "bytes.") {

				tmp := strings.Split(record.Name, ".")

				// bytes.0-123.124.example.com
				if len(tmp) > 3 && tmp[0] == "bytes" && strings.Contains(tmp[1], "-") {

					num_from, err_from := strconv.ParseInt(tmp[1][0:strings.Index(tmp[1], "-")], 10, 64)
					num_to,   err_to   := strconv.ParseInt(tmp[1][strings.Index(tmp[1], "-")+1:], 10, 64)
					num_size, err_size := strconv.ParseInt(tmp[2], 10, 64)

					if err_from == nil && err_to == nil && err_size == nil {
						result_from = int(num_from)
						result_to   = int(num_to)
						result_size = int(num_size)
					}

				}

			}

		}

	}

	return result_from, result_to, result_size

}
