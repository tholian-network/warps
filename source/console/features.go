package console

const (
	FeatureAll int = iota
	FeatureGroup
	FeatureLog
	FeatureInfo
	FeatureWarn
	FeatureError
	FeatureInspect
	FeatureProgress
)

var features map[int]bool

func init() {

	features = make(map[int]bool)

	features[FeatureGroup] = true
	features[FeatureLog] = true
	features[FeatureInfo] = true
	features[FeatureWarn] = true
	features[FeatureError] = true
	features[FeatureInspect] = true
	features[FeatureProgress] = true

}
