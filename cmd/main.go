package main

import (
	"log"

	"github.com/sugyk/rest_vpn/lib/configs"
	"github.com/sugyk/rest_vpn/lib/database"
)

func main() {
	// This is the entry point of the application
	// There is parsing of configs, connecting to database, and starting the server

	// Parse the db config file
	db_cnf, err := configs.NewDbConfigs()

	if err != nil {
		log.Fatalln("Error parsing db config file")
	}

	log.Println("Successfully parsed db config file")

	_, err = database.NewPostgresConnection(database.ConfingDB(db_cnf))

	if err != nil {
		log.Fatalln("Error connecting to database")
	}

	log.Println("Successfully connected to database")
}
