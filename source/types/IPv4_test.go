package types

import "reflect"
import "testing"

func TestIPv4(t *testing.T) {

	t.Run("IsIPv4()", func(t *testing.T) {

		is1 := IsIPv4("192.168.0.1")
		is2 := IsIPv4("255.255.255.255")
		is3 := IsIPv4("example.com")
		is4 := IsIPv4("1337.0.0.1")

		if is1 != true {
			t.Errorf("Expected %t to be true", is1)
		}

		if is2 != true {
			t.Errorf("Expected %t to be true", is2)
		}

		if is3 != false {
			t.Errorf("Expected %t to be false", is3)
		}

		if is4 != false {
			t.Errorf("Expected %t to be false", is4)
		}

	})

	t.Run("IsIPv4AndPort()", func(t *testing.T) {

		is1 := IsIPv4AndPort("192.168.0.1")
		is2 := IsIPv4AndPort("192.168.0.1:1337")
		is3 := IsIPv4AndPort("example.com")
		is4 := IsIPv4AndPort("example.com:1337")

		if is1 != false {
			t.Errorf("Expected %t to be false", is1)
		}

		if is2 != true {
			t.Errorf("Expected %t to be true", is2)
		}

		if is3 != false {
			t.Errorf("Expected %t to be false", is3)
		}

		if is4 != false {
			t.Errorf("Expected %t to be false", is4)
		}

	})

	t.Run("ParseIPv4()", func(t *testing.T) {

		ip1 := ParseIPv4("192.168.0.1")
		ip2 := ParseIPv4("255.255.255.255")
		ip3 := ParseIPv4("example.com")
		ip4 := ParseIPv4("1337.0.0.1")

		if ip1 == nil {
			t.Errorf("Expected %s to be valid", "192.168.0.1")
		} else if ip1.String() != "192.168.0.1" {
			t.Errorf("Expected %s to be %s", ip1.String(), "192.168.0.1")
		}

		if ip2 == nil {
			t.Errorf("Expected %s to be valid", "255.255.255.255")
		} else if ip2.String() != "255.255.255.255" {
			t.Errorf("Expected %s to be %s", ip2.String(), "255.255.255.255")
		}

		if ip3 != nil {
			t.Errorf("Expected %s to be invalid", "example.com")
		}

		if ip4 != nil {
			t.Errorf("Expected %s to be invalid", "1337.0.0.1")
		}

	})

	t.Run("ParseIPv4AndPort()", func(t *testing.T) {

		ip1, port1 := ParseIPv4AndPort("192.168.0.1")
		ip2, port2 := ParseIPv4AndPort("255.255.255.255")
		ip3, port3 := ParseIPv4AndPort("192.168.0.1:1337")
		ip4, port4 := ParseIPv4AndPort("255.255.255.255:13")

		if ip1 != nil || port1 != 0 {
			t.Errorf("Expected %s to be invalid", "192.168.0.1")
		}

		if ip2 != nil || port2 != 0 {
			t.Errorf("Expected %s to be invalid", "255.255.255.255")
		}

		if ip3 == nil || port3 == 0 {
			t.Errorf("Expected %s to be valid", "192.168.0.1:1337")
		} else if ip3.String() != "192.168.0.1" || port3 != 1337 {
			t.Errorf("Expected %s:%d to be %s:%d", ip3.String(), port3, "192.168.0.1", 1337)
		}

		if ip4 == nil || port4 == 0 {
			t.Errorf("Expected %s to be valid", "255.255.255.255:13")
		} else if ip4.String() != "255.255.255.255" || port4 != 13 {
			t.Errorf("Expected %s:%d to be %s:%d", ip4.String(), port4, "255.255.255.255", 13)
		}

	})

	t.Run("Bytes(bitmasks)", func(t *testing.T) {

		ip := ParseIPv4("255.255.255.255")

		want1 := []byte{128, 0, 0, 0}
		want2 := []byte{192, 0, 0, 0}
		want3 := []byte{224, 0, 0, 0}
		want4 := []byte{240, 0, 0, 0}
		want5 := []byte{248, 0, 0, 0}
		want6 := []byte{252, 0, 0, 0}
		want7 := []byte{254, 0, 0, 0}
		want8 := []byte{255, 0, 0, 0}

		want9 := []byte{255, 128, 0, 0}
		want10 := []byte{255, 192, 0, 0}
		want11 := []byte{255, 224, 0, 0}
		want12 := []byte{255, 240, 0, 0}
		want13 := []byte{255, 248, 0, 0}
		want14 := []byte{255, 252, 0, 0}
		want15 := []byte{255, 254, 0, 0}
		want16 := []byte{255, 255, 0, 0}

		want17 := []byte{255, 255, 128, 0}
		want18 := []byte{255, 255, 192, 0}
		want19 := []byte{255, 255, 224, 0}
		want20 := []byte{255, 255, 240, 0}
		want21 := []byte{255, 255, 248, 0}
		want22 := []byte{255, 255, 252, 0}
		want23 := []byte{255, 255, 254, 0}
		want24 := []byte{255, 255, 255, 0}

		want25 := []byte{255, 255, 255, 128}
		want26 := []byte{255, 255, 255, 192}
		want27 := []byte{255, 255, 255, 224}
		want28 := []byte{255, 255, 255, 240}
		want29 := []byte{255, 255, 255, 248}
		want30 := []byte{255, 255, 255, 252}
		want31 := []byte{255, 255, 255, 254}
		want32 := []byte{255, 255, 255, 255}

		if reflect.DeepEqual(ip.Bytes(1), want1) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(1), want1)
		}

		if reflect.DeepEqual(ip.Bytes(2), want2) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(2), want2)
		}

		if reflect.DeepEqual(ip.Bytes(3), want3) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(3), want3)
		}

		if reflect.DeepEqual(ip.Bytes(4), want4) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(4), want4)
		}

		if reflect.DeepEqual(ip.Bytes(5), want5) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(5), want5)
		}

		if reflect.DeepEqual(ip.Bytes(6), want6) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(6), want6)
		}

		if reflect.DeepEqual(ip.Bytes(7), want7) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(7), want7)
		}

		if reflect.DeepEqual(ip.Bytes(8), want8) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(8), want8)
		}

		if reflect.DeepEqual(ip.Bytes(9), want9) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(9), want9)
		}

		if reflect.DeepEqual(ip.Bytes(10), want10) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(10), want10)
		}

		if reflect.DeepEqual(ip.Bytes(11), want11) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(11), want11)
		}

		if reflect.DeepEqual(ip.Bytes(12), want12) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(12), want12)
		}

		if reflect.DeepEqual(ip.Bytes(13), want13) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(13), want13)
		}

		if reflect.DeepEqual(ip.Bytes(14), want14) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(14), want14)
		}

		if reflect.DeepEqual(ip.Bytes(15), want15) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(15), want15)
		}

		if reflect.DeepEqual(ip.Bytes(16), want16) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(16), want16)
		}

		if reflect.DeepEqual(ip.Bytes(17), want17) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(17), want17)
		}

		if reflect.DeepEqual(ip.Bytes(18), want18) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(18), want18)
		}

		if reflect.DeepEqual(ip.Bytes(19), want19) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(19), want19)
		}

		if reflect.DeepEqual(ip.Bytes(20), want20) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(20), want20)
		}

		if reflect.DeepEqual(ip.Bytes(21), want21) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(21), want21)
		}

		if reflect.DeepEqual(ip.Bytes(22), want22) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(22), want22)
		}

		if reflect.DeepEqual(ip.Bytes(23), want23) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(23), want23)
		}

		if reflect.DeepEqual(ip.Bytes(24), want24) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(24), want24)
		}

		if reflect.DeepEqual(ip.Bytes(25), want25) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(25), want25)
		}

		if reflect.DeepEqual(ip.Bytes(26), want26) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(26), want26)
		}

		if reflect.DeepEqual(ip.Bytes(27), want27) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(27), want27)
		}

		if reflect.DeepEqual(ip.Bytes(28), want28) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(28), want28)
		}

		if reflect.DeepEqual(ip.Bytes(29), want29) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(29), want29)
		}

		if reflect.DeepEqual(ip.Bytes(30), want30) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(30), want30)
		}

		if reflect.DeepEqual(ip.Bytes(31), want31) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(31), want31)
		}

		if reflect.DeepEqual(ip.Bytes(32), want32) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(32), want32)
		}

	})

	t.Run("Bytes()", func(t *testing.T) {

		ip1 := ParseIPv4("192.168.0.1")
		ip2 := ParseIPv4("0.0.0.0")

		want11 := []byte{192, 168, 0, 1}
		want12 := []byte{192, 168, 0, 0}
		want13 := []byte{192, 128, 0, 0}
		want14 := []byte{192, 160, 0, 0}
		want15 := []byte{192, 168, 0, 0}

		want2 := []byte{0, 0, 0, 0}

		if reflect.DeepEqual(ip1.Bytes(32), want11) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(32), want11)
		}

		if reflect.DeepEqual(ip1.Bytes(24), want12) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(24), want12)
		}

		if reflect.DeepEqual(ip1.Bytes(9), want13) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(9), want13)
		}

		if reflect.DeepEqual(ip1.Bytes(11), want14) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(11), want14)
		}

		if reflect.DeepEqual(ip1.Bytes(13), want15) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(13), want15)
		}

		if reflect.DeepEqual(ip2.Bytes(32), want2) != true {
			t.Errorf("Expected %x to be %x", ip2.Bytes(32), want2)
		}

		if reflect.DeepEqual(ip2.Bytes(0), want2) != true {
			t.Errorf("Expected %x to be %x", ip2.Bytes(0), want2)
		}

	})

	t.Run("Scopes()", func(t *testing.T) {

		ip1 := ParseIPv4("192.168.0.1")
		ip2 := ParseIPv4("0.0.0.0")
		ip3 := ParseIPv4("169.254.0.1")
		ip4 := ParseIPv4("172.24.23.45")
		ip5 := ParseIPv4("1.3.3.7")

		if ip1.Scope() != "private" {
			t.Errorf("Expected %s to be private", ip1.String())
		}

		if ip2.Scope() != "private" {
			t.Errorf("Expected %s to be private", ip2.String())
		}

		if ip3.Scope() != "private" {
			t.Errorf("Expected %s to be private", ip3.String())
		}

		if ip4.Scope() != "private" {
			t.Errorf("Expected %s to be private", ip4.String())
		}

		if ip5.Scope() != "public" {
			t.Errorf("Expected %s to be public", ip5.String())
		}

	})

	t.Run("String()", func(t *testing.T) {

		ip1 := ParseIPv4("192.168.0.1")
		ip2 := ParseIPv4("255.255.255.255")
		ip3 := ParseIPv4("0.0.0.0")

		if ip1.String() != "192.168.0.1" {
			t.Errorf("Expected %s to be %s", ip1.String(), "192.168.0.1")
		}

		if ip2.String() != "255.255.255.255" {
			t.Errorf("Expected %s to be %s", ip2.String(), "255.255.255.255")
		}

		if ip3.String() != "0.0.0.0" {
			t.Errorf("Expected %s to be %s", ip3.String(), "0.0.0.0")
		}

	})

}
