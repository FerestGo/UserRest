package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config map[string]string

func GetConfig() Config {
	configFile := flag.String("env", ".env", "environment file path")
	flag.Parse()
	log.Println("CONFIG", *configFile)

	config, err := godotenv.Read(*configFile)
	if err != nil {
		log.Println(err)
		config = make(Config)
		environ := os.Environ()
		for _, variable := range environ {
			varName := variable[0:strings.Index(variable, "=")]
			varValue := variable[strings.Index(variable, "=")+1:]
			config[varName] = varValue
		}
	}

	return config
}
