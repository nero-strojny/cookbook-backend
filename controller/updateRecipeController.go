package controller

import (
	"context"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateRecipe(recipeID string, updatedRecipe models.Recipe) (models.Recipe, error) {
	currentTime := time.Now()
	updatedRecipe.LastUpdatedDate = currentTime.Format("2006.01.02 15:04:05")
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	//Could do this as an update but that requires checking what fields are different between recipes
	//Could be a hassle with a long list of ingredients or measurements. Easier to just replace the entire recipe with the new update
	opts := options.Replace().SetUpsert(true)
	result, err := RecipeCollection.ReplaceOne(context.Background(), filter, updatedRecipe, opts)

	if err != nil {
		return models.Recipe{}, err
	}
	if result.UpsertedID != nil {
		updatedRecipe.RecipeID = result.UpsertedID.(primitive.ObjectID)
	}
	return updatedRecipe, nil
}
