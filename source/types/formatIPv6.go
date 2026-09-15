package types

import "strings"
import utils_strings "tholian-warps/utils/strings"

func formatIPv6(value string) string {

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = value[1 : len(value)-1]
	}

	if strings.Contains(value, "::") {

		tmp := strings.Split(value, "::")
		prefix := strings.Split(tmp[0], ":")
		suffix := strings.Split(tmp[1], ":")
		complete := []string{"0000", "0000", "0000", "0000", "0000", "0000", "0000", "0000"}
		valid := true

		if len(prefix)+len(suffix) <= 8 {

			for p := 0; p < len(prefix); p++ {

				chunk := prefix[p]
				c := p

				if utils_strings.IsHex(chunk) {

					if len(chunk) == 0 {
						// Do Nothing
					} else if len(chunk) == 1 {
						complete[c] = "000" + chunk
					} else if len(chunk) == 2 {
						complete[c] = "00" + chunk
					} else if len(chunk) == 3 {
						complete[c] = "0" + chunk
					} else if len(chunk) == 4 {
						complete[c] = chunk
					} else {
						valid = false
						break
					}

				} else {
					valid = false
					break
				}

			}

			for s := 0; s < len(suffix); s++ {

				chunk := suffix[s]
				c := 8 - len(suffix) + s

				if utils_strings.IsHex(chunk) {

					if len(chunk) == 0 {
						// Do Nothing
					} else if len(chunk) == 1 {
						complete[c] = "000" + chunk
					} else if len(chunk) == 2 {
						complete[c] = "00" + chunk
					} else if len(chunk) == 3 {
						complete[c] = "0" + chunk
					} else if len(chunk) == 4 {
						complete[c] = chunk
					} else {
						valid = false
						break
					}

				} else {
					valid = false
					break
				}

			}

		} else {
			valid = false
		}

		if valid == true && len(complete) == 8 {
			return "[" + strings.Join(complete, ":") + "]"
		}

	} else if strings.Contains(value, ":") {

		chunks := strings.Split(value, ":")
		complete := []string{"0000", "0000", "0000", "0000", "0000", "0000", "0000", "0000"}
		valid := true

		if len(chunks) == 8 {

			for c := 0; c < len(chunks); c++ {

				chunk := chunks[c]

				if utils_strings.IsHex(chunk) {

					if len(chunk) == 1 {
						complete[c] = "000" + chunk
					} else if len(chunk) == 2 {
						complete[c] = "00" + chunk
					} else if len(chunk) == 3 {
						complete[c] = "0" + chunk
					} else if len(chunk) == 4 {
						complete[c] = chunk
					} else {
						valid = false
						break
					}

				} else {
					valid = false
					break
				}

			}

		} else {
			valid = false
		}

		if valid == true && len(complete) == 8 {
			return "[" + strings.Join(complete, ":") + "]"
		}

	}

	return ""

}
