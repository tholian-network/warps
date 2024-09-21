package test

import "tholian-endpoint/protocols/dns"
import "tholian-endpoint/types"
import "slices"
import "strings"

type SpyResolver struct {
	Resolved    string `json:"resolved"`
	WasResolved bool   `json:"was_resolved"`
	isOnline    bool
}

func NewSpyResolver(isOnline bool) SpyResolver {

	var resolver SpyResolver

	resolver.isOnline = isOnline

	return resolver

}

func (resolver *SpyResolver) Resolve(subject string) dns.Packet {

	var response dns.Packet

	resolved := ""

	if types.IsIPv4(subject) {

		ipv4 := types.ParseIPv4(subject)

		if ipv4 != nil {

			tmp := strings.Split(ipv4.String(), ".")
			slices.Reverse(tmp)

			resolved = "PTR:" + strings.Join(tmp, ".") + ".in-addr.arpa"

		}

	} else if types.IsIPv6(subject) {

		ipv6 := types.ParseIPv6(subject)

		if ipv6 != nil {

			tmp := strings.Split(strings.Join(strings.Split(ipv6.String(), ":"), ""), "")
			slices.Reverse(tmp)

			resolved = "PTR:" + strings.Join(tmp, ".") + "ip6.arpa"

		}

	} else if types.IsDomain(subject) {

		resolved = "A:" + subject + ",AAAA:" + subject

	}

	if resolver.isOnline == true {
		response = dns.Resolve(subject)
	}

	resolver.Resolved = resolved
	resolver.WasResolved = true

	return response

}

func (resolver *SpyResolver) ResolvePacket(query dns.Packet) dns.Packet {

	var response dns.Packet

	if query.Type == "query" && len(query.Questions) > 0 {

		resolved := ""

		for q := 0; q < len(query.Questions); q++ {

			resolved += query.Questions[q].Type.String() + ":" + query.Questions[q].Name

			if q <= len(query.Questions) - 1 {
				resolved += ","
			}

		}

		if resolver.isOnline == true {
			response = dns.ResolvePacket(query)
		}

		resolver.WasResolved = true
		resolver.Resolved = resolved

	}

	return response

}
