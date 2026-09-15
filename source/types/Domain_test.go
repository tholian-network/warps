package types

import "reflect"
import "testing"

func TestDomain(t *testing.T) {

	t.Run("IsDomain()", func(t *testing.T) {

		result1 := IsDomain("example.com")
		result2 := IsDomain("another.example.com")
		result3 := IsDomain("localhost")
		result4 := IsDomain("xy")

		if result1 != true {
			t.Errorf("Expected %t to be true", result1)
		}

		if result2 != true {
			t.Errorf("Expected %t to be true", result2)
		}

		if result3 != true {
			t.Errorf("Expected %t to be true", result3)
		}

		if result4 != false {
			t.Errorf("Expected %t to be false", result4)
		}

	})

	t.Run("IsDomainAndPort()", func(t *testing.T) {

		result1 := IsDomainAndPort("example.com")
		result2 := IsDomainAndPort("another.example.com")
		result3 := IsDomainAndPort("example.com:1337")
		result4 := IsDomainAndPort("another.example.com:13")

		if result1 != false {
			t.Errorf("Expected %t to be false", result1)
		}

		if result2 != false {
			t.Errorf("Expected %t to be false", result2)
		}

		if result3 != true {
			t.Errorf("Expected %t to be true", result3)
		}

		if result4 != true {
			t.Errorf("Expected %t to be true", result4)
		}

	})

	t.Run("ParseDomain()", func(t *testing.T) {

		domain1 := ParseDomain("example.com")
		domain2 := ParseDomain("another.example.com")
		domain3 := ParseDomain("localhost")
		domain4 := ParseDomain("xy")

		if domain1 == nil {
			t.Errorf("Expected %s to be valid", "example.com")
		} else if domain1.String() != "example.com" {
			t.Errorf("Expected %s to be %s", domain1.String(), "example.com")
		}

		if domain2 == nil {
			t.Errorf("Expected %s to be valid", "another.example.com")
		} else if domain2.String() != "another.example.com" {
			t.Errorf("Expected %s to be %s", domain2.String(), "another.example.com")
		}

		if domain3 == nil {
			t.Errorf("Expected %s to be valid", "localhost")
		} else if domain3.String() != "localhost" {
			t.Errorf("Expected %s to be %s", domain3.String(), "localhost")
		}

		if domain4 != nil {
			t.Errorf("Expected %s to be invalid", "xy")
		}

	})

	t.Run("ParseDomainAndPort()", func(t *testing.T) {

		domain1, port1 := ParseDomainAndPort("example.com")
		domain2, port2 := ParseDomainAndPort("another.example.com")
		domain3, port3 := ParseDomainAndPort("example.com:1337")
		domain4, port4 := ParseDomainAndPort("another.example.com:13")

		if domain1 != nil || port1 != 0 {
			t.Errorf("Expected %s to be invalid", "example.com")
		}

		if domain2 != nil || port2 != 0 {
			t.Errorf("Expected %s to be invalid", "another.example.com")
		}

		if domain3 == nil || port3 == 0 {
			t.Errorf("Expected %s to be valid", "example.com:1337")
		} else if domain3.String() != "example.com" || port3 != 1337 {
			t.Errorf("Expected %s:%d to be %s:%d", domain3.String(), port3, "example.com", 1337)
		}

		if domain4 == nil || port4 == 0 {
			t.Errorf("Expected %s to be valid", "another.example.com:13")
		} else if domain4.String() != "another.example.com" || port4 != 13 {
			t.Errorf("Expected %s:%d to be %s:%d", domain4.String(), port4, "another.example.com", 13)
		}

	})

	t.Run("Bytes()", func(t *testing.T) {

		domain1 := ParseDomain("example.com")
		domain2 := ParseDomain("another.example.com")
		domain3 := ParseDomain("localhost")

		want1 := []byte("example.com")
		want2 := []byte("another.example.com")
		want3 := []byte("localhost")

		if reflect.DeepEqual(domain1.Bytes(), want1) != true {
			t.Errorf("Expected %x to be %x", domain1.Bytes(), want1)
		}

		if reflect.DeepEqual(domain2.Bytes(), want2) != true {
			t.Errorf("Expected %x to be %x", domain2.Bytes(), want2)
		}

		if reflect.DeepEqual(domain3.Bytes(), want3) != true {
			t.Errorf("Expected %x to be %x", domain3.Bytes(), want3)
		}

	})

	t.Run("Scope()", func(t *testing.T) {

		domain1 := ParseDomain("example.com")
		domain2 := ParseDomain("another.example.com")
		domain3 := ParseDomain("1.0.31.172.in-addr.arpa")
		domain4 := ParseDomain("1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.e.f.ip6.arpa")
		domain5 := ParseDomain("fritz.box.local")
		domain6 := ParseDomain("localhost")

		if domain1.Scope() != "public" {
			t.Errorf("Expected %s to be %s", domain1.Scope(), "public")
		}

		if domain2.Scope() != "public" {
			t.Errorf("Expected %s to be %s", domain2.Scope(), "public")
		}

		if domain3.Scope() != "private" {
			t.Errorf("Expected %s to be %s", domain3.Scope(), "private")
		}

		if domain4.Scope() != "private" {
			t.Errorf("Expected %s to be %s", domain4.Scope(), "private")
		}

		if domain5.Scope() != "private" {
			t.Errorf("Expected %s to be %s", domain5.Scope(), "private")
		}

		if domain6.Scope() != "private" {
			t.Errorf("Expected %s to be %s", domain6.Scope(), "private")
		}

	})

}
