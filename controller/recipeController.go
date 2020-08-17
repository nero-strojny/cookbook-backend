package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const recipeDBName = "recipesTable"
const recipeCollectionName = "cookbookCollection"

// recipeCollection object/instance
var recipeCollection *mongo.Collection
var recipeClient *mongo.Client

//SetRecipeClient
func SetRecipeClient(c *mongo.Client) {
	recipeClient = c
	recipeCollection = recipeClient.Database(recipeDBName).Collection(recipeCollectionName)

	fmt.Println("Collection instance created!")
}

//GetAllRecipes - gets all recipes
func GetAllRecipes() ([]primitive.M, error) {
	var emptyResults []primitive.M
	cur, err := recipeCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return emptyResults, err
	}

	// individually decode mongo results
	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			return emptyResults, e
		}
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		return emptyResults, err
	}

	cur.Close(context.Background())
	return results, nil
}

//GetRecipe - gets recipes by its ID
func GetRecipe(recipeID string) (models.Recipe, error) {
	result := models.Recipe{}
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	err := recipeCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//SearchRecipe - searches for a recipe by exact name
func SearchRecipe(name string) (models.Recipe, error) {
	result := models.Recipe{}
	filter := bson.M{"recipename": name}
	log.Print(name)
	err := recipeCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil

}

//DeleteRecipe - deletes a recipe by its ID.
func DeleteRecipe(recipeID string) error {
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	result, err := recipeCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}

//CreateRecipe a new recipe
func CreateRecipe(recipe models.Recipe) (models.Recipe, []string, error) {
	currentTime := time.Now()
	recipe.CreatedDate = currentTime.Format("2006.01.02 15:04:05")
	recipe.LastUpdatedDate = currentTime.Format("2006.01.02 15:04:05")
	valid, invalidFields := isValidRecipe(recipe)
	if valid == false {
		return models.Recipe{}, invalidFields, errors.New("Invalid fields")
	}
	result, err := recipeCollection.InsertOne(context.Background(), recipe)

	if err != nil {
		return models.Recipe{}, invalidFields, err
	}

	recipe.RecipeID = result.InsertedID.(primitive.ObjectID)
	return recipe, invalidFields, nil
}

//UpdateRecipe - updates an existing recipe by its id
func UpdateRecipe(recipeID string, updatedRecipe models.Recipe) (models.Recipe, error) {
	currentTime := time.Now()
	updatedRecipe.LastUpdatedDate = currentTime.Format("2006.01.02 15:04:05")
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	//Could do this as an update but that requires checking what fields are different between recipes
	//Could be a hassle with a long list of ingredients or measurements. Easier to just replace the entire recipe with the new update
	opts := options.Replace().SetUpsert(true)
	result, err := recipeCollection.ReplaceOne(context.Background(), filter, updatedRecipe, opts)

	if err != nil {
		return models.Recipe{}, err
	}
	if result.UpsertedID != nil {
		updatedRecipe.RecipeID = result.UpsertedID.(primitive.ObjectID)
	}
	return updatedRecipe, nil
}

func isValidRecipe(recipe models.Recipe) (valid bool, invalidFields []string) {

	if recipe.RecipeName == "" {
		invalidFields = append(invalidFields, "recipeName")
	}

	if len(invalidFields) > 0 {
		return false, invalidFields
	}
	return true, invalidFields

}
