package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"server/controller"
	"server/router"
)

func main() {
	var dBFlag string
	var envFlag string

	flag.StringVar(&dBFlag, "DB_STRING", "", "Database connection string")
	flag.StringVar(&envFlag, "ENV", "", "Environment string")
	flag.Parse()

	r := router.Router()
	controller.SetClients(dBFlag, envFlag)
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
