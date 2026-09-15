package types

import "reflect"
import "testing"

func TestIPv6(t *testing.T) {

	t.Run("IsIPv6()", func(t *testing.T) {

		is1 := IsIPv6("[fe80:0000:0000:0000:1234:5678:9abc:0def]")
		is2 := IsIPv6("fe80:0000:0000:0000:1234:5678:9abc:0def")
		is3 := IsIPv6("fe80::1")
		is4 := IsIPv6("fe80:133337::1")
		is5 := IsIPv6("example.com")
		is6 := IsIPv6("192.168.0.1")

		if is1 != true {
			t.Errorf("Expected %t to be true", is1)
		}

		if is2 != true {
			t.Errorf("Expected %t to be true", is2)
		}

		if is3 != true {
			t.Errorf("Expected %t to be true", is3)
		}

		if is4 != false {
			t.Errorf("Expected %t to be false", is4)
		}

		if is5 != false {
			t.Errorf("Expected %t to be false", is5)
		}

		if is6 != false {
			t.Errorf("Expected %t to be false", is6)
		}

	})

	t.Run("IsIPv6AndPort()", func(t *testing.T) {

		is1 := IsIPv6AndPort("[fe80:0000:0000:0000:1234:5678:9abc:0def]:1337")
		is2 := IsIPv6AndPort("[fe80::1]:1337")
		is3 := IsIPv6AndPort("fe80:0000:0000:0000:1234:5678:9abc:0def")

		if is1 != true {
			t.Errorf("Expected %t to be true", is1)
		}

		if is2 != true {
			t.Errorf("Expected %t to be true", is2)
		}

		if is3 != false {
			t.Errorf("Expected %t to be false", is3)
		}

	})

	t.Run("ParseIPv6()", func(t *testing.T) {

		ip1 := ParseIPv6("[fe80:0000:0000:0000:1234:5678:9abc:0def]")
		ip2 := ParseIPv6("fe80:0000:0000:0000:1234:5678:9abc:0def")
		ip3 := ParseIPv6("fe80::1")
		ip4 := ParseIPv6("example.com")
		ip5 := ParseIPv6("192.168.0.1")

		if ip1 == nil {
			t.Errorf("Expected %s to be valid", "fe80:0000:0000:0000:1234:5678:9abc:0def")
		} else if ip1.String() != "fe80:0000:0000:0000:1234:5678:9abc:0def" {
			t.Errorf("Expected %s to be %s", ip1.String(), "fe80:0000:0000:0000:1234:5678:9abc:0def")
		}

		if ip2 == nil {
			t.Errorf("Expected %s to be valid", "fe80:0000:0000:0000:1234:5678:9abc:0def")
		} else if ip2.String() != "fe80:0000:0000:0000:1234:5678:9abc:0def" {
			t.Errorf("Expected %s to be %s", ip2.String(), "fe80:0000:0000:0000:1234:5678:9abc:0def")
		}

		if ip3 == nil {
			t.Errorf("Expected %s to be valid", "fe80:0000:0000:0000:0000:0000:0000:0001")
		} else if ip3.String() != "fe80:0000:0000:0000:0000:0000:0000:0001" {
			t.Errorf("Expected %s to be %s", ip3.String(), "fe80:0000:0000:0000:0000:0000:0000:0001")
		}

		if ip4 != nil {
			t.Errorf("Expected %s to be invalid", "example.com")
		}

		if ip5 != nil {
			t.Errorf("Expected %s to be invalid", "192.168.0.1")
		}

	})

	t.Run("ParseIPv6AndPort()", func(t *testing.T) {

		ip1, port1 := ParseIPv6AndPort("[fe80:0000:0000:0000:1234:5678:9abc:0def]:1337")
		ip2, port2 := ParseIPv6AndPort("[fe80::1]:1337")
		ip3, port3 := ParseIPv6AndPort("fe80:0000:0000:0000:1234:5678:9abc:0def")

		if ip1 == nil || port1 == 0 {
			t.Errorf("Expected %s to be valid", "[fe80:0000:0000:0000:1234:5678:9abc:0def]:1337")
		} else if ip1.String() != "fe80:0000:0000:0000:1234:5678:9abc:0def" || port1 != 1337 {
			t.Errorf("Expected %s:%d to be %s:%d", ip1.String(), port1, "fe80:0000:0000:0000:1234:5678:9abc:0def", 1337)
		}

		if ip2 == nil || port2 == 0 {
			t.Errorf("Expected %s to be valid", "[fe80::1]:1337")
		} else if ip2.String() != "fe80:0000:0000:0000:0000:0000:0000:0001" || port2 != 1337 {
			t.Errorf("Expected %s:%d to be %s:%d", ip2.String(), port2, "fe80:0000:0000:0000:0000:0000:0000:0001", 1337)
		}

		if ip3 != nil || port3 != 0 {
			t.Errorf("Expected %s to be invalid", "fe80:0000:0000:0000:1234:5678:9abc:0def")
		}

	})

	t.Run("Bytes(bitmasks)", func(t *testing.T) {

		ip := ParseIPv6("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")

		want1 := []byte{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want2 := []byte{192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want3 := []byte{224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want4 := []byte{240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want5 := []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want6 := []byte{252, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want7 := []byte{254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want8 := []byte{255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want9 := []byte{255, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want10 := []byte{255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want11 := []byte{255, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want12 := []byte{255, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want13 := []byte{255, 248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want14 := []byte{255, 252, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want15 := []byte{255, 254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want16 := []byte{255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		want17 := []byte{255, 255, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want18 := []byte{255, 255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want19 := []byte{255, 255, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want20 := []byte{255, 255, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want21 := []byte{255, 255, 248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want22 := []byte{255, 255, 252, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want23 := []byte{255, 255, 254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want24 := []byte{255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want25 := []byte{255, 255, 255, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want26 := []byte{255, 255, 255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want27 := []byte{255, 255, 255, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want28 := []byte{255, 255, 255, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want29 := []byte{255, 255, 255, 248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want30 := []byte{255, 255, 255, 252, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want31 := []byte{255, 255, 255, 254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want32 := []byte{255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		want33 := []byte{255, 255, 255, 255, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want34 := []byte{255, 255, 255, 255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want35 := []byte{255, 255, 255, 255, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want36 := []byte{255, 255, 255, 255, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want37 := []byte{255, 255, 255, 255, 248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want38 := []byte{255, 255, 255, 255, 252, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want39 := []byte{255, 255, 255, 255, 254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want40 := []byte{255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want41 := []byte{255, 255, 255, 255, 255, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want42 := []byte{255, 255, 255, 255, 255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want43 := []byte{255, 255, 255, 255, 255, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want44 := []byte{255, 255, 255, 255, 255, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want45 := []byte{255, 255, 255, 255, 255, 248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want46 := []byte{255, 255, 255, 255, 255, 252, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want47 := []byte{255, 255, 255, 255, 255, 254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want48 := []byte{255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		want49 := []byte{255, 255, 255, 255, 255, 255, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want50 := []byte{255, 255, 255, 255, 255, 255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want51 := []byte{255, 255, 255, 255, 255, 255, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want52 := []byte{255, 255, 255, 255, 255, 255, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want53 := []byte{255, 255, 255, 255, 255, 255, 248, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want54 := []byte{255, 255, 255, 255, 255, 255, 252, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want55 := []byte{255, 255, 255, 255, 255, 255, 254, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want56 := []byte{255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want57 := []byte{255, 255, 255, 255, 255, 255, 255, 128, 0, 0, 0, 0, 0, 0, 0, 0}
		want58 := []byte{255, 255, 255, 255, 255, 255, 255, 192, 0, 0, 0, 0, 0, 0, 0, 0}
		want59 := []byte{255, 255, 255, 255, 255, 255, 255, 224, 0, 0, 0, 0, 0, 0, 0, 0}
		want60 := []byte{255, 255, 255, 255, 255, 255, 255, 240, 0, 0, 0, 0, 0, 0, 0, 0}
		want61 := []byte{255, 255, 255, 255, 255, 255, 255, 248, 0, 0, 0, 0, 0, 0, 0, 0}
		want62 := []byte{255, 255, 255, 255, 255, 255, 255, 252, 0, 0, 0, 0, 0, 0, 0, 0}
		want63 := []byte{255, 255, 255, 255, 255, 255, 255, 254, 0, 0, 0, 0, 0, 0, 0, 0}
		want64 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0}

		want65 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 128, 0, 0, 0, 0, 0, 0, 0}
		want66 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 192, 0, 0, 0, 0, 0, 0, 0}
		want67 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 224, 0, 0, 0, 0, 0, 0, 0}
		want68 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 240, 0, 0, 0, 0, 0, 0, 0}
		want69 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 248, 0, 0, 0, 0, 0, 0, 0}
		want70 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 252, 0, 0, 0, 0, 0, 0, 0}
		want71 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 254, 0, 0, 0, 0, 0, 0, 0}
		want72 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0}
		want73 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 128, 0, 0, 0, 0, 0, 0}
		want74 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 192, 0, 0, 0, 0, 0, 0}
		want75 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 224, 0, 0, 0, 0, 0, 0}
		want76 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 240, 0, 0, 0, 0, 0, 0}
		want77 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 248, 0, 0, 0, 0, 0, 0}
		want78 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 252, 0, 0, 0, 0, 0, 0}
		want79 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 254, 0, 0, 0, 0, 0, 0}
		want80 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0}

		want81 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 128, 0, 0, 0, 0, 0}
		want82 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 192, 0, 0, 0, 0, 0}
		want83 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 224, 0, 0, 0, 0, 0}
		want84 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 240, 0, 0, 0, 0, 0}
		want85 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 248, 0, 0, 0, 0, 0}
		want86 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 252, 0, 0, 0, 0, 0}
		want87 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254, 0, 0, 0, 0, 0}
		want88 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0}
		want89 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 128, 0, 0, 0, 0}
		want90 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 192, 0, 0, 0, 0}
		want91 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 224, 0, 0, 0, 0}
		want92 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 240, 0, 0, 0, 0}
		want93 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 248, 0, 0, 0, 0}
		want94 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 252, 0, 0, 0, 0}
		want95 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254, 0, 0, 0, 0}
		want96 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0}

		want97 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 128, 0, 0, 0}
		want98 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 192, 0, 0, 0}
		want99 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 224, 0, 0, 0}
		want100 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 240, 0, 0, 0}
		want101 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 248, 0, 0, 0}
		want102 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 252, 0, 0, 0}
		want103 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254, 0, 0, 0}
		want104 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0}
		want105 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 128, 0, 0}
		want106 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 192, 0, 0}
		want107 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 224, 0, 0}
		want108 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 240, 0, 0}
		want109 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 248, 0, 0}
		want110 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 252, 0, 0}
		want111 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254, 0, 0}
		want112 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0}

		want113 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 128, 0}
		want114 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 192, 0}
		want115 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 224, 0}
		want116 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 240, 0}
		want117 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 248, 0}
		want118 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 252, 0}
		want119 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254, 0}
		want120 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0}
		want121 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 128}
		want122 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 192}
		want123 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 224}
		want124 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 240}
		want125 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 248}
		want126 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 252}
		want127 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254}
		want128 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}

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

		if reflect.DeepEqual(ip.Bytes(33), want33) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(33), want33)
		}

		if reflect.DeepEqual(ip.Bytes(34), want34) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(34), want34)
		}

		if reflect.DeepEqual(ip.Bytes(35), want35) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(35), want35)
		}

		if reflect.DeepEqual(ip.Bytes(36), want36) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(36), want36)
		}

		if reflect.DeepEqual(ip.Bytes(37), want37) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(37), want37)
		}

		if reflect.DeepEqual(ip.Bytes(38), want38) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(38), want38)
		}

		if reflect.DeepEqual(ip.Bytes(39), want39) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(39), want39)
		}

		if reflect.DeepEqual(ip.Bytes(40), want40) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(40), want40)
		}

		if reflect.DeepEqual(ip.Bytes(41), want41) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(41), want41)
		}

		if reflect.DeepEqual(ip.Bytes(42), want42) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(42), want42)
		}

		if reflect.DeepEqual(ip.Bytes(43), want43) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(43), want43)
		}

		if reflect.DeepEqual(ip.Bytes(44), want44) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(44), want44)
		}

		if reflect.DeepEqual(ip.Bytes(45), want45) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(45), want45)
		}

		if reflect.DeepEqual(ip.Bytes(46), want46) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(46), want46)
		}

		if reflect.DeepEqual(ip.Bytes(47), want47) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(47), want47)
		}

		if reflect.DeepEqual(ip.Bytes(48), want48) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(48), want48)
		}

		if reflect.DeepEqual(ip.Bytes(49), want49) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(49), want49)
		}

		if reflect.DeepEqual(ip.Bytes(50), want50) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(50), want50)
		}

		if reflect.DeepEqual(ip.Bytes(51), want51) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(51), want51)
		}

		if reflect.DeepEqual(ip.Bytes(52), want52) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(52), want52)
		}

		if reflect.DeepEqual(ip.Bytes(53), want53) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(53), want53)
		}

		if reflect.DeepEqual(ip.Bytes(54), want54) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(54), want54)
		}

		if reflect.DeepEqual(ip.Bytes(55), want55) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(55), want55)
		}

		if reflect.DeepEqual(ip.Bytes(56), want56) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(56), want56)
		}

		if reflect.DeepEqual(ip.Bytes(57), want57) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(57), want57)
		}

		if reflect.DeepEqual(ip.Bytes(58), want58) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(58), want58)
		}

		if reflect.DeepEqual(ip.Bytes(59), want59) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(59), want59)
		}

		if reflect.DeepEqual(ip.Bytes(60), want60) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(60), want60)
		}

		if reflect.DeepEqual(ip.Bytes(61), want61) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(61), want61)
		}

		if reflect.DeepEqual(ip.Bytes(62), want62) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(62), want62)
		}

		if reflect.DeepEqual(ip.Bytes(63), want63) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(63), want63)
		}

		if reflect.DeepEqual(ip.Bytes(64), want64) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(64), want64)
		}

		if reflect.DeepEqual(ip.Bytes(65), want65) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(65), want65)
		}

		if reflect.DeepEqual(ip.Bytes(66), want66) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(66), want66)
		}

		if reflect.DeepEqual(ip.Bytes(67), want67) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(67), want67)
		}

		if reflect.DeepEqual(ip.Bytes(68), want68) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(68), want68)
		}

		if reflect.DeepEqual(ip.Bytes(69), want69) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(69), want69)
		}

		if reflect.DeepEqual(ip.Bytes(70), want70) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(70), want70)
		}

		if reflect.DeepEqual(ip.Bytes(71), want71) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(71), want71)
		}

		if reflect.DeepEqual(ip.Bytes(72), want72) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(72), want72)
		}

		if reflect.DeepEqual(ip.Bytes(73), want73) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(73), want73)
		}

		if reflect.DeepEqual(ip.Bytes(74), want74) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(74), want74)
		}

		if reflect.DeepEqual(ip.Bytes(75), want75) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(75), want75)
		}

		if reflect.DeepEqual(ip.Bytes(76), want76) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(76), want76)
		}

		if reflect.DeepEqual(ip.Bytes(77), want77) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(77), want77)
		}

		if reflect.DeepEqual(ip.Bytes(78), want78) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(78), want78)
		}

		if reflect.DeepEqual(ip.Bytes(79), want79) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(79), want79)
		}

		if reflect.DeepEqual(ip.Bytes(80), want80) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(80), want80)
		}

		if reflect.DeepEqual(ip.Bytes(81), want81) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(81), want81)
		}

		if reflect.DeepEqual(ip.Bytes(82), want82) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(82), want82)
		}

		if reflect.DeepEqual(ip.Bytes(83), want83) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(83), want83)
		}

		if reflect.DeepEqual(ip.Bytes(84), want84) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(84), want84)
		}

		if reflect.DeepEqual(ip.Bytes(85), want85) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(85), want85)
		}

		if reflect.DeepEqual(ip.Bytes(86), want86) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(86), want86)
		}

		if reflect.DeepEqual(ip.Bytes(87), want87) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(87), want87)
		}

		if reflect.DeepEqual(ip.Bytes(88), want88) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(88), want88)
		}

		if reflect.DeepEqual(ip.Bytes(89), want89) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(89), want89)
		}

		if reflect.DeepEqual(ip.Bytes(90), want90) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(90), want90)
		}

		if reflect.DeepEqual(ip.Bytes(91), want91) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(91), want91)
		}

		if reflect.DeepEqual(ip.Bytes(92), want92) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(92), want92)
		}

		if reflect.DeepEqual(ip.Bytes(93), want93) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(93), want93)
		}

		if reflect.DeepEqual(ip.Bytes(94), want94) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(94), want94)
		}

		if reflect.DeepEqual(ip.Bytes(95), want95) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(95), want95)
		}

		if reflect.DeepEqual(ip.Bytes(96), want96) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(96), want96)
		}

		if reflect.DeepEqual(ip.Bytes(97), want97) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(97), want97)
		}

		if reflect.DeepEqual(ip.Bytes(98), want98) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(98), want98)
		}

		if reflect.DeepEqual(ip.Bytes(99), want99) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(99), want99)
		}

		if reflect.DeepEqual(ip.Bytes(100), want100) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(100), want100)
		}

		if reflect.DeepEqual(ip.Bytes(101), want101) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(101), want101)
		}

		if reflect.DeepEqual(ip.Bytes(102), want102) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(102), want102)
		}

		if reflect.DeepEqual(ip.Bytes(103), want103) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(103), want103)
		}

		if reflect.DeepEqual(ip.Bytes(104), want104) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(104), want104)
		}

		if reflect.DeepEqual(ip.Bytes(105), want105) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(105), want105)
		}

		if reflect.DeepEqual(ip.Bytes(106), want106) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(106), want106)
		}

		if reflect.DeepEqual(ip.Bytes(107), want107) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(107), want107)
		}

		if reflect.DeepEqual(ip.Bytes(108), want108) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(108), want108)
		}

		if reflect.DeepEqual(ip.Bytes(109), want109) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(109), want109)
		}

		if reflect.DeepEqual(ip.Bytes(110), want110) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(110), want110)
		}

		if reflect.DeepEqual(ip.Bytes(111), want111) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(111), want111)
		}

		if reflect.DeepEqual(ip.Bytes(112), want112) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(112), want112)
		}

		if reflect.DeepEqual(ip.Bytes(113), want113) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(113), want113)
		}

		if reflect.DeepEqual(ip.Bytes(114), want114) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(114), want114)
		}

		if reflect.DeepEqual(ip.Bytes(115), want115) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(115), want115)
		}

		if reflect.DeepEqual(ip.Bytes(116), want116) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(116), want116)
		}

		if reflect.DeepEqual(ip.Bytes(117), want117) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(117), want117)
		}

		if reflect.DeepEqual(ip.Bytes(118), want118) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(118), want118)
		}

		if reflect.DeepEqual(ip.Bytes(119), want119) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(119), want119)
		}

		if reflect.DeepEqual(ip.Bytes(120), want120) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(120), want120)
		}

		if reflect.DeepEqual(ip.Bytes(121), want121) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(121), want121)
		}

		if reflect.DeepEqual(ip.Bytes(122), want122) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(122), want122)
		}

		if reflect.DeepEqual(ip.Bytes(123), want123) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(123), want123)
		}

		if reflect.DeepEqual(ip.Bytes(124), want124) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(124), want124)
		}

		if reflect.DeepEqual(ip.Bytes(125), want125) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(125), want125)
		}

		if reflect.DeepEqual(ip.Bytes(126), want126) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(126), want126)
		}

		if reflect.DeepEqual(ip.Bytes(127), want127) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(127), want127)
		}

		if reflect.DeepEqual(ip.Bytes(128), want128) != true {
			t.Errorf("Expected %x to be %x", ip.Bytes(128), want128)
		}

	})

	t.Run("Bytes()", func(t *testing.T) {

		ip1 := ParseIPv6("2a00:20:3000:209f:6dd8:c7bc:c51a::")
		ip2 := ParseIPv6("0001:0012:0123:1234::")

		want11 := []byte{42, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want12 := []byte{42, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want13 := []byte{42, 0, 0, 32, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		want14 := []byte{42, 0, 0, 32, 48, 0, 32, 159, 0, 0, 0, 0, 0, 0, 0, 0}

		want2 := []byte{0, 1, 0, 18, 1, 35, 18, 52, 0, 0, 0, 0, 0, 0, 0, 0}

		if reflect.DeepEqual(ip1.Bytes(16), want11) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(16), want11)
		}

		if reflect.DeepEqual(ip1.Bytes(32), want12) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(32), want12)
		}

		if reflect.DeepEqual(ip1.Bytes(48), want13) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(48), want13)
		}

		if reflect.DeepEqual(ip1.Bytes(64), want14) != true {
			t.Errorf("Expected %x to be %x", ip1.Bytes(64), want14)
		}

		if reflect.DeepEqual(ip2.Bytes(64), want2) != true {
			t.Errorf("Expected %x to be %x", ip2.Bytes(64), want2)
		}

	})

	t.Run("Scopes()", func(t *testing.T) {

		ip1 := ParseIPv6("::1")
		ip2 := ParseIPv6("fe80:0000:0000:0000:1337::1337")
		ip3 := ParseIPv6("1337::12:345")

		if ip1.Scope() != "private" {
			t.Errorf("Expected %s to be private", ip1.String())
		}

		if ip2.Scope() != "private" {
			t.Errorf("Expected %s to be private", ip2.String())
		}

		if ip3.Scope() != "public" {
			t.Errorf("Expected %s to be public", ip3.String())
		}

	})

	t.Run("String()", func(t *testing.T) {

		ip1 := ParseIPv6("2a00:20:3000:209f:6dd8:c7bc:c51a::")
		ip2 := ParseIPv6("0001:0012:0123:1234::")

		if ip1.String() != "2a00:0020:3000:209f:6dd8:c7bc:c51a:0000" {
			t.Errorf("Expected %s to be %s", ip1.String(), "2a00:0020:3000:209f:6dd8:c7bc:c51a:0000")
		}

		if ip2.String() != "0001:0012:0123:1234:0000:0000:0000:0000" {
			t.Errorf("Expected %s to be %s", ip2.String(), "0001:0012:0123:1234:0000:0000:0000:0000")
		}

	})

}
