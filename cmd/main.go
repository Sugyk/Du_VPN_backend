package main

import (
	"log"

	rest "github.com/sugyk/rest_vpn/api"
)

func main() {
	// This is the entry point of the application
	// There is parsing of configs, connecting to database, and starting the server

	log.Println("Starting the application...")
	serv := rest.NewServer()

	serv.Run("8080")
}
