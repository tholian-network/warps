package strings

import "sort"

func Unique(values []string) []string {

	result := make([]string, 0)
	hashmap := make(map[string]bool)

	for _, value := range values {

		_, ok := hashmap[value]

		if ok == false {
			result = append(result, value)
			hashmap[value] = true
		}

	}

	sort.Strings(result)

	return result

}
