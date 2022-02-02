package test

import (
	"context"
	"flag"
	"log"
	"server/controller"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbPointer = flag.String("DB_STRING", "", "Database connection string")
var envPointer = flag.String("ENV", "", "Environment string")

func TestDatabaseSetup(t *testing.T) {
	clientOptions := options.Client().ApplyURI(*dbPointer)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	controller.SetClients(mongoClient)
	controller.GetCollections(*envPointer)
}
