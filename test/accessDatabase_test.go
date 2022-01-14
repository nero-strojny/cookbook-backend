package test

import (
	"context"
	"flag"
	"log"
	"server/config"
	"server/controller"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbPointer = flag.String("DB_STRING", "", "Database connection string")

func TestDatabaseSetup(t *testing.T) {
	// If the dBString is empty, then we need to fall back on a file if one is present
	var dBFlag string
	dBFlag = *dbPointer
	if dBFlag == "" {
		dBFlag = config.GetConfig("../config.json").ConnectionString
	}
	clientOptions := options.Client().ApplyURI(dBFlag)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	controller.SetClients(mongoClient)
	controller.GetCollections("dev")
}
