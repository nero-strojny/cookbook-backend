package controller

import (
	"server/db"
	"server/models"
)

type IngredientControl interface {
	CreateIngredient(ingredient models.Ingredient, repository db.IngredientCreator) (models.Ingredient, error)
	DeleteIngredient(ingredientID string, repository db.IngredientDeleter) error
	GetIngredient(ingredientID string, repository db.IngredientGetter) (models.Ingredient, error)
	QueryIngredient(prefixIngredient string, repository db.IngredientGetter) ([]models.Ingredient, error)
}

type IngredientController struct {
}

func NewIngredientController() IngredientController {
	return IngredientController{}
}

//CreateIngredient creates a new ingredient
func (ic IngredientController) CreateIngredient(ingredient models.Ingredient, repository db.IngredientCreator) (models.Ingredient, error) {
	newIngredient, err := repository.CreateIngredient(ingredient)
	if err != nil {
		return models.Ingredient{}, err
	}
	return newIngredient, nil
}

//DeleteIngredient - deletes a ingredient by its ID.
func (ic IngredientController) DeleteIngredient(ingredientID string, repository db.IngredientDeleter) error {
	return repository.DeleteIngredient(ingredientID)
}

//GetIngredient - gets ingredients by its ID
func (ic IngredientController) GetIngredient(ingredientID string, repository db.IngredientGetter) (models.Ingredient, error) {
	return repository.GetIngredient(ingredientID)
}

// QueryIngredient returns all ingredients matching a prefix
func (ic IngredientController) QueryIngredient(prefixIngredient string, repository db.IngredientGetter) ([]models.Ingredient, error) {
	return repository.QueryIngredients(prefixIngredient)
}

