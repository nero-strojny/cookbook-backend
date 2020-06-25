package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "recipesTable"
const collName = "cookbookCollection"

// collection object/instance
var collection *mongo.Collection

func openFile() string {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return "Could not open file"
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	if str, ok := result["connectionString"].(string); ok {
		return str
	}

	return "Could not find db connection string"
}

// create connection with mongo db
func init() {
	dbConnectionString := openFile()
	clientOptions := options.Client().ApplyURI(dbConnectionString)
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

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created!")
}

//GetAll - gets all recipes
func GetAll() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

//Get - gets recipes by its ID
func Get(recipeID string) models.Recipe {
	result := models.Recipe{}
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

//Delete a recipe by its ID.
func Delete(recipeID string) {
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
}

//Create a new recipe
func Create(recipe models.Recipe) (invalidData []string, errorCode int) {
	currentTime := time.Now()
	recipe.CreatedDate = currentTime.Format("2006.01.02 15:04:05")
	recipe.LastUpdatedDate = currentTime.Format("2006.01.02 15:04:05")
	valid, invalidFields := validateRecipe(recipe)
	if valid == false {
		return invalidFields, 400
	}
	_, err := collection.InsertOne(context.Background(), recipe)

	if err != nil {
		log.Fatal(err)
	}
	return invalidFields, 204
}

//Update an existing recipe by its id
func Update(recipeID string, updatedRecipe models.Recipe) {
	currentTime := time.Now()
	updatedRecipe.LastUpdatedDate = currentTime.Format("2006.01.02 15:04:05")
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	//Could do this as an update but that requires checking what fields are different between recipes
	//Could be a hassle with a long list of ingredients or measurements. Easier to just replace the entire recipe with the new update
	_, err := collection.ReplaceOne(context.Background(), filter, updatedRecipe)

	if err != nil {
		log.Fatal(err)
	}
}

func validateRecipe(recipe models.Recipe) (valid bool, invalidFields []string) {

	if recipe.RecipeName == "" {
		invalidFields = append(invalidFields, "recipeName")
	}

	if len(invalidFields) > 0 {
		return false, invalidFields
	} 
		return true, invalidFields
	
}
