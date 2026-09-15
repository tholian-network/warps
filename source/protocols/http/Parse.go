package http

import "bytes"
import "net/url"
import "strconv"
import "strings"

func Parse(buffer []byte) Packet {

	packet := NewPacket()

	if len(buffer) > 0 && bytes.Index(buffer, []byte("\r\n\r\n")) != -1 {

		chunks := bytes.Split(buffer, []byte("\r\n\r\n"))

		if len(chunks) == 2 {

			lines := strings.Split(string(chunks[0]), "\r\n")

			if len(lines) >= 1 {

				first_line := lines[0]

				if strings.HasSuffix(first_line, " HTTP/1.1") && strings.Index(first_line, " ") != -1 {

					request_method := Method(first_line[0:strings.Index(first_line, " ")])
					request_url, err := url.Parse(first_line[strings.Index(first_line, " ")+1:len(first_line)-9])

					if err == nil {

						packet.Type = "request"
						packet.URL = request_url
						packet.Method = request_method

					} else {

						request_url, err = url.Parse("/")

						packet.Type = "request"
						packet.URL = request_url
						packet.Method = MethodGet

					}

				} else if strings.HasPrefix(first_line, "HTTP/1.1 ") {

					status_line := first_line[9:]
					num, err := strconv.ParseInt(status_line[0:strings.Index(status_line, " ")], 10, 64)

					if err == nil {
						packet.Type = "response"
						packet.URL = nil
						packet.Method = Method("")
						packet.Status = Status(int(num))
					}

				}

				for l := 1; l < len(lines); l++ {

					line := lines[l]

					if strings.Contains(line, ":") {

						header_name := strings.ToLower(strings.TrimSpace(line[0:strings.Index(line, ":")]))
						header_value := strings.TrimSpace(line[strings.Index(line, ":")+1:])

						packet.SetHeader(header_name, header_value)

					}

				}

			}

			if len(chunks[1]) > 0 {
				packet.SetPayload(chunks[1])
			}

			content_encoding, ok := packet.Headers["Content-Encoding"]

			if ok {

				if content_encoding == "identity" {
					packet.Encoding = EncodingIdentity
				} else if content_encoding == "bzip2" {
					packet.Encoding = EncodingBzip2
				} else if content_encoding == "deflate" {
					packet.Encoding = EncodingDeflate
				} else if content_encoding == "gzip" {
					packet.Encoding = EncodingGzip
				} else if content_encoding == "zstd" {
					packet.Encoding = EncodingZstd
				}

			}

		}

	}

	return packet

}
