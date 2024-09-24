package http

import "strconv"
import "strings"

func ParseContentRange(value string) (int, int, int) {

	var from   int = -1
	var to     int = -1
	var length int = -1

	if strings.HasPrefix(value, "bytes ") && strings.Contains(value, "-") && strings.Contains(value, "/") {

		// "bytes 0-511/512"
		tmp := strings.Split(strings.TrimSpace(value[6:]), "/")

		if len(tmp) == 2 && strings.Contains(tmp[0], "-") {

			num_from,   err_from   := strconv.ParseInt(tmp[0][0:strings.Index(tmp[0], "-")], 10, 64)
			num_to,     err_to     := strconv.ParseInt(tmp[0][strings.Index(tmp[0], "-")+1:], 10, 64)
			num_length, err_length := strconv.ParseInt(tmp[1], 10, 64)

			if err_from == nil && err_to == nil && err_length == nil {
				from   = int(num_from)
				to     = int(num_to)
				length = int(num_length)
			}

		}

	}

	return from, to, length

}
