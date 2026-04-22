package config

// Логика парсинга и получения конфига

import (
	"encoding/json"
	"fmt"
	"mycli/internal/models"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func parseWithSchemas(data map[string]interface{}, schemas []models.Schema) (models.Config, bool) {

	for _, schema := range schemas {

		var cfg models.Config
		flag := true

		if value, ok := getByPath(data, schema.Name); ok {
			name, ok := value.(string)
			if !ok {
				flag = false
			}
			cfg.Name = name
		} else {
			continue
		}

		if value, ok := getByPath(data, schema.Author); ok {
			author, ok := value.(string)
			if !ok {
				flag = false

			}
			cfg.Metadata.Author = author
		} else {
			continue
		}

		if value, ok := getByPath(data, schema.Version); ok {
			version, ok := value.(float64)
			if !ok {
				flag = false

			}
			cfg.Version = int(version)
		} else {
			continue
		}

		if value, ok := getByPath(data, schema.Description); ok {
			description, ok := value.(string)
			if !ok {
				flag = false

			}
			cfg.Description = description
		} else {
			continue
		}

		if value, ok := getByPath(data, schema.Tags); ok {
			rawTags, ok := value.([]interface{})
			if !ok {
				flag = false
				continue
			}
			tags := make([]string, 0, len(rawTags))
			for _, rawTag := range rawTags {
				tag, ok := rawTag.(string)
				if !ok {
					flag = false
					continue
				}
				tags = append(tags, tag)
			}
			cfg.Metadata.Tags = tags
		} else {
			continue
		}

		if flag {
			return cfg, true
		}
	}

	return models.Config{}, false
}

func getByPath(data map[string]interface{}, path string) (interface{}, bool) {

	splitPath := strings.Split(path, ".")
	current := data
	for id, value := range splitPath {
		result, exist := current[value]
		if !exist {
			return nil, false
		}

		if id < len(splitPath)-1 {
			next, ok := result.(map[string]interface{})
			if !ok {
				return nil, false
			}
			current = next
		} else {
			return result, true
		}

	}
	return nil, false
}

func ParseConfig(paths []string, yamlPath string) []models.Config {

	var schema models.JSONSchemas

	jsonPathFile, err := os.Open(yamlPath)
	if err != nil {
		panic(err)
	}
	defer jsonPathFile.Close()

	err = yaml.NewDecoder(jsonPathFile).Decode(&schema)
	if err != nil {
		panic(err)
	}

	configs := make([]models.Config, 0)

	for _, path := range paths { // можно распараллелить черех воркеров

		if filepath.Ext(path) == ".json" {
			file, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			conf := make(map[string]interface{})

			err = json.NewDecoder(file).Decode(&conf)
			if err != nil {
				panic(err)

			}

			value, ok := parseWithSchemas(conf, schema.Schemas)
			if !ok {
				fmt.Println("wrong path")
				return nil
			}

			configs = append(configs, value)

		} else {
			file, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			var yamlConfig models.Config

			err = yaml.NewDecoder(file).Decode(&yamlConfig)
			if err != nil {
				panic(err)
			}

			configs = append(configs, yamlConfig)
		}

	}
	return configs
}
