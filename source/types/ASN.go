package types

import "strconv"
import "strings"

type ASN int

func IsASN(value string) bool {

	var result bool

	if strings.HasPrefix(value, "AS") && len(value) >= 3 {

		result = true

		for v := 2; v < len(value); v++ {

			chr := string(value[v])

			if chr >= "0" && chr <= "9" {
				continue
			} else {
				result = false
				break
			}

		}

	}

	return result

}

func ParseASN(value string) *ASN {

	var result *ASN = nil

	if strings.HasPrefix(value, "AS") {

		num, err := strconv.ParseInt(value[2:], 10, 64)

		if err == nil {
			asn := ASN(num)
			result = &asn
		}

	} else {

		num, err := strconv.ParseInt(value, 10, 64)

		if err == nil {
			asn := ASN(num)
			result = &asn
		}

	}

	return result

}

func ToASN(value string) ASN {

	var asn ASN = ASN(0)

	if value != "" {

		tmp := ParseASN(value)

		if tmp != nil {
			asn = *tmp
		}

	}

	return asn

}

func (asn ASN) String() string {
	return "AS" + strconv.Itoa(int(asn))
}

func (asn *ASN) IsValid() bool {

	var result bool

	if IsASN(asn.String()) {
		result = true
	}

	return result

}

func (asn ASN) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("AS" + strconv.Itoa(int(asn)))), nil
}

func (asn *ASN) UnmarshalJSON(data []byte) error {

	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	tmp := ParseASN(unquoted)

	if tmp != nil {
		*asn = *tmp
	}

	return nil

}
