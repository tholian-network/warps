package types

import "strconv"

type Protocol string

const (
	ProtocolANY        Protocol = "any"
	ProtocolDNS        Protocol = "dns"
	ProtocolDNSoverTLS Protocol = "dns-over-tls"
	ProtocolHTTPS      Protocol = "https"
	ProtocolHTTP       Protocol = "http"
	ProtocolICMP       Protocol = "icmp"
	ProtocolSSH        Protocol = "ssh"
	ProtocolSOCKS      Protocol = "socks"
	ProtocolTCP        Protocol = "tcp"
	ProtocolUDP        Protocol = "udp"
	ProtocolWHOIS      Protocol = "whois"
)

func IsProtocol(value string) bool {

	var result bool

	if value == string(ProtocolANY) {
		result = true
	} else if value == string(ProtocolDNS) {
		result = true
	} else if value == string(ProtocolDNSoverTLS) {
		result = true
	} else if value == string(ProtocolHTTPS) {
		result = true
	} else if value == string(ProtocolHTTP) {
		result = true
	} else if value == string(ProtocolICMP) {
		result = true
	} else if value == string(ProtocolSSH) {
		result = true
	} else if value == string(ProtocolSOCKS) {
		result = true
	} else if value == string(ProtocolTCP) {
		result = true
	} else if value == string(ProtocolUDP) {
		result = true
	} else if value == string(ProtocolWHOIS) {
		result = true
	}

	return result

}

func ParseProtocol(value string) *Protocol {

	var result *Protocol = nil

	if value == "*" || value == "all" || value == "any" {
		protocol := Protocol(ProtocolANY)
		result = &protocol
	} else if value == "https" {
		protocol := Protocol(ProtocolHTTPS)
		result = &protocol
	} else if value == "http" {
		protocol := Protocol(ProtocolHTTP)
		result = &protocol
	} else if value == "dns" {
		protocol := Protocol(ProtocolDNS)
		result = &protocol
	} else if value == "dns-over-tls" || value == "dot" {
		protocol := Protocol(ProtocolDNSoverTLS)
		result = &protocol
	} else if value == "icmp" {
		protocol := Protocol(ProtocolICMP)
		result = &protocol
	} else if value == "ssh" {
		protocol := Protocol(ProtocolSSH)
		result = &protocol
	} else if value == "socks" || value == "socks4" || value == "socks5" {
		protocol := Protocol(ProtocolSOCKS)
		result = &protocol
	} else if value == "whois" {
		protocol := Protocol(ProtocolWHOIS)
		result = &protocol
	} else if value == "tcp" || value == "tcp4" || value == "tcp6" {
		protocol := Protocol(ProtocolTCP)
		result = &protocol
	} else if value == "udp" || value == "udp4" || value == "udp6" {
		protocol := Protocol(ProtocolUDP)
		result = &protocol
	}

	return result

}

func (protocol Protocol) Bytes() []byte {
	return []byte(protocol)
}

func (protocol Protocol) String() string {
	return string(protocol)
}

func (protocol Protocol) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(protocol))), nil
}

func (protocol *Protocol) UnmarshalJSON(data []byte) error {

	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	tmp := ParseProtocol(unquoted)

	if tmp != nil {
		*protocol = *tmp
	}

	return nil

}

func (protocol *Protocol) IsValid() bool {

	var result bool

	if IsProtocol(protocol.String()) {
		result = true
	}

	return result

}
