package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"server/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/controller"
	"server/router"
)

func main() {
	var dBFlag string
	var envFlag string

	flag.StringVar(&dBFlag, "tastyboi-server-DB_STRING", "", "Database connection string")
	flag.StringVar(&envFlag, "ENV", "", "Environment string")
	flag.Parse()

	// If the dBString is empty, then we need to fall back on a file if one is present
	if dBFlag == "" {
		fmt.Println("Environment variable not found")
		dBFlag = config.GetConfig().ConnectionString
	} else {
		fmt.Println("Environment variable found")
	}
	clientOptions := options.Client().ApplyURI(dBFlag)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	controller.SetClients(mongoClient)
	controller.GetCollections(envFlag)
	r := router.Router()
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
