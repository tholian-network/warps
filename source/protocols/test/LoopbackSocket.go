package test

import "tholian-warps/protocols/icmp"
import "errors"
import "net"

type LoopbackAddr string

func (addr LoopbackAddr) Network() string {
	return "loopback"
}

func (addr LoopbackAddr) String() string {
	return string(addr)
}

type LoopbackPacket struct {
	Packet  icmp.Packet
	Address net.Addr
}

type LoopbackSocket struct {
	Address net.Addr
	peer    *LoopbackSocket
	packets chan LoopbackPacket
	closed  bool
}

func NewLoopbackSocket() (*LoopbackSocket, *LoopbackSocket) {

	socket_a := &LoopbackSocket{
		Address: LoopbackAddr("loopback-a"),
		packets: make(chan LoopbackPacket, 1024),
	}

	socket_b := &LoopbackSocket{
		Address: LoopbackAddr("loopback-b"),
		packets: make(chan LoopbackPacket, 1024),
	}

	socket_a.peer = socket_b
	socket_b.peer = socket_a

	return socket_a, socket_b

}

func (socket *LoopbackSocket) Listen(host string) error {
	return nil
}

func (socket *LoopbackSocket) ReadPacket() (icmp.Packet, net.Addr, error) {

	incoming := <-socket.packets

	return incoming.Packet, incoming.Address, nil

}

func (socket *LoopbackSocket) WritePacket(packet icmp.Packet, address net.Addr) error {

	if socket.closed {
		return errors.New("loopback socket is closed")
	}

	socket.peer.packets <- LoopbackPacket{
		Packet:  packet,
		Address: address,
	}

	return nil

}

func (socket *LoopbackSocket) Close() error {

	socket.closed = true

	return nil

}
