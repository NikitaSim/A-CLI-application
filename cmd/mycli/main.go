package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mycli/internal/app"
	"mycli/internal/config"
	"mycli/internal/db"
)

var flags struct {
	ConfigPath    string
	DirectoryPath string
}

func init() {
	flag.StringVar(&flags.ConfigPath, "c", "", "path to yaml schema config files")
	flag.StringVar(&flags.DirectoryPath, "d", "", "path to directory with files")
}

const postgresqlConnString = "postgres://user:password@localhost:5433/postgres"

func main() {
	// Инициализация 5 аргументов
	if len(os.Args) != 5 {
		log.Fatalln("not enough arguments")
	} else {
		flag.Parse()
	}

	//fmt.Println(flags.ConfigPath)
	configsPath := app.ProcessConfigs(flags.DirectoryPath)
	//fmt.Println(configsPath)

	configs := config.ParseConfig(configsPath, flags.ConfigPath)

	// for _, value := range configs {
	// 	fmt.Println(value)
	// }
	err := db.DbInsert(configs, postgresqlConnString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PASS")
}
