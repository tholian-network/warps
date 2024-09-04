package url

import "tholian-warps/insights"

func IsTrackingParameter(domain string, key string, value string) bool {

	var result bool = false

	entry, ok := insights.TrackingParameters[domain]

	if ok == true {

		for p := 0; p < len(entry.Parameters); p++ {

			if key == entry.Parameters[p] {
				result = true
				break
			}

		}

	}

	if result == false {

		for _, entry := range insights.GenericTrackingParameters {

			for p := 0; p < len(entry.Parameters); p++ {

				if key == entry.Parameters[p] {
					result = true
					break
				}

			}

			if result == true {
				break
			}

		}

	}

	return result

}
