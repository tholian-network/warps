package console

func Enable(feature int) {

	if feature == FeatureAll {

		for feature := range features {

			if feature != FeatureAll {
				features[feature] = true
			}

		}

	} else {

		_, ok := features[feature]

		if ok == true {
			features[feature] = true
		}

	}

}
