package url

import "tholian-endpoint/types"
import net_url "net/url"
import "strconv"

func ToHostAndPort(url *net_url.URL) string {

	var result string

	domain := ""
	port := 0

	if types.IsIPv6AndPort(url.Host) {

		tmp_ip, tmp_port := types.ParseIPv6AndPort(url.Host)

		if tmp_ip != nil && tmp_port != 0 {
			domain = "[" + tmp_ip.String() + "]"
			port   = int(tmp_port)
		}

	} else if types.IsIPv6(url.Host) {

		tmp_ip := types.ParseIPv6(url.Host)

		if tmp_ip != nil {
			domain = "[" + tmp_ip.String() + "]"
		}

	} else if types.IsIPv4AndPort(url.Host) {

		tmp_ip, tmp_port := types.ParseIPv4AndPort(url.Host)

		if tmp_ip != nil && tmp_port != 0 {
			domain = tmp_ip.String()
			port   = int(tmp_port)
		}

	} else if types.IsIPv4(url.Host) {

		tmp_ip := types.ParseIPv4(url.Host)

		if tmp_ip != nil {
			domain = tmp_ip.String()
		}

	} else if types.IsDomainAndPort(url.Host) {

		tmp_domain, tmp_port := types.ParseDomainAndPort(url.Host)

		if tmp_domain != nil && tmp_port != 0 {
			domain = tmp_domain.String()
			port   = int(tmp_port)
		}

	} else if types.IsDomain(url.Host) {

		tmp_domain := types.ParseDomain(url.Host)

		if tmp_domain != nil {
			domain = tmp_domain.String()
		}

	}

	if url.Scheme == "dns" {

		if port == 0 {
			port = 53
		}

	} else if url.Scheme == "dot" || url.Scheme == "dns-over-tls" {

		if port == 0 {
			port = 853
		}

	} else if url.Scheme == "https" {

		if port == 0 {
			port = 443
		}

	} else if url.Scheme == "http" {

		if port == 0 {
			port = 80
		}

	} else if url.Scheme == "socks" {

		if port == 0 {
			port = 1080
		}

	} else if url.Scheme == "ssh" {

		if port == 0 {
			port = 22
		}

	}

	if domain != "" && port != 0 {
		result = domain + ":" + strconv.Itoa(port)
	}

	return result

}
