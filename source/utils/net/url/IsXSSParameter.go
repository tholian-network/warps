package url

import "strings"

var xss_values []string = []string{
	";--",
	"SELECT",
	"/../",
	"<script",
}

func IsXSSParameter(key string, value string) bool {

	var result bool = false

	for v := 0; v < len(xss_values); v++ {

		if strings.Contains(value, xss_values[v]) {
			result = true
			break
		}

	}

	return result

}
