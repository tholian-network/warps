package strings

import "strings"

func Split(line string, separator string) []string {

	if separator == "" {
		separator = " "
	}

	if separator == " " {
		line = strings.ReplaceAll(line, "\t", " ")
	}

	result := make([]string, 0)
	source := strings.Split(strings.TrimSpace(line), separator)

	for s := 0; s < len(source); s++ {

		chunk := strings.TrimSpace(source[s])

		if chunk != "" {
			result = append(result, chunk)
		}

	}

	return result

}
