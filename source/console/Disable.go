package console

func Disable(feature int) {

	if feature == FeatureAll {

		for feature := range features {

			if feature != FeatureAll {
				features[feature] = false
			}

		}

	} else {

		_, ok := features[feature]

		if ok == true {
			features[feature] = false
		}

	}

}
