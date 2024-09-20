package tunnel

import "tholian-endpoint/protocols/dns"
import "encoding/json"
import "strings"

func EncodeHeaders(response *dns.Packet, headers map[string]string) bool {

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

			headers_record := dns.NewRecord("headers." + range_domain[6:], dns.TypeTXT)

			data := make([]byte, 0)

			if len(headers) > 0 {

				buffer, err := json.Marshal(headers)

				if err == nil {
					data = append(data, buffer...)
				}

			}

			headers_record.Data = data

			response.AddAnswer(headers_record)

		}

	}

	return result

}
