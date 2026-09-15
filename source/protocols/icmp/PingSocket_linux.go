//go:build linux

package icmp

import "errors"
import "net"
import "syscall"

type PingSocket struct {
	socket int
	closed bool
}

func NewPingSocket() PingSocket {

	var socket PingSocket

	socket.socket = -1
	socket.closed = false

	return socket

}

func (socket *PingSocket) Listen(host string) error {

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_ICMP)

	if err != nil {
		return err
	}

	socket.socket = fd

	ip := net.ParseIP(host)

	if ip == nil || ip.To4() == nil {
		ip = net.IPv4zero
	}

	var address syscall.SockaddrInet4
	address.Port = 0
	copy(address.Addr[:], ip.To4())

	err = syscall.Bind(fd, &address)

	if err != nil {
		syscall.Close(fd)
		socket.socket = -1
		return err
	}

	return nil

}

func (socket *PingSocket) ReadPacket() (Packet, net.Addr, error) {

	var packet Packet
	var address net.Addr

	if socket.closed || socket.socket < 0 {
		return packet, address, errors.New("ICMP socket is closed")
	}

	buffer := make([]byte, 65535)

	length, remote, err := syscall.Recvfrom(socket.socket, buffer, 0)

	if err != nil {
		return packet, address, err
	}

	packet = Parse(buffer[0:length])

	if addr, ok := remote.(*syscall.SockaddrInet4); ok {
		address = &net.IPAddr{
			IP: net.IPv4(addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3]),
		}
	}

	return packet, address, nil

}

func (socket *PingSocket) WritePacket(packet Packet, address net.Addr) error {

	if socket.closed || socket.socket < 0 {
		return errors.New("ICMP socket is closed")
	}

	ipaddr, ok := address.(*net.IPAddr)

	if !ok || ipaddr.IP.To4() == nil {
		return errors.New("Invalid ICMP address")
	}

	var remote syscall.SockaddrInet4
	remote.Port = 0
	copy(remote.Addr[:], ipaddr.IP.To4())

	return syscall.Sendto(socket.socket, packet.Bytes(), 0, &remote)

}

func (socket *PingSocket) Close() error {

	if socket.closed {
		return nil
	}

	socket.closed = true

	if socket.socket >= 0 {
		err := syscall.Close(socket.socket)
		socket.socket = -1
		return err
	}

	return nil

}
