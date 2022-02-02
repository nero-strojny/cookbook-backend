package controller

import (
	"context"
	"fmt"
	"log"
	"server/models"

	"go.mongodb.org/mongo-driver/mongo"
)

const dbName = "tastyBoiDatabase"

// RecipeCollection mongo db object to connect to the collection housing the recipe data
var RecipeCollection *mongo.Collection

// UserCollection mongo db object to connect to the collection housing the user data
var UserCollection *mongo.Collection

// IngredientCollection mongo db object to connect to the collection housing the user data
var IngredientCollection *mongo.Collection

// HouseholdCollection mongo db object to connect to the collection housing the user data
var HouseholdCollection *mongo.Collection

var client *mongo.Client

// SetClients connect to mongo db, set up the collections
func SetClients(mongoClient *mongo.Client) {
	client = mongoClient
	// Check the connection
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func GetCollections(env string) {
	if env == "dev" {
		RecipeCollection = client.Database(dbName).Collection("testCookbookCollection")
		UserCollection = client.Database(dbName).Collection("testUserCollection")
		IngredientCollection = client.Database(dbName).Collection("testIngredientCollection")
		HouseholdCollection = client.Database(dbName).Collection("testHouseholdCollectionCollection")
	} else {
		RecipeCollection = client.Database(dbName).Collection("cookbookCollection")
		UserCollection = client.Database(dbName).Collection("userCollection")
		IngredientCollection = client.Database(dbName).Collection("ingredientCollection")
		HouseholdCollection = client.Database(dbName).Collection("householdCollectionCollection")

	}
	fmt.Println("Collection instances created!")
}

func HealthCheck() models.HealthStatus {
	mongoStatus := "OK"
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		mongoStatus = fmt.Sprintf("Mongo ERROR: %s", err)
	}
	return models.HealthStatus{DB: mongoStatus}
}
