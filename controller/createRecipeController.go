package controller

import (
	"context"
	"errors"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateRecipe a new recipe
func CreateRecipe(recipe models.Recipe) (models.Recipe, []string, error) {
	currentTime := time.Now()
	recipe.CreatedDate = currentTime.Format("2006.01.02 15:04:05")
	recipe.LastUpdatedDate = currentTime.Format("2006.01.02 15:04:05")
	valid, invalidFields := isValidRecipe(recipe)
	if valid == false {
		return models.Recipe{}, invalidFields, errors.New("Invalid fields")
	}
	result, err := RecipeCollection.InsertOne(context.Background(), recipe)

	if err != nil {
		return models.Recipe{}, invalidFields, err
	}

	recipe.RecipeID = result.InsertedID.(primitive.ObjectID)
	return recipe, invalidFields, nil
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
