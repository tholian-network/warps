package icmp

import "net"

type Socket interface {
	Listen(host string) error
	ReadPacket() (Packet, net.Addr, error)
	WritePacket(Packet, net.Addr) error
	Close() error
}
