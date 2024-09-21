package tunnel

import "tholian-endpoint/types"
import net_url "net/url"
import "slices"
import "strings"

func ToRecordName(url *net_url.URL) string {

	var result string

	if types.IsIPv6AndPort(url.Host) {

		tmp_ip, tmp_port := types.ParseIPv6AndPort(url.Host)

		if tmp_ip != nil && tmp_port != 0 {

			tmp := strings.Split(strings.Join(strings.Split(tmp_ip.String(), ":"), ""), "")
			slices.Reverse(tmp)

			result = strings.Join(tmp, ".") + ".ip6.arpa"

		}

	} else if types.IsIPv6(url.Host) {

		tmp_ip := types.ParseIPv6(url.Host)

		if tmp_ip != nil {

			tmp := strings.Split(strings.Join(strings.Split(tmp_ip.String(), ":"), ""), "")
			slices.Reverse(tmp)

			result = strings.Join(tmp, ".") + ".ip6.arpa"

		}

	} else if types.IsIPv4AndPort(url.Host) {

		tmp_ip, tmp_port := types.ParseIPv4AndPort(url.Host)

		if tmp_ip != nil && tmp_port != 0 {

			tmp := strings.Split(tmp_ip.String(), ".")
			slices.Reverse(tmp)

			return strings.Join(tmp, ".") + ".in-addr.arpa"

		}

	} else if types.IsIPv4(url.Host) {

		tmp_ip := types.ParseIPv4(url.Host)

		if tmp_ip != nil {

			tmp := strings.Split(tmp_ip.String(), ".")
			slices.Reverse(tmp)

			return strings.Join(tmp, ".") + ".in-addr.arpa"

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

	return result

}
