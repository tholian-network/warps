//go:build !linux

package icmp

import "errors"
import "net"

type PingSocket struct {
	closed bool
}

func NewPingSocket() PingSocket {
	return PingSocket{}
}

func (socket *PingSocket) Listen(host string) error {
	return errors.New("ICMP is not supported on this platform")
}

func (socket *PingSocket) ReadPacket() (Packet, net.Addr, error) {
	return Packet{}, nil, errors.New("ICMP is not supported on this platform")
}

func (socket *PingSocket) WritePacket(packet Packet, address net.Addr) error {
	return errors.New("ICMP is not supported on this platform")
}

func (socket *PingSocket) Close() error {
	return nil
}
