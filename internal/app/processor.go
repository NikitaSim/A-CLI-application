package app

// Логика обработки директории с файлами (должны получать файлы из папки) пуская это будет массив с путями
import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func ProcessConfigs(pathToDir string) []string {
	configs := make([]string, 0)
	err := filepath.WalkDir(pathToDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			//fmt.Println("file:", path)
			if filepath.Ext(path) == ".json" || filepath.Ext(path) == ".yaml" {
				configs = append(configs, path)
			}

		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return configs
}
