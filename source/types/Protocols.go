package types

type Protocol uint

const (
	ProtocolNone Protocol = iota
	ProtocolHTTPS
	ProtocolHTTP
	ProtocolDNS
	ProtocolICMP
	ProtocolTCP
)

func (protocol Protocol) String() string {

	var result string

	if protocol == ProtocolNone {
		result = ""
	} else if protocol == ProtocolHTTPS {
		result = "https"
	} else if protocol == ProtocolHTTP {
		result = "http"
	} else if protocol == ProtocolDNS {
		result = "dns"
	} else if protocol == ProtocolICMP {
		result = "icmp"
	} else if protocol == ProtocolTCP {
		result = "tcp"
	}

	return result

}
