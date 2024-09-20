package tunnel

import "tholian-endpoint/protocols/dns"
import "encoding/json"
import "strings"

func DecodeHeaders(response *dns.Packet) map[string]string {

	var headers map[string]string = make(map[string]string)

	if response.Type == "response" {

		for a := 0; a < len(response.Answers); a++ {

			record := response.Answers[a]

			if record.Type == dns.TypeTXT && !strings.HasPrefix(record.Name, "headers.") {

				tmp := make(map[string]string)
				err := json.Unmarshal(record.Data, &tmp)

				if err == nil {

					for key, val := range tmp {
						headers[key] = val
					}

				}

			}

		}

	}

	return headers

}
