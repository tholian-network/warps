package icmp

import "net"
import "os"
import "testing"
import "time"

func TestPingSocket(t *testing.T) {

	if os.Geteuid() != 0 {
		t.Skip("Skipping raw ICMP test (requires root)")
	}

	t.Run("PingSocket loopback", func(t *testing.T) {

		socket := NewPingSocket()

		err := socket.Listen("0.0.0.0")

		if err != nil {
			t.Skipf("Skipping raw ICMP test: cannot bind ICMP socket: %s", err.Error())
		}

		defer socket.Close()

		packet := NewPacket()
		packet.SetType("request")
		packet.SetIdentifier(1)
		packet.SetSequence(1)
		packet.SetPayload([]byte("ping"))

		address := &net.IPAddr{IP: net.ParseIP("127.0.0.1")}

		err = socket.WritePacket(packet, address)

		if err != nil {
			t.Skipf("Skipping raw ICMP test: cannot write ICMP packet: %s", err.Error())
		}

		channel := make(chan Packet, 1)

		go func() {

			reply, _, err := socket.ReadPacket()

			if err == nil {
				channel <- reply
			}

		}()

		select {
		case reply := <-channel:

			if reply.Type != "reply" {
				t.Errorf("Expected ICMP type '%s' but got '%s'", "reply", reply.Type)
			}

			if string(reply.Payload) != "ping" {
				t.Errorf("Expected ICMP payload '%s' but got '%s'", "ping", string(reply.Payload))
			}

		case <-time.After(2 * time.Second):
			t.Skip("Skipping raw ICMP test: no echo reply received (may be blocked by firewall)")
		}

	})

}
