package controller

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//DeleteRecipe - deletes a recipe by its ID.
func DeleteRecipe(recipeID string) error {
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	result, err := RecipeCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}
