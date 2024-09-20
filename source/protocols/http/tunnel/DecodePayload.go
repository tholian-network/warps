package tunnel

import "tholian-endpoint/protocols/http"
import "encoding/base64"

func DecodePayload(packet *http.Packet) []byte {

	var payload []byte

	if packet.Type == "request" {

		if packet.Method == http.MethodGet {

			if packet.URL.Path == "/dns-query" {

				values := packet.URL.Query()

				if values.Get("dns") != "" {

					tmp, err := base64.URLEncoding.DecodeString(values.Get("dns"))

					if err == nil {
						payload = tmp
					}

				}

			}

		} else if packet.Method == http.MethodPost {

			content_type := packet.GetHeader("Content-Type")

			if content_type == "application/dns-message" {

				packet.Decode()

				data := make([]byte, 0)
				data = append(data, packet.Payload...)

				payload = data

			}

		}

	} else if packet.Type == "response" {

		if packet.Status == http.StatusOK {

			content_type := packet.GetHeader("Content-Type")

			if content_type == "application/dns-message" {

				packet.Decode()

				data := make([]byte, 0)
				data = append(data, packet.Payload...)

				payload = data

			}

		}

	}

	return payload

}
