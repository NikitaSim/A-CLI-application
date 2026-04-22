package models

// Config - общая модель файлов (конфигов)
type Config struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Version     int      `yaml:"version"`
	Metadata    Metadata `yaml:"metadata"`
}

type Metadata struct {
	Author string   `yaml:"author"`
	Tags   []string `yaml:"tags"`
}

// Schema - модель описания JSON схемы (поле->путь)
type Schema struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Tags        string `yaml:"tags"`
}

type JSONSchemas struct {
	Schemas []Schema `yaml:"json_schemas"`
}
