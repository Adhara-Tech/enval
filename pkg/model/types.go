package model

type Manifest struct {
	Name        string         `yaml:"name"`
	CustomSpecs string         `yaml:"custom_specs"`
	Tools       []ManifestTool `yaml:"tools"`
}

type VersionParser struct {
	Type        string
	Regexp      *string
	CodeSnippet *string
}

type ManifestTool struct {
	Name   string            `yaml:"name"`
	Flavor *string           `yaml:"flavor"`
	Checks map[string]string `yaml:"checks"`
}

func (mt ManifestTool) IsFlavoredCheck() bool {
	return mt.Flavor != nil
}
