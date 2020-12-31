package traffic

// ConfigFileSchema Lists of positions w/confidence through which vehicle has passed
type ConfigFileSchema struct {
	InputFiles     InputFileConfig      `json:"input_files"`
	TrailDetectAlt TrailDetectAltConfig `json:"traildetect_alt"` // ID given to the vehicle
	OutputDirs     OutputFileConfig     `json:"output_paths"`
}

// OutputFileConfig
type OutputFileConfig struct {
	TrailDetectAlt string `json:"traildetect_alt"`
	TrailDetect    string `json:"traildetect"`
}

// InputFileConfig bargain
type InputFileConfig struct {
	TrailDetectAlt string `json:"traildetect_alt"`
	TrailDetect    string `json:"traildetect"`
}

// TrailDetectAltConfig is
type TrailDetectAltConfig struct {
	Rewards            int     `json:"rewards"`
	Penalty            int     `json:"penalty"`
	XThreshold         float64 `json:"x_threshold"`
	YThreshold         float64 `json:"y_threshold"`
	EliminateThreshold int     `json:"eliminate_threshold"`
}
