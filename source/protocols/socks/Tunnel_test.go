package socks

import "io"
import "testing"
import "time"

func TestTunnel(t *testing.T) {

	t.Run("Tunnel Dial to Proxy", func(t *testing.T) {

		echo_port := startEchoServer(t)
		proxy := NewProxy("localhost", 14434, nil)

		go func() {
			err := proxy.Listen()
			if err != nil {
				t.Errorf("Unexpected error '%s'", err.Error())
			}
		}()

		go func() {
			time.Sleep(2 * time.Second)
			proxy.Destroy()
		}()

		time.Sleep(100 * time.Millisecond)

		tunnel := NewTunnel("127.0.0.1", 14434)
		connection, err := tunnel.dial("127.0.0.1", uint16(echo_port))

		if err != nil {
			t.Fatalf("Cannot dial through SOCKS5 tunnel: %s", err.Error())
		}

		defer connection.Close()

		connection.Write([]byte("Hello"))

		data := make([]byte, 5)

		if _, err := io.ReadFull(connection, data); err != nil {
			t.Fatalf("Cannot read echo data: %s", err.Error())
		}

		if string(data) != "Hello" {
			t.Errorf("Expected echo '%s' but got '%s'", "Hello", string(data))
		}

		time.Sleep(2 * time.Second)

	})

}
