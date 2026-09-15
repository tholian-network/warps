package http

import "tholian-warps/utils/compress/zstd"
import "tholian-warps/protocols/dns"
import "tholian-warps/types"
import "bytes"
import "crypto/x509"
import "compress/bzip2"
import "compress/flate"
import "compress/gzip"
import "io"
import "net/url"
import "strconv"
import "strings"

type ResolveMethod func(domain string) (dns.Packet, error)

type Packet struct {
	Type            string            `json:"type"`
	URL             *url.URL          `json:"url"`
	Method          Method            `json:"method"`
	Status          Status            `json:"status"`
	Headers         map[string]string `json:"headers"`
	Payload         []byte            `json:"payload"`
	Server          *types.Server     `json:"server"`
	Certificate     *x509.Certificate `json:"certificate"`
	Encoding        Encoding          `json:"encoding"`
	resolve_method  ResolveMethod
}

func NewPacket() Packet {

	var packet Packet

	tmp, _ := url.Parse("/")

	packet.Type     = "request"
	packet.Headers  = make(map[string]string)
	packet.Payload  = make([]byte, 0)
	packet.URL      = tmp
	packet.Method   = Method("")
	packet.Server   = nil

	packet.Headers["Accept-Encoding"] = "bzip2, gzip, deflate, zstd"
	packet.Headers["Content-Encoding"] = "identity"
	packet.Encoding = EncodingIdentity

	packet.resolve_method = ResolveMethod(func(domain string) (dns.Packet, error) {
		return dns.Resolve(domain)
	})

	return packet

}

func (packet *Packet) Bytes() []byte {

	var buffer []byte

	if packet.Type == "request" {

		buffer = append(buffer, []byte(packet.Method.String() + " " + packet.URL.String() + " HTTP/1.1\r\n")...)

		if len(packet.Headers) > 0 {

			for key, val := range packet.Headers {
				buffer = append(buffer, []byte(key + ": " + val + "\r\n")...)
			}

			buffer = append(buffer, []byte("\r\n")...)

		} else {
			buffer = append(buffer, []byte("\r\n")...)
			buffer = append(buffer, []byte("\r\n")...)
		}

		if len(packet.Payload) > 0 {
			buffer = append(buffer, packet.Payload...)
		}

	} else if packet.Type == "response" {

		buffer = append(buffer, []byte("HTTP/1.1 " + strconv.Itoa(int(packet.Status)) + " " + packet.Status.String())...)

		if len(packet.Headers) > 0 {

			for key, val := range packet.Headers {
				buffer = append(buffer, []byte(key + ": " + val + "\r\n")...)
			}

			buffer = append(buffer, []byte("\r\n")...)

		} else {
			buffer = append(buffer, []byte("\r\n")...)
			buffer = append(buffer, []byte("\r\n")...)
		}

		if len(packet.Payload) > 0 {
			buffer = append(buffer, packet.Payload...)
		}

	}

	return buffer

}

