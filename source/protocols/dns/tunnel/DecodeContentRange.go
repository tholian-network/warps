package tunnel

import "tholian-endpoint/protocols/dns"
import utils_http "tholian-warps/utils/protocols/http"
import "encoding/json"
import "strconv"
import "strings"

func DecodeContentRange(packet *dns.Packet) (int, int, int) {

	var from int = -1
	var to   int = -1
	var size int = -1

	if packet.Type == "request" {

		headers_domain := ""

		for q := 0; q < len(packet.Questions); q++ {

			question := packet.Questions[q]

			if question.Type == dns.TypeTXT && strings.Contains(question.Name, ".headers.") {
				headers_domain = question.Name
				break
			}

		}

		if from == -1 && to == -1 && headers_domain != "" {

			tmp1 := headers_domain[0:strings.Index(headers_domain, ".headers.")]

			if strings.Contains(tmp1, "-") {

				tmp2 := strings.Split(tmp1, "-")

				if len(tmp2) == 2 && tmp2[1] != "" {

					// frame request: 0-511.headers.domain.tld
					num_from, err_from := strconv.ParseInt(tmp2[0], 10, 64)
					num_to,   err_to   := strconv.ParseInt(tmp2[1], 10, 64)

					if err_from == nil && err_to == nil {
						from = int(num_from)
						to   = int(num_to)
					}

				} else if len(tmp2) == 2 {

					// first request: 0-.headers.domain.tld
					num_from, err_from := strconv.ParseInt(tmp2[0], 10, 64)

					if err_from == nil {
						from = int(num_from)
					}

				}

			}

		}

	} else if packet.Type == "response" {

		headers_domain := ""
		headers := make(map[string]string)

		for a := 0; a < len(packet.Answers); a++ {

			record := packet.Answers[a]

			if record.Type == dns.TypeTXT && strings.Contains(record.Name, ".headers.") {
				headers_domain = record.Name
				json.Unmarshal(record.Data, &headers)
				break
			}

		}

		content_range, ok1 := headers["Content-Range"]
		content_length, ok2 := headers["Content-Length"]

		if ok1 == true {
			from, to, size = utils_http.ParseContentRange(content_range)
		}

		if from == -1 && to == -1 && headers_domain != "" {

			tmp1 := headers_domain[0:strings.Index(headers_domain, ".headers.")]

			if strings.Contains(tmp1, "-") {

				tmp2 := strings.Split(tmp1, "-")

				if len(tmp2) == 2 && tmp2[1] != "" {

					// frame response: 0-511.headers.domain.tld
					num_from, err_from := strconv.ParseInt(tmp2[0], 10, 64)
					num_to,   err_to   := strconv.ParseInt(tmp2[1], 10, 64)

					if err_from == nil && err_to == nil {
						from = int(num_from)
						to   = int(num_to)
					}

				} else if len(tmp2) == 2 {

					// first response: 0-.headers.domain.tld
					num_from, err_from := strconv.ParseInt(tmp2[0], 10, 64)

					if err_from == nil {
						from = int(num_from)
					}

				}

			}

		}

		if size == -1 && ok2 == true {

			num_length, err_length := strconv.ParseInt(content_length, 10, 64)

			if err_length == nil {
				size = int(num_length)
			}

		}

	}

	return from, to, size

}
