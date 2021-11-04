package main

import (
	"context"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"server/config"

	"server/controller"
	"server/router"
)

func main() {
	var dBFlag string
	var envFlag string

	flag.StringVar(&dBFlag, "DB_STRING", "", "Database connection string")
	flag.StringVar(&envFlag, "ENV", "", "Environment string")
	flag.Parse()

	// If the dBString is empty, then we need to fall back on a file if one is present
	if dBFlag == "" {
		dBFlag = config.GetConfig().ConnectionString
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
