package socks

import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "errors"
import "io"
import "net"
import "strconv"

type Tunnel struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func NewTunnel(host string, port uint16) Tunnel {

	var tunnel Tunnel

	tunnel.Host = host
	tunnel.Port = port

	return tunnel

}

func (tunnel *Tunnel) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	response = dns.NewPacket()
	response.SetType("response")
	response.SetIdentifier(query.Identifier)
	response.SetResponseCode(dns.ResponseCodeNonExistDomain)
	response.Flags.RecursionAvailable = true

	for q := 0; q < len(query.Questions); q++ {
		response.AddQuestion(query.Questions[q])
	}

	return response

}

func (tunnel *Tunnel) RequestPacket(request http.Packet) http.Packet {

	var response http.Packet

	if request.URL == nil {

		response = http.NewPacket()
		response.SetStatus(http.StatusNotFound)
		response.SetPayload([]byte{})

	} else {

		host := request.URL.Hostname()
		port, err1 := strconv.Atoi(request.URL.Port())

		if err1 != nil || port == 0 {

			if request.URL.Scheme == "https" {
				port = 443
			} else {
				port = 80
			}

		}

		connection, err2 := tunnel.dial(host, uint16(port))

		if err2 == nil {

			request.SetHeader("Connection", "close")
			connection.Write(request.Bytes())

			buffer := make([]byte, 0)

			for {

				chunk := make([]byte, 1*1024*1024)
				size, err3 := connection.Read(chunk)

				if err3 == nil {

					if size > 0 {
						buffer = append(buffer, chunk[0:size]...)
					}

				} else {
					break
				}

			}

			connection.Close()

			if len(buffer) > 0 {
				response = http.Parse(buffer)
			} else {
				response = http.NewPacket()
				response.SetURL(*request.URL)
				response.SetStatus(http.StatusNotFound)
				response.SetPayload([]byte{})
			}

		} else {

			response = http.NewPacket()
			response.SetURL(*request.URL)
			response.SetStatus(http.StatusNotFound)
			response.SetPayload([]byte{})

		}

	}

	return response

}

func (tunnel *Tunnel) dial(host string, port uint16) (net.Conn, error) {

	address := net.JoinHostPort(tunnel.Host, strconv.FormatUint(uint64(tunnel.Port), 10))
	connection, err := net.Dial("tcp", address)

	if err != nil {
		return nil, err
	}

	_, err = connection.Write([]byte{0x05, 0x01, 0x00})

	if err != nil {
		connection.Close()
		return nil, err
	}

	reply := make([]byte, 2)

	if _, err = io.ReadFull(connection, reply); err != nil {
		connection.Close()
		return nil, err
	}

	if reply[0] != 0x05 || reply[1] != 0x00 {
		connection.Close()
		return nil, errors.New("SOCKS5 handshake failed")
	}

	request := []byte{0x05, 0x01, 0x00}

	if ip := net.ParseIP(host); ip != nil && ip.To4() != nil {
		request = append(request, 0x01)
		request = append(request, ip.To4()...)
	} else if ip := net.ParseIP(host); ip != nil && ip.To16() != nil {
		request = append(request, 0x04)
		request = append(request, ip.To16()...)
	} else {
		request = append(request, 0x03)
		request = append(request, byte(len(host)))
		request = append(request, []byte(host)...)
	}

	request = append(request, byte(port>>8), byte(port&0xff))

	if _, err = connection.Write(request); err != nil {
		connection.Close()
		return nil, err
	}

	reply = make([]byte, 4)

	if _, err = io.ReadFull(connection, reply); err != nil {
		connection.Close()
		return nil, err
	}

	if reply[0] != 0x05 || reply[1] != 0x00 {
		connection.Close()
		return nil, errors.New("SOCKS5 connect failed")
	}

	var skip_length int

	if reply[3] == 0x01 {
		skip_length = 4 + 2
	} else if reply[3] == 0x04 {
		skip_length = 16 + 2
	} else if reply[3] == 0x03 {

		length := make([]byte, 1)

		if _, err = io.ReadFull(connection, length); err != nil {
			connection.Close()
			return nil, err
		}

		skip_length = int(length[0]) + 2

	} else {
		connection.Close()
		return nil, errors.New("SOCKS5 invalid address type")
	}

	skip := make([]byte, skip_length)

	if _, err = io.ReadFull(connection, skip); err != nil {
		connection.Close()
		return nil, err
	}

	return connection, nil

}
