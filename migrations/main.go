package main

import (
	"os"

	migrations "github.com/robinjoseph08/go-pg-migrations"
)

const directory = "migrations"

func main() {
	config := GetConfig()
	db := DB(config)

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		return
	}
}
