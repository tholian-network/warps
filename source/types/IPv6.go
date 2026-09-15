package types

import "math"
import "strconv"
import "strings"

type IPv6 [16]byte

func IsIPv6(raw string) bool {

	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {

		value := formatIPv6(raw[1 : len(raw)-1])

		if value != "" {

			// Ignore embedded IPv6::IPv4 syntax
			if !strings.Contains(value[1:len(value)-1], ".") {

				tmp := strings.Split(value[1:len(value)-1], ":")

				if len(tmp) == 8 {

					valid := true

					for t := 0; t < len(tmp); t++ {

						_, err := strconv.ParseUint(tmp[t], 16, 64)

						if err != nil {
							valid = false
							break
						}

					}

					return valid

				}

			}

		}

	} else if strings.HasPrefix(raw, "[") && strings.Contains(raw, "]:") {

		// Do Nothing

	} else if strings.Contains(raw, ":") {

		value := formatIPv6(raw)

		if value != "" {

			// Ignore embedded IPv6::IPv4 syntax
			if !strings.Contains(value[1:len(value)-1], ".") {

				tmp := strings.Split(value[1:len(value)-1], ":")

				if len(tmp) == 8 {

					valid := true

					for t := 0; t < len(tmp); t++ {

						_, err := strconv.ParseUint(tmp[t], 16, 64)

						if err != nil {
							valid = false
							break
						}

					}

					return valid

				}

			}

		}

	}

	return false

}

func IsIPv6AndPort(raw string) bool {

	if strings.HasPrefix(raw, "[") && strings.Contains(raw, "]:") {

		value := formatIPv6(raw[1:strings.Index(raw, "]:")])

		if value != "" {

			tmp := raw[strings.Index(raw, "]:")+2:]
			_, err := strconv.ParseUint(tmp, 10, 16)

			if IsIPv6(value) && err == nil {
				return true
			}

		}

	}

	return false

}

func ParseIPv6(raw string) *IPv6 {

	var result *IPv6 = nil

	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {

		value := formatIPv6(raw[1 : len(raw)-1])

		if value != "" {

			tmp1 := strings.Join(strings.Split(value[1:len(value)-1], ":"), "")
			tmp2 := make([]uint8, 0)

			for t := 0; t < len(tmp1); t += 2 {

				hex := string(tmp1[t : t+2])
				num, err := strconv.ParseUint(hex, 16, 8)

				if err == nil {
					tmp2 = append(tmp2, uint8(num))
				} else {
					tmp2 = append(tmp2, uint8(0))
				}

			}

			if len(tmp2) == 16 {
				ipv6 := IPv6(tmp2)
				result = &ipv6
			}

		}

	} else if strings.HasPrefix(raw, "[") && strings.Contains(raw, "]:") {

		// Do Nothing

	} else if strings.Contains(raw, ":") {

		value := formatIPv6(raw)

		if value != "" {

			tmp1 := strings.Join(strings.Split(value[1:len(value)-1], ":"), "")
			tmp2 := make([]uint8, 0)

			for t := 0; t < len(tmp1); t += 2 {

				hex := string(tmp1[t : t+2])
				num, err := strconv.ParseUint(hex, 16, 8)

				if err == nil {
					tmp2 = append(tmp2, uint8(num))
				} else {
					tmp2 = append(tmp2, uint8(0))
				}

			}

			if len(tmp2) == 16 {
				ipv6 := IPv6(tmp2)
				result = &ipv6
			}

		}

	}

	return result

}

func ParseIPv6AndPort(raw string) (*IPv6, uint16) {

	var result_ipv6 *IPv6 = nil
	var result_port uint16 = 0

	if strings.HasPrefix(raw, "[") && strings.Contains(raw, "]:") {

		value := formatIPv6(raw[1:strings.Index(raw, "]:")])

		if value != "" {

			tmp1 := ParseIPv6(value)
			tmp2 := raw[strings.Index(raw, "]:")+2:]
			num, err := strconv.ParseUint(tmp2, 10, 16)

			if tmp1 != nil && err == nil {
				result_ipv6 = tmp1
				result_port = uint16(num)
			}

		}

	}

	return result_ipv6, result_port

}

func (ipv6 IPv6) Bytes(prefix uint8) []byte {

	bytes := make([]byte, 16)

	for i := 0; i < len(ipv6); i++ {
		bytes[i] = byte(ipv6[i])
	}

	if prefix%8 == 0 {

		bytes_length := int(prefix / 8)
		bytes_mask := make([]byte, bytes_length)
		copy(bytes_mask, bytes)

		for len(bytes_mask) < len(bytes) {
			bytes_mask = append(bytes_mask, uint8(0))
		}

		return bytes_mask

	} else {

		bytes_length := int(math.Floor(float64(prefix / 8)))
		bytes_mask := make([]byte, bytes_length)
		copy(bytes_mask, bytes)

		var last_byte byte = bytes[bytes_length]
		remaining_bits := int(int(prefix) - bytes_length*8)

		switch remaining_bits {
		case 1:
			bytes_mask = append(bytes_mask, last_byte&0b10000000)
		case 2:
			bytes_mask = append(bytes_mask, last_byte&0b11000000)
		case 3:
			bytes_mask = append(bytes_mask, last_byte&0b11100000)
		case 4:
			bytes_mask = append(bytes_mask, last_byte&0b11110000)
		case 5:
			bytes_mask = append(bytes_mask, last_byte&0b11111000)
		case 6:
			bytes_mask = append(bytes_mask, last_byte&0b11111100)
		case 7:
			bytes_mask = append(bytes_mask, last_byte&0b11111110)
		case 8:
			bytes_mask = append(bytes_mask, last_byte&0b11111111)
		}

		for len(bytes_mask) < len(bytes) {
			bytes_mask = append(bytes_mask, uint8(0))
		}

		return bytes_mask

	}

}

func (ipv6 IPv6) Scope() string {

	var result string = "public"

	private_ipv6s := []string{
		// RFC3513
		"0000:0000:0000:0000:0000:0000:0000:0000",
		"0000:0000:0000:0000:0000:0000:0000:0001",
		"fe80:0000:0000:0000",
	}
	value := ipv6.String()

	for p := 0; p < len(private_ipv6s); p++ {

		if strings.HasPrefix(value, private_ipv6s[p]) {
			result = "private"
			break
		}

	}

	return result

}

func (ipv6 IPv6) String() string {

	tmp := make([]string, 8)

	for i := 0; i < len(ipv6); i += 2 {

		hex1 := strconv.FormatUint(uint64(ipv6[i]), 16)
		hex2 := strconv.FormatUint(uint64(ipv6[i+1]), 16)

		if len(hex1) == 1 {
			hex1 = "0" + hex1
		}

		if len(hex2) == 1 {
			hex2 = "0" + hex2
		}

		tmp[i/2] = hex1 + hex2

	}

	return strings.Join(tmp, ":")

}

func (ipv6 IPv6) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(ipv6.String())), nil
}

func (ipv6 *IPv6) UnmarshalJSON(data []byte) error {

	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	tmp := ParseIPv6(unquoted)

	if tmp != nil {
		*ipv6 = *tmp
	}

	return nil

}

func (ipv6 *IPv6) IsValid() bool {

	var result bool

	if IsIPv6(ipv6.String()) {
		result = true
	}

	return result

}
