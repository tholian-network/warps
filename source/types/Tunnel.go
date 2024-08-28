package types

type Tunnel struct {
	Host     string   `json:"host"`
	Port     uint16   `json:"port"`
	Protocol Protocol `json:"protocol"`
}

func (tunnel *Tunnel) SetPort(port uint16) {

	if port > 0 && port < 65535 {
		tunnel.Port = port
	}

}

func (tunnel *Tunnel) SetProtocol(protocol Protocol) {
	tunnel.Protocol = protocol
}
