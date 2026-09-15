package http

import "crypto/tls"
import "crypto/x509"
import "net"
import "strconv"

func requestTLS(domain string, ip string, port uint16, request Packet) Packet {

	var response Packet

	server := net.ParseIP(ip)

	if server != nil {

		var certificates *x509.CertPool

		if request.Certificate != nil {

			certificates = x509.NewCertPool()
			certificates.AddCert(request.Certificate)

		} else {

			tmp, err := x509.SystemCertPool()

			if err == nil {
				certificates = tmp
			}

		}

		connection, err1 := tls.Dial("tcp", server.String() + ":" + strconv.FormatUint(uint64(port), 10), &tls.Config{
			// InsecureSkipVerify: true,
			RootCAs:    certificates,
			ServerName: domain,
		})

		if err1 == nil {

			request.SetHeader("Connection", "close")
			connection.Write(request.Bytes())

			response_buffer := make([]byte, 0)

			for {

				chunk_buffer := make([]byte, 1 * 1024 * 1024)
				chunk_size, err3 := connection.Read(chunk_buffer)

				if err3 == nil {

					if chunk_size > 0 {
						response_buffer = append(response_buffer, chunk_buffer[0:chunk_size]...)
					}

				} else {
					break
				}

			}

			if len(response_buffer) > 0 {
				response = Parse(response_buffer)
			}

			connection.Close()

		}

	}

	return response

}
