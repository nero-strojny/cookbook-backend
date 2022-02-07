package controller

import (
	"errors"
	"math/rand"
	"server/db"
	"server/models"
	"time"
)

type RecipeControl interface {
	CreateRecipe(recipe models.Recipe, repository db.RecipeCreator) (models.Recipe, []string, error)
	DeleteRecipe(recipeID string, repository db.RecipeDeleter) error
	UpdateRecipe(recipeID string, updatedRecipe models.Recipe, repository db.RecipeUpdater) (models.Recipe, error)
	GetRandomRecipes(numberOfRecipes int, repository db.RecipeGetter) ([]models.Recipe, error)
	PostPaginatedRecipes(paginatedRequest models.PaginatedRecipeRequest, repository db.RecipeGetter) (models.PaginatedRecipeResponse, error)
	GetRecipe(recipeID string, repository db.RecipeGetter) (models.Recipe, error)
}

type RecipeController struct {
}

func NewRecipeController() RecipeController {
	return RecipeController{}
}

//CreateRecipe a new recipe
func (rc RecipeController) CreateRecipe(recipe models.Recipe, repository db.RecipeCreator) (models.Recipe, []string, error) {
	currentTime := time.Now()
	recipe.CreatedDate = currentTime.Format("2006.01.02 15:04:05")
	recipe.LastUpdatedDate = currentTime.Format("2006.01.02 15:04:05")
	valid, invalidFields := isValidRecipe(recipe)
	if !valid {
		return models.Recipe{}, invalidFields, errors.New("invalid fields")
	}
	err := repository.CreateRecipe(&recipe)

	if err != nil {
		return models.Recipe{}, invalidFields, err
	}

	return recipe, invalidFields, nil
}

//DeleteRecipe - deletes a recipe by its ID.
func (rc RecipeController) DeleteRecipe(recipeID string, repository db.RecipeDeleter) error {
	err := repository.DeleteRecipe(recipeID)
	if err != nil {
		return err
	}
	return nil
}

func (rc RecipeController) UpdateRecipe(recipeID string, updatedRecipe models.Recipe, repository db.RecipeUpdater) (models.Recipe, error) {
	recipe, err := repository.UpdateRecipe(recipeID, updatedRecipe)
	if err != nil {
		return models.Recipe{}, err
	}
	return recipe, nil
}

// GetRandomRecipes - returns a slice of random recipes
func (rc RecipeController) GetRandomRecipes(numberOfRecipes int, repository db.RecipeGetter) ([]models.Recipe, error) {
	var emptyResults []models.Recipe
	itemCount, err := repository.CountRecipes()
	if err != nil  {
		return emptyResults, err
	}
	if numberOfRecipes > int(itemCount) {
		numberOfRecipes = int(itemCount)
	} else if numberOfRecipes < 1 {
		return emptyResults, err
	}
	var seededRand = rand.New(
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
		recipes, getErr := repository.GetPaginatedRecipes(paginatedRequest)
		if getErr != nil {
			return emptyResults, err
		}
		randomRecipes = append(randomRecipes, recipes[selectedRecipe])
	}
	return randomRecipes, nil
}

//PostPaginateRecipes - gets all recipes
func (rc RecipeController) PostPaginatedRecipes(paginatedRequest models.PaginatedRecipeRequest, repository db.RecipeGetter) (models.PaginatedRecipeResponse, error) {
	itemNumber, countErr := repository.GetFilteredRecipeCount(paginatedRequest)
	if countErr != nil {
		return models.PaginatedRecipeResponse{}, countErr
	}
	if itemNumber == int64(0) {
		return models.PaginatedRecipeResponse{
			Recipes:         []models.Recipe{},
			NumberOfRecipes: itemNumber,
			PageCount:       paginatedRequest.PageCount,
			PageSize:        paginatedRequest.PageSize,
		}, nil
	}
	recipes, getErr := repository.GetPaginatedRecipes(paginatedRequest)

	if getErr != nil {
		return models.PaginatedRecipeResponse{}, getErr
	}

	return models.PaginatedRecipeResponse{
		Recipes:         recipes,
		NumberOfRecipes: itemNumber,
		PageCount:       paginatedRequest.PageCount,
		PageSize:        paginatedRequest.PageSize,
	}, nil
}

//GetRecipe - gets recipes by its ID
func (rc RecipeController) GetRecipe(recipeID string, repository db.RecipeGetter) (models.Recipe, error) {
	return repository.GetRecipe(recipeID)
}


func contains(recipeNumbers []int, recipeNumber int) bool {
	for _, num := range recipeNumbers {
		if num == recipeNumber {
			return true
		}
	}
	return false
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