func (packet *Packet) Resolve() {

	if packet.Type == "request" && packet.Server == nil && packet.URL != nil {

		domain := ""
		host := ""
		addresses := make([]string, 0)
		port := 0
		protocol := types.ProtocolANY

		if types.IsIPv6AndPort(packet.URL.Host) {

			tmp_ip, tmp_port := types.ParseIPv6AndPort(packet.URL.Host)

			if tmp_ip != nil && tmp_port != 0 {
				domain    = tmp_ip.String()
				host      = "[" + tmp_ip.String() + "]:" + strconv.FormatUint(uint64(tmp_port), 10)
				addresses = append(addresses, tmp_ip.String())
				port      = int(tmp_port)
			}

		} else if types.IsIPv6(packet.URL.Host) {

			tmp_ip := types.ParseIPv6(packet.URL.Host)

			if tmp_ip != nil {
				domain    = tmp_ip.String()
				host      = "[" + tmp_ip.String() + "]"
				addresses = append(addresses, tmp_ip.String())
			}

		} else if types.IsIPv4AndPort(packet.URL.Host) {

			tmp_ip, tmp_port := types.ParseIPv4AndPort(packet.URL.Host)

			if tmp_ip != nil && tmp_port != 0 {
				domain    = tmp_ip.String()
				host      = tmp_ip.String() + ":" + strconv.FormatUint(uint64(tmp_port), 10)
				addresses = append(addresses, tmp_ip.String())
				port      = int(tmp_port)
			}

		} else if types.IsIPv4(packet.URL.Host) {

			tmp_ip := types.ParseIPv4(packet.URL.Host)

			if tmp_ip != nil {
				domain    = tmp_ip.String()
				host      = tmp_ip.String()
				addresses = append(addresses, tmp_ip.String())
			}

		} else if types.IsDomainAndPort(packet.URL.Host) {

			tmp_domain, tmp_port := types.ParseDomainAndPort(packet.URL.Host)

			if tmp_domain != nil && tmp_port != 0 {

				domain = tmp_domain.String()
				host   = tmp_domain.String() + ":" + strconv.FormatUint(uint64(tmp_port), 10)
				port   = int(tmp_port)

				response, err := packet.resolve_method(tmp_domain.String())

				if err == nil && response.Type == "response" {

					for a := 0; a < len(response.Answers); a++ {

						record := response.Answers[a]

						if record.Type == dns.TypeA {
							addresses = append(addresses, record.ToIPv4())
						} else if record.Type == dns.TypeAAAA {
							addresses = append(addresses, record.ToIPv6())
						}

					}

				}

			}

		} else if types.IsDomain(packet.URL.Host) {

			tmp_domain := types.ParseDomain(packet.URL.Host)

			if tmp_domain != nil {

				domain = tmp_domain.String()
				host   = tmp_domain.String()

				response, err := packet.resolve_method(tmp_domain.String())

				if err == nil && response.Type == "response" {

					for a := 0; a < len(response.Answers); a++ {

						record := response.Answers[a]

						if record.Type == dns.TypeA {
							addresses = append(addresses, record.ToIPv4())
						} else if record.Type == dns.TypeAAAA {
							addresses = append(addresses, record.ToIPv6())
						}

					}

				}

			}

		}

		if packet.URL.Scheme == "https" {

			protocol = types.ProtocolHTTPS

			if host != "" {
				packet.Headers["Host"] = host
			}

			if port == 0 {
				port = 443
			}

		} else if packet.URL.Scheme == "http" {

			protocol = types.ProtocolHTTP

			if host != "" {
				packet.Headers["Host"] = host
			}

			if port == 0 {
				port = 80
			}

		}

		if domain != "" && len(addresses) > 0 && port != 0 {

			packet.SetServer(types.Server{
				Domain:    domain,
				Addresses: addresses,
				Port:      uint16(port),
				Protocol:  protocol,
				Schema:    "",
			})

		}

	}

}

func (packet *Packet) Decode() bool {

	var result bool = false

	if packet.Encoding == EncodingIdentity {

		result = true

	} else if packet.Encoding == EncodingBzip2 {

		buffer := bytes.NewBuffer(packet.Payload)
		reader := bzip2.NewReader(buffer)

		decoded, err2 := io.ReadAll(reader)

		if err2 == nil {
			packet.Headers["Content-Encoding"] = "identity"
			packet.Encoding = EncodingIdentity
			packet.Payload = decoded
			result = true
		}

	} else if packet.Encoding == EncodingDeflate {

		buffer := bytes.NewBuffer(packet.Payload)
		reader := flate.NewReader(buffer)

		decoded, err2 := io.ReadAll(reader)
		reader.Close()

		if err2 == nil {
			packet.Headers["Content-Encoding"] = "identity"
			packet.Encoding = EncodingIdentity
			packet.Payload = decoded
			result = true
		}

	} else if packet.Encoding == EncodingGzip {

		buffer := bytes.NewBuffer(packet.Payload)
		reader, err1 := gzip.NewReader(buffer)

		if err1 == nil {

			decoded, err2 := io.ReadAll(reader)
			reader.Close()

			if err2 == nil {
				packet.Headers["Content-Encoding"] = "identity"
				packet.Encoding = EncodingIdentity
				packet.Payload = decoded
				result = true
			}

		}

	} else if packet.Encoding == EncodingZstd {

		buffer := bytes.NewBuffer(packet.Payload)
		reader := zstd.NewReader(buffer)

		decoded, err1 := io.ReadAll(reader)

		if err1 == nil {
			packet.Headers["Content-Encoding"] = "identity"
			packet.Encoding = EncodingIdentity
			packet.Payload = decoded
			result = true
		}

	}

	return result

}

