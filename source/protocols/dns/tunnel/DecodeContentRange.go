package tunnel

import "tholian-endpoint/protocols/dns"
import "strconv"
import "strings"

func DecodeContentRange(packet *dns.Packet) (string, int, int, int) {

	var url  string = ""
	var from int = -1
	var to   int = -1
	var size int = -1

	if packet.Type == "query" || packet.Type == "response" {

		for a := 0; a < len(packet.Answers); a++ {

			record := packet.Answers[a]

			if record.Type == dns.TypeURI && strings.HasPrefix(record.Name, "bytes.") {

				tmp := strings.Split(record.Name, ".")

				// bytes.0-123.124.example.com
				if len(tmp) > 3 && tmp[0] == "bytes" && strings.Contains(tmp[1], "-") {

					tmp_from := tmp[1][0:strings.Index(tmp[1], "-")]
					tmp_to   := tmp[1][strings.Index(tmp[1], "-")+1:]
					tmp_size := tmp[2]

					if tmp_from != "" && tmp_to != "" {

						num_from, err_from := strconv.ParseInt(tmp_from, 10, 64)
						num_to,   err_to   := strconv.ParseInt(tmp_to, 10, 64)

						if err_from == nil && err_to == nil {
							url  = record.ToURL()
							from = int(num_from)
							to   = int(num_to)
						}

						if tmp_size != "x" {

							num_size, err_size := strconv.ParseInt(tmp_size, 10, 64)

							if err_size == nil {
								size = int(num_size)
							}

						}

					} else if tmp_from != "" {

						num_from, err_from := strconv.ParseInt(tmp_from, 10, 64)

						if err_from == nil {
							url  = record.ToURL()
							from = int(num_from)
						}

					}

				}

			}

		}

	}

	return url, from, to, size

}
