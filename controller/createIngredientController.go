package controller

import (
	"context"

	"server/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateIngredient creates a new ingredient
func CreateIngredient(ingredient models.Ingredient) (models.Ingredient, error) {
	result, err := IngredientCollection.InsertOne(context.Background(), ingredient)

	if err != nil {
		return models.Ingredient{}, err
	}

	ingredient.IngredientID = result.InsertedID.(primitive.ObjectID)
	return ingredient, nil
}
