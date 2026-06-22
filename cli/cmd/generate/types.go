package generate

type plateYAML struct {
	Name    string                 `yaml:"name" json:"Name"`
	Schema  map[string]schemaField `yaml:"schema" json:"Schema"`
	Modules map[string]moduleDef   `yaml:"modules" json:"Modules"`
	Files   []fileEntry            `yaml:"files" json:"Files"`
}

type schemaField struct {
	Type     string   `yaml:"type" json:"Type"`
	Required bool     `yaml:"required" json:"Required"`
	Values   []string `yaml:"values" json:"Values"`
	Default  any      `yaml:"default" json:"Default"`
}

type moduleDef struct {
	Enabled bool `yaml:"enabled" json:"Enabled"`
}

type fileEntry struct {
	Path      string `yaml:"path"`
	Template  string `yaml:"template"`
	Condition string `yaml:"condition"`
}
