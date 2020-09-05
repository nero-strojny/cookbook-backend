package controller

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const dbName = "tastyBoiDatabase"

// RecipeCollection mongo db object to connect to the collection housing the recipe data
var RecipeCollection *mongo.Collection

// UserCollection mongo db object to connect to the collection housing the user data
var UserCollection *mongo.Collection

// CalorieLogCollection mongo db object to connect to the collection housing the calorie log data
var CalorieLogCollection *mongo.Collection

// SetClients connect to mongo db, set up the collections
func SetClients(dbString string, env string) {
	log.Print(dbString, env)
	clientOptions := options.Client().ApplyURI(dbString)
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
		CalorieLogCollection = client.Database(dbName).Collection("testCalorieLogCollection")
		UserCollection = client.Database(dbName).Collection("testUserCollection")
	} else {
		RecipeCollection = client.Database(dbName).Collection("cookbookCollection")
		CalorieLogCollection = client.Database(dbName).Collection("calorieLogCollection")
		UserCollection = client.Database(dbName).Collection("userCollection")

	}
	fmt.Println("Collection instances created!")
}
