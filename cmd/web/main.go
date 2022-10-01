package main

import (
	"database/sql"
	"github.com/DeLuci/coog-music/internal/config"
	"github.com/DeLuci/coog-music/internal/driver"
	"log"
)

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {
			panic(err)
		}
	}(db.SQL)
}

func run() (*driver.DB, error) {
	log.Println("Connecting to database...")
	dbSource := config.GetDBEnvVariable("DB_SOURCE")
	db, err := driver.ConnectSQL(dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("connected to database!!!")

	return db, nil
}