func (packet *Packet) Encode(encoding Encoding) bool {

	var result bool = false

	if packet.Decode() {

		if encoding == EncodingIdentity {

			packet.Headers["Content-Encoding"] = "identity"
			packet.Encoding = EncodingIdentity
			result = true

		} else if encoding == EncodingBzip2 {

			// TODO: compress/bzip2 doesn't support io.Writer yet
			result = false

		} else if encoding == EncodingDeflate {

			var encoded bytes.Buffer

			reader := bytes.NewBuffer(packet.Payload)
			writer, err1 := flate.NewWriter(&encoded, flate.BestSpeed)

			if err1 == nil {

				_, err2 := io.Copy(writer, reader)

				writer.Flush()
				writer.Close()

				if err2 == nil {
					packet.Headers["Content-Encoding"] = "deflate"
					packet.Encoding = EncodingDeflate
					packet.Payload = encoded.Bytes()
					result = true
				}

			}

		} else if encoding == EncodingGzip {

			var encoded bytes.Buffer

			reader := bytes.NewBuffer(packet.Payload)
			writer := gzip.NewWriter(&encoded)

			_, err1 := io.Copy(writer, reader)

			writer.Flush()
			writer.Close()

			if err1 == nil {
				packet.Headers["Content-Encoding"] = "gzip"
				packet.Encoding = EncodingGzip
				packet.Payload = encoded.Bytes()
				result = true
			}

		} else if encoding == EncodingZstd {

			// TODO: internal/zstd doesn't support io.Writer yet
			result = false

		}

	}

	return result

}

func (packet *Packet) GetHeader(name string) string {

	var result string

	if strings.Contains(name, "-") {

		chunks := strings.Split(name, "-")

		for c := 0; c < len(chunks); c++ {
			chunks[c] = strings.ToUpper(chunks[c][0:1]) + strings.ToLower(chunks[c][1:])
		}

		tmp, ok := packet.Headers[strings.Join(chunks, "-")]

		if ok == true {
			result = tmp
		}

	} else {

		tmp, ok := packet.Headers[strings.ToUpper(name[0:1]) + strings.ToLower(name[1:])]

		if ok == true {
			result = tmp
		}

	}

	return result

}

func (packet *Packet) SetCertificate(value x509.Certificate) {
	packet.Certificate = &value
}

func (packet *Packet) SetHeader(name string, value string) {

	if strings.Contains(name, "-") {

		chunks := strings.Split(name, "-")

		for c := 0; c < len(chunks); c++ {
			chunks[c] = strings.ToUpper(chunks[c][0:1]) + strings.ToLower(chunks[c][1:])
		}

		packet.Headers[strings.Join(chunks, "-")] = value

	} else {

		packet.Headers[strings.ToUpper(name[0:1]) + strings.ToLower(name[1:])] = value

	}

}

func (packet *Packet) SetMethod(value Method) {

	packet.Type = "request"
	packet.Method = value
	packet.Status = Status(0)

	packet.Headers["Accept-Encoding"] = "bzip2, gzip, deflate, zstd"

}

func (packet *Packet) SetPayload(value []byte) {

	packet.Headers["Content-Length"] = strconv.Itoa(len(value))
	packet.Payload = value

}

func (packet *Packet) SetResolveMethod(callback ResolveMethod) {
	packet.resolve_method = callback
}

func (packet *Packet) SetURL(value url.URL) {
	packet.URL = &value
}

func (packet *Packet) SetServer(value types.Server) {
	packet.Server = &value
}

func (packet *Packet) SetStatus(value Status) {

	packet.Type = "response"
	packet.Method = Method("")
	packet.Status = value

	delete(packet.Headers, "Accept-Encoding")

}
