package models

// ProjectType represents the detected type of project
type ProjectType int

const (
	ProjectTypeUnknown ProjectType = iota
	ProjectTypeWebReact
	ProjectTypeWebVue
	ProjectTypeWebAngular
	ProjectTypeWebSvelte
	ProjectTypeReactNative
	ProjectTypeFlutter
	ProjectTypeIOS
	ProjectTypeAndroid
	ProjectTypeGo
	ProjectTypeRust
)

// String returns the string representation of ProjectType
func (pt ProjectType) String() string {
	return [...]string{
		"Unknown",
		"React (Web)",
		"Vue (Web)",
		"Angular (Web)",
		"Svelte (Web)",
		"React Native",
		"Flutter",
		"iOS (Swift)",
		"Android (Kotlin/Java)",
		"Go",
		"Rust",
	}[pt]
}

// ProjectConfig holds the configuration for scanning behavior
type ProjectConfig struct {
	// Asset Discovery
	AssetPaths   []string `yaml:"asset_paths" json:"asset_paths"`
	Extensions   []string `yaml:"extensions" json:"extensions"`
	ExcludePaths []string `yaml:"exclude_paths" json:"exclude_paths"`

	// Reference Detection
	ConstantFiles  []string `yaml:"constant_files" json:"constant_files"`
	BasePathVars   []string `yaml:"base_path_vars" json:"base_path_vars"`
	CustomPatterns []string `yaml:"custom_patterns" json:"custom_patterns"`

	// Behavior
	FollowSymlinks        bool        `yaml:"follow_symlinks" json:"follow_symlinks"`
	AutoDetectProjectType bool        `yaml:"auto_detect_project_type" json:"auto_detect_project_type"`
	ProjectType           ProjectType `yaml:"project_type" json:"project_type"`

	// Performance
	MaxWorkers  int   `yaml:"max_workers" json:"max_workers"`
	MemoryLimit int64 `yaml:"memory_limit" json:"memory_limit"`

	// Output
	Verbose     bool `yaml:"verbose" json:"verbose"`
	ShowProgress bool `yaml:"show_progress" json:"show_progress"`
	ColorOutput bool `yaml:"color_output" json:"color_output"`
}
