package url

import "tholian-endpoint/types"
import net_url "net/url"

func ToHost(raw_url string) string {

	var result string

	url, err := net_url.Parse(raw_url)

	if err == nil {

		if types.IsIPv6AndPort(url.Host) {

			tmp_ip, tmp_port := types.ParseIPv6AndPort(url.Host)

			if tmp_ip != nil && tmp_port != 0 {
				result = tmp_ip.String()
			}

		} else if types.IsIPv6(url.Host) {

			tmp_ip := types.ParseIPv6(url.Host)

			if tmp_ip != nil {
				result = tmp_ip.String()
			}

		} else if types.IsIPv4AndPort(url.Host) {

			tmp_ip, tmp_port := types.ParseIPv4AndPort(url.Host)

			if tmp_ip != nil && tmp_port != 0 {
				result = tmp_ip.String()
			}

		} else if types.IsIPv4(url.Host) {

			tmp_ip := types.ParseIPv4(url.Host)

			if tmp_ip != nil {
				result = tmp_ip.String()
			}

		} else if types.IsDomainAndPort(url.Host) {

			tmp_domain, tmp_port := types.ParseDomainAndPort(url.Host)

			if tmp_domain != nil && tmp_port != 0 {
				result = tmp_domain.String()
			}

		} else if types.IsDomain(url.Host) {

			tmp_domain := types.ParseDomain(url.Host)

			if tmp_domain != nil {
				result = tmp_domain.String()
			}

		}

	}

	return result

}
