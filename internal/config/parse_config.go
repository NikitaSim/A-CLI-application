package config

// Логика парсинга и получения конфига

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"mycli/internal/models"

	"gopkg.in/yaml.v3"
)

const workerCount = 5

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

func worker(jobs <-chan string, wg *sync.WaitGroup, result chan<- models.Config, schema models.JSONSchemas) {
	defer wg.Done()

	for path := range jobs { // можно распараллелить черех воркеров

		if filepath.Ext(path) == ".json" {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			conf := make(map[string]interface{})

			err = json.NewDecoder(file).Decode(&conf)
			if err != nil {
				file.Close()
				fmt.Println(err)
				continue
			}

			value, ok := parseWithSchemas(conf, schema.Schemas)
			if !ok {
				fmt.Println("wrong path")
				continue
			}
			result <- value

			file.Close()

		} else {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			var yamlConfig models.Config

			err = yaml.NewDecoder(file).Decode(&yamlConfig)
			if err != nil {
				file.Close()
				fmt.Println(err)
				continue
			}
			result <- yamlConfig

			file.Close()
		}
	}
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

	job := make(chan string)
	result := make(chan models.Config, workerCount)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go worker(job, &wg, result, schema)
	}

	for _, path := range paths { // можно распараллелить черех воркеров
		job <- path
	}
	close(job)

	go func() {
		wg.Wait()
		close(result)
	}()

	for value := range result {
		configs = append(configs, value)
	}

	return configs
}
