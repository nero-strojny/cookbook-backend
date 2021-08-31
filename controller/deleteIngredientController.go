package controller

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//DeleteIngredient - deletes a ingredient by its ID.
func DeleteIngredient(ingredientID string) error {
	id, _ := primitive.ObjectIDFromHex(ingredientID)
	filter := bson.M{"_id": id}
	result, err := IngredientCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}
