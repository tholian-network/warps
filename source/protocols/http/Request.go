package http

import "tholian-warps/types"
import net_url "net/url"
import "strconv"

func Request(raw_url string) Packet {

	var response Packet

	url, err := net_url.Parse(raw_url)

	if err == nil {

		if types.IsIPv4AndPort(url.Host) {

			ipv4, port := types.ParseIPv4AndPort(url.Host)

			if ipv4 != nil && port != 0 {

				request := NewPacket()
				request.SetURL(*url)
				request.SetMethod(MethodGet)
				request.SetHeader("Host", ipv4.String() + ":" + strconv.FormatUint(uint64(port), 10))

				request.SetServer(types.Server{
					Domain:    ipv4.String(),
					Addresses: []string{ipv4.String()},
					Port:      port,
					Protocol:  types.Protocol(url.Scheme),
					Schema:    "",
				})

				response = RequestPacket(request)

			}

		} else if types.IsIPv4(url.Host) {

			ipv4 := types.ParseIPv4(url.Host)

			if ipv4 != nil {

				port := 0

				if url.Scheme == "https" {
					port = 443
				} else if url.Scheme == "http" {
					port = 80
				}

				request := NewPacket()
				request.SetURL(*url)
				request.SetMethod(MethodGet)
				request.SetHeader("Host", ipv4.String() + ":" + strconv.FormatUint(uint64(port), 10))

				request.SetServer(types.Server{
					Domain:    ipv4.String(),
					Addresses: []string{ipv4.String()},
					Port:      uint16(port),
					Protocol:  types.Protocol(url.Scheme),
					Schema:    "",
				})

				response = RequestPacket(request)

			}

		} else if types.IsIPv6AndPort(url.Host) {

			ipv6, port := types.ParseIPv6AndPort(url.Host)

			if ipv6 != nil && port != 0 {

				request := NewPacket()
				request.SetURL(*url)
				request.SetMethod(MethodGet)
				request.SetHeader("Host", "[" + ipv6.String() + "]:" + strconv.FormatUint(uint64(port), 10))

				request.SetServer(types.Server{
					Domain:    ipv6.String(),
					Addresses: []string{ipv6.String()},
					Port:      port,
					Protocol:  types.Protocol(url.Scheme),
					Schema:    "",
				})

				response = RequestPacket(request)

			}

		} else if types.IsIPv6(url.Host) {

			ipv6 := types.ParseIPv6(url.Host)

			if ipv6 != nil {

				port := 0

				if url.Scheme == "https" {
					port = 443
				} else if url.Scheme == "http" {
					port = 80
				}

				request := NewPacket()
				request.SetURL(*url)
				request.SetMethod(MethodGet)
				request.SetHeader("Host", "[" + ipv6.String() + "]:" + strconv.FormatUint(uint64(port), 10))

				request.SetServer(types.Server{
					Domain:    ipv6.String(),
					Addresses: []string{ipv6.String()},
					Port:      uint16(port),
					Protocol:  types.Protocol(url.Scheme),
					Schema:    "",
				})

				response = RequestPacket(request)

			}

		} else if types.IsDomainAndPort(url.Host) {

			domain, port := types.ParseDomainAndPort(url.Host)

			if domain != nil && port != 0 {

				request := NewPacket()
				request.SetURL(*url)
				request.SetMethod(MethodGet)
				request.SetHeader("Host", domain.String() + ":" + strconv.FormatUint(uint64(port), 10))

				response = RequestPacket(request)

			}

		} else if types.IsDomain(url.Host) {

			domain := types.ParseDomain(url.Host)

			if domain != nil {

				port := 0

				if url.Scheme == "https" {
					port = 443
				} else if url.Scheme == "http" {
					port = 80
				}

				request := NewPacket()
				request.SetURL(*url)
				request.SetMethod(MethodGet)
				request.SetHeader("Host", domain.String() + ":" + strconv.FormatUint(uint64(port), 10))

				response = RequestPacket(request)

			}

		}

	}

	return response

}
