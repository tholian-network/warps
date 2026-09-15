package types

import "strconv"
import "strings"

type Domain string

func isASCIILabel(value string) bool {

	var result bool = true

	for _, chr := range value {

		character := string(chr)

		if character >= "0" && character <= "9" {
			continue
		} else if character >= "A" && character <= "Z" {
			continue
		} else if character >= "a" && character <= "z" {
			continue
		} else if character == "-" || character == "_" {
			continue
		} else {
			result = false
			break
		}

	}

	return result

}

func IsDomain(value string) bool {

	if strings.Contains(value, ".") {

		labels := strings.Split(value, ".")
		valid := true

		if len(labels) >= 2 {

			for l := 0; l < len(labels)-1; l++ {

				label := labels[l]

				if isASCIILabel(label) && len(label) >= 1 {
					// Do Nothing
				} else {
					valid = false
					break
				}

			}

			tld := labels[len(labels)-1]

			if isASCIILabel(tld) && len(tld) >= 2 {
				// Do Nothing
			} else {
				valid = false
			}

		}

		return valid

	} else {

		label := value

		if isASCIILabel(label) && len(label) >= 3 {
			return true
		}

	}

	return false

}

func IsDomainAndPort(value string) bool {

	if strings.Contains(value, ":") {

		tmp := strings.Split(value, ":")

		if len(tmp) == 2 {

			_, err := strconv.ParseUint(tmp[1], 10, 16)

			if IsDomain(tmp[0]) && err == nil {
				return true
			}

		}

	}

	return false

}

func ParseDomain(value string) *Domain {

	var result *Domain = nil

	if IsDomain(value) {
		domain := Domain(value)
		result = &domain
	}

	return result

}

func ParseDomainAndPort(value string) (*Domain, uint16) {

	var result_domain *Domain = nil
	var result_port uint16 = 0

	if strings.Contains(value, ":") {

		tmp1 := strings.Split(value, ":")

		if len(tmp1) == 2 {

			tmp2 := ParseDomain(tmp1[0])
			num, err := strconv.ParseUint(tmp1[1], 10, 16)

			if tmp2 != nil && err == nil {
				result_domain = tmp2
				result_port = uint16(num)
			}

		}

	}

	return result_domain, result_port

}

func ToDomain(value string) Domain {

	var domain Domain

	if IsDomain(value) {
		domain = Domain(value)
	}

	return domain

}

func (domain Domain) Bytes() []byte {
	return []byte(domain)
}

func (domain Domain) Scope() string {

	tmp := string(domain)

	if strings.HasSuffix(tmp, ".in-addr.arpa") {

		suffixes := []string{
			// RFC 6761
			"10.in-addr.arpa",
			"16.172.in-addr.arpa",
			"17.172.in-addr.arpa",
			"18.172.in-addr.arpa",
			"19.172.in-addr.arpa",
			"20.172.in-addr.arpa",
			"21.172.in-addr.arpa",
			"22.172.in-addr.arpa",
			"23.172.in-addr.arpa",
			"24.172.in-addr.arpa",
			"25.172.in-addr.arpa",
			"26.172.in-addr.arpa",
			"27.172.in-addr.arpa",
			"28.172.in-addr.arpa",
			"29.172.in-addr.arpa",
			"30.172.in-addr.arpa",
			"31.172.in-addr.arpa",
			"168.192.in-addr.arpa",

			// RFC 6762
			"254.169.in-addr.arpa",

			// RFC 8880
			"170.0.0.192.in-addr.arpa",
			"171.0.0.192.in-addr.arpa",
		}

		scope := "public"

		for s := 0; s < len(suffixes); s++ {

			if strings.HasSuffix(tmp, suffixes[s]) {
				scope = "private"
				break
			}

		}

		return scope

	} else if strings.HasSuffix(tmp, ".ip6.arpa") {

		suffixes := []string{
			// RFC 6762
			"8.e.f.ip6.arpa",
			"9.e.f.ip6.arpa",
			"a.e.f.ip6.arpa",
			"b.e.f.ip6.arpa",
		}

		scope := "public"

		for s := 0; s < len(suffixes); s++ {

			if strings.HasSuffix(tmp, suffixes[s]) {
				scope = "private"
				break
			}

		}

		return scope

	} else if strings.Contains(tmp, ".") {

		suffixes := []string{
			// RFC 6761
			".example",
			".invalid",
			".localhost",
			".test",

			// RFC 6762
			".local",

			// RFC 8375
			".home.arpa",

			// RFC 9476
			".alt",
		}

		scope := "public"

		for s := 0; s < len(suffixes); s++ {

			if strings.HasSuffix(tmp, suffixes[s]) {
				scope = "private"
				break
			}

		}

		return scope

	} else {

		return "private"

	}

}

func (domain Domain) String() string {
	return string(domain)
}

func (domain Domain) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(domain))), nil
}

func (domain *Domain) UnmarshalJSON(data []byte) error {

	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	tmp := ParseDomain(unquoted)

	if tmp != nil {
		*domain = *tmp
	}

	return nil

}

func (domain *Domain) IsValid() bool {

	var result bool

	if IsDomain(domain.String()) {
		result = true
	}

	return result

}
