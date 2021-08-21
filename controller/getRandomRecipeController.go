package controller

import (
	"context"
	"math/rand"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
)

// GetRandomRecipes
func GetRandomRecipes(numberOfRecipes int) ([]models.Recipe, error) {
	var emptyResults []models.Recipe
	itemCount, err := RecipeCollection.CountDocuments(context.Background(), bson.D{{}})
	if err != nil {
		return emptyResults, err
	}
	if numberOfRecipes > int(itemCount) {
		numberOfRecipes = int(itemCount)
	} else if numberOfRecipes < 1 {
		return emptyResults, err
	}
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	var recipeNumbers []int
	doneComputing := false
	for !doneComputing {
		recipeNumber := seededRand.Intn(int(itemCount))
		if !contains(recipeNumbers, recipeNumber) {
			recipeNumbers = append(recipeNumbers, recipeNumber)
			if len(recipeNumbers) == numberOfRecipes {
				doneComputing = true
			}
		}
	}
	var randomRecipes []models.Recipe
	for _, num := range recipeNumbers {
		pageSize := 10
		pageCount := num / pageSize
		selectedRecipe := num % pageSize
		paginatedRequest := models.PaginatedRecipeRequest{
			PageCount:   pageCount,
			PageSize:    int64(pageSize),
			QueryRecipe: models.Recipe{},
		}
		recipes, getErr := PaginatedRecipes(paginatedRequest)
		if getErr != nil {
			return emptyResults, err
		}
		randomRecipes = append(randomRecipes, recipes[selectedRecipe])
	}
	return randomRecipes, nil
}

func contains(recipeNumbers []int, recipeNumber int) bool {
	for _, num := range recipeNumbers {
		if num == recipeNumber {
			return true
		}
	}
	return false
}
