package types

import "strconv"
import "strings"

type Prefix uint8

func IsPrefix(value string) bool {

	var result bool

	if strings.HasPrefix(value, "/") {

		_, err := strconv.ParseUint(value[1:], 10, 8)

		if err == nil {
			result = true
		}

	} else if strings.HasPrefix(value, "[") && strings.Contains(value, "]/") {

		_, err := strconv.ParseUint(strings.Split(value, "]/")[1], 10, 8)

		if err == nil {
			result = true
		}

	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {

		result = true

	} else if strings.Contains(value, ":") && strings.Contains(value, "/") {

		_, err := strconv.ParseUint(strings.Split(value, "/")[1], 10, 8)

		if err == nil {
			result = true
		}

	} else if strings.Contains(value, ":") {

		result = true

	} else if strings.Contains(value, ".") && strings.Contains(value, "/") {

		_, err := strconv.ParseUint(strings.Split(value, "/")[1], 10, 8)

		if err == nil {
			result = true
		}

	} else if strings.Contains(value, ".") {

		result = true

	} else {

		_, err := strconv.ParseUint(value, 10, 8)

		if err == nil {
			result = true
		}

	}

	return result

}

func ParsePrefix(value string) *Prefix {

	var result *Prefix = nil

	if strings.HasPrefix(value, "/") {

		tmp, err := strconv.ParseUint(value[1:], 10, 8)

		if err == nil {
			prefix := Prefix(tmp)
			result = &prefix
		}

	} else if strings.HasPrefix(value, "[") && strings.Contains(value, "]/") {

		tmp, err := strconv.ParseUint(strings.Split(value, "]/")[1], 10, 8)

		if err == nil {
			prefix := Prefix(tmp)
			result = &prefix
		}

	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {

		prefix := Prefix(128)
		result = &prefix

	} else if strings.Contains(value, ":") && strings.Contains(value, "/") {

		tmp, err := strconv.ParseUint(strings.Split(value, "/")[1], 10, 8)

		if err == nil {
			prefix := Prefix(tmp)
			result = &prefix
		}

	} else if strings.Contains(value, ":") {

		prefix := Prefix(128)
		result = &prefix

	} else if strings.Contains(value, ".") && strings.Contains(value, "/") {

		tmp, err := strconv.ParseUint(strings.Split(value, "/")[1], 10, 8)

		if err == nil {
			prefix := Prefix(tmp)
			result = &prefix
		}

	} else if strings.Contains(value, ".") {

		prefix := Prefix(32)
		result = &prefix

	} else {

		tmp, err := strconv.ParseUint(value, 10, 8)

		if err == nil {
			prefix := Prefix(tmp)
			result = &prefix
		}

	}

	return result

}

func ToPrefix(value string) Prefix {

	var prefix Prefix = Prefix(0)

	if value != "" {

		tmp := ParsePrefix(value)

		if tmp != nil {
			prefix = *tmp
		}

	}

	return prefix

}

func (prefix Prefix) String() string {
	return "/" + strconv.FormatUint(uint64(prefix), 10)
}

func (prefix *Prefix) IsValid() bool {

	var result bool

	if IsPrefix(prefix.String()) {
		result = true
	}

	return result

}

func (prefix Prefix) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("/" + strconv.FormatUint(uint64(prefix), 10))), nil
}

func (prefix *Prefix) UnmarshalJSON(data []byte) error {

	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	tmp := ParsePrefix(unquoted)

	if tmp != nil {
		*prefix = *tmp
	}

	return nil

}
