package main

import (
	"fmt"
	"sync"

	"github.com/go-pg/pg"
)

var db *pg.DB

var once sync.Once

func DB(config ...Config) *pg.DB {
	once.Do(func() {
		if len(config) > 0 {
			db = pgConnect(config[0])
		}
	})
	return db
}

func pgConnect(config Config) *pg.DB {
	var dbConfig = &pg.Options{
		User:     config["DATABASE_USER"],
		Password: config["DATABASE_PASSWORD"],
		Database: config["DATABASE_NAME"],
		Addr:     fmt.Sprintf("%s:%s", config["POSTGRES_HOST"], config["POSTGRES_PORT"]),
	}
	return pg.Connect(dbConfig)
}
