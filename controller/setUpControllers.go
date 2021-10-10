package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"os"
	"server/models"
)

var client *mongo.Client

const dbName = "tastyBoiDatabase"

// RecipeCollection mongo db object to connect to the collection housing the recipe data
var RecipeCollection *mongo.Collection

// UserCollection mongo db object to connect to the collection housing the user data
var UserCollection *mongo.Collection

// UserCollection mongo db object to connect to the collection housing the user data
var IngredientCollection *mongo.Collection

var Config models.Config

// SetClients connect to mongo db, set up the collections
func SetClients(dbString string, env string) {
	var connectionString string

	// If the dBString is empty, then we need to fall back on a file if one is present
	if dbString == "" {
		jsonFile, _ := os.Open("../config.json")
		defer jsonFile.Close()

		jsonData, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(jsonData, &Config)
		connectionString = Config.ConnectionString
	} else {
		connectionString = dbString
	}
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	if env == "dev" {
		RecipeCollection = client.Database(dbName).Collection("testCookbookCollection")
		UserCollection = client.Database(dbName).Collection("testUserCollection")
		IngredientCollection = client.Database(dbName).Collection("testIngredientCollection")
	} else {
		RecipeCollection = client.Database(dbName).Collection("cookbookCollection")
		UserCollection = client.Database(dbName).Collection("userCollection")
		IngredientCollection = client.Database(dbName).Collection("ingredientCollection")

	}
	fmt.Println("Collection instances created!")
}
