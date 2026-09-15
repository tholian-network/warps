package http

import "net"

func requestTCP(ip string, port uint16, request Packet) Packet {

	var response Packet

	server := net.ParseIP(ip)

	if server != nil {

		connection, err1 := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   server,
			Port: int(port),
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
