package strings

func IsHex(value string) bool {

	var result bool = true

	for v := 0; v < len(value); v++ {

		character := string(value[v])

		if character >= "0" && character <= "9" {
			continue
		} else if character >= "A" && character <= "F" {
			continue
		} else if character >= "a" && character <= "f" {
			continue
		} else {
			result = false
			break
		}

	}

	return result

}
