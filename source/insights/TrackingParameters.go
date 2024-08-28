package insights

import _ "embed"
import "encoding/json"

//go:embed TrackingParameters.json
var embedded_tracking_parameters []byte

var TrackingParameters map[string]*TrackingParameter

func init() {

	TrackingParameters = make(map[string]*TrackingParameter)

	list := make([]TrackingParameter, 0)

	err := json.Unmarshal(embedded_tracking_parameters, &list)

	if err == nil {

		for l := 0; l < len(list); l++ {

			entry := list[l]

			for d := 0; d < len(entry.Domains); d++ {
				TrackingParameters[entry.Domains[d]] = &entry
			}

		}

	}

}
