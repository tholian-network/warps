package tunnel

import "tholian-endpoint/protocols/dns"
import net_url "net/url"
import "strings"

func DecodeURL(packet *dns.Packet) *net_url.URL {

	var result *net_url.URL = nil

	if packet.Type == "query" {

		for a := 0; a < len(packet.Additionals); a++ {

			record := packet.Additionals[a]

			if record.Type == dns.TypeURI && strings.Contains(record.Name, ".bytes.") {

				url, err := net_url.Parse(record.ToURL())

				if err == nil {
					result = url
					break
				}

			}

		}

	} else if packet.Type == "response" {

		for a := 0; a < len(packet.Answers); a++ {

			record := packet.Answers[a]

			if record.Type == dns.TypeURI && strings.Contains(record.Name, ".bytes.") {

				url, err := net_url.Parse(record.ToURL())

				if err == nil {
					result = url
					break
				}

			}

		}

	}

	return result

}
