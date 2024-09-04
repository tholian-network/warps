package net

import "net"

func ReadConnection(connection net.Conn) []byte {

	buffer := make([]byte, 0)

	for {

		chunk_buffer := make([]byte, 1 * 1024 * 1024)
		chunk_size, err3 := connection.Read(chunk_buffer)

		if err3 == nil {

			if chunk_size > 0 {
				buffer = append(buffer, chunk_buffer[0:chunk_size]...)
			}

		} else {
			break
		}

	}

	return buffer

}
