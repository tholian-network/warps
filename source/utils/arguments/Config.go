package arguments

import "tholian-endpoint/types"
import net_url "net/url"
import "strconv"
import "strings"

type Config struct {
	Host     string         `json:"host"`
	Port     uint16         `json:"port"`
	Protocol types.Protocol `json:"protocol"`
}

func ParseConfig(raw_url string) *Config {

	var result *Config = nil

	if raw_url == "any" {

		tmp := Config{
			Host:     "0.0.0.0",
			Port:     0,
			Protocol: types.ProtocolANY,
		}

		result = &tmp

	} else {

		url, err := net_url.Parse(raw_url)

		if err == nil {

			host := ""
			port := uint16(0)
			protocol := types.ProtocolANY

			switch url.Scheme {
			case "dns":
				protocol = types.ProtocolDNS
			case "http":
				protocol = types.ProtocolHTTP
			case "https":
				protocol = types.ProtocolHTTPS
			case "socks":
				protocol = types.ProtocolSOCKS
			default:
				protocol = types.ProtocolANY
			}

			if strings.HasPrefix(url.Host, "localhost:") {

				num_port, err_port := strconv.ParseUint(url.Host[10:], 10, 16)

				if err_port == nil {
					host = "0.0.0.0"
					port = uint16(num_port)
				}

			} else if url.Host == "localhost" {

				host = "0.0.0.0"
				port = 0

			} else if types.IsIPv6AndPort(url.Host) {

				tmp_ipv6, tmp_port := types.ParseIPv6AndPort(url.Host)

				if tmp_ipv6 != nil && tmp_port != 0 {
					host = tmp_ipv6.String()
					port = uint16(tmp_port)
				}

			} else if types.IsIPv6(url.Host) {

				tmp_ipv6 := types.ParseIPv6(url.Host)

				if tmp_ipv6 != nil {
					host = tmp_ipv6.String()
					port = 0
				}

			} else if types.IsIPv4AndPort(url.Host) {

				tmp_ipv4, tmp_port := types.ParseIPv4AndPort(url.Host)

				if tmp_ipv4 != nil && tmp_port != 0 {
					host = tmp_ipv4.String()
					port = uint16(tmp_port)
				}

			} else if types.IsIPv4(url.Host) {

				tmp_ipv4 := types.ParseIPv4(url.Host)

				if tmp_ipv4 != nil {
					host = tmp_ipv4.String()
					port = 0
				}

			} else if types.IsDomain(url.Host) {

				// TODO: Domain Parsing
				// TODO: This should also look in /etc/letsencrypt
				// for matching SSL certificates automatically

			}

			if host != "" && protocol != types.ProtocolANY {

				if port == 0 {

					switch protocol {
					case types.ProtocolDNS:
						port = 1053
					case types.ProtocolHTTP:
						port = 1080
					case types.ProtocolHTTPS:
						port = 1443
					case types.ProtocolSOCKS:
						port = 1090
					case types.ProtocolANY:
						port = 0
					}

				}

				tmp := Config{
					Host:     host,
					Port:     port,
					Protocol: protocol,
				}

				result = &tmp

			}

		}

	}

	return result

}

func (config Config) String() string {

	if types.IsIPv6(config.Host) {
		return types.Protocol(config.Protocol).String() + "://[" + config.Host + "]:" + strconv.FormatUint(uint64(config.Port), 10)
	} else {
		return types.Protocol(config.Protocol).String() + "://" + config.Host + ":" + strconv.FormatUint(uint64(config.Port), 10)
	}

}
