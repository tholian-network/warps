package socks

import "io"
import "net"
import "testing"
import "time"

func startEchoServer(t *testing.T) int {

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	})

	if err != nil {
		t.Fatalf("Cannot start echo server: %s", err.Error())
	}

	port := listener.Addr().(*net.TCPAddr).Port

	go func() {

		for {

			connection, err := listener.Accept()

			if err != nil {
				return
			}

			go func(connection net.Conn) {
				defer connection.Close()
				io.Copy(connection, connection)
			}(connection)

		}

	}()

	return port

}

func TestProxy(t *testing.T) {

	t.Run("Proxy with SOCKS5 client", func(t *testing.T) {

		echo_port := startEchoServer(t)
		proxy := NewProxy("localhost", 14433, nil)

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

		connection, err := net.Dial("tcp", "127.0.0.1:14433")

		if err != nil {
			t.Fatalf("Cannot connect to SOCKS5 proxy: %s", err.Error())
		}

		defer connection.Close()

		connection.Write([]byte{0x05, 0x01, 0x00})

		greeting := make([]byte, 2)

		if _, err := io.ReadFull(connection, greeting); err != nil {
			t.Fatalf("Cannot read greeting: %s", err.Error())
		}

		if greeting[0] != 0x05 || greeting[1] != 0x00 {
			t.Fatalf("Expected SOCKS5 no-auth greeting")
		}

		request := []byte{0x05, 0x01, 0x00, 0x01, 127, 0, 0, 1, byte(echo_port >> 8), byte(echo_port & 0xff)}
		connection.Write(request)

		reply := make([]byte, 10)

		if _, err := io.ReadFull(connection, reply); err != nil {
			t.Fatalf("Cannot read connect reply: %s", err.Error())
		}

		if reply[0] != 0x05 || reply[1] != 0x00 {
			t.Fatalf("Expected SOCKS5 connect success")
		}

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
