package controller

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// getMongoPage takes the max id and page size to find a selection of recipes
func getMongoPage(pageSize int64, idLimit primitive.ObjectID) ([]models.Recipe, error) {
	var emptyResults []models.Recipe
	// individually decode mongo results
	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	filter := bson.M{"_id": bson.M{"$gt": idLimit}}
	cur, err := RecipeCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return emptyResults, err
	}

	return decodeCurToRecipes(cur)
}

// PaginatedRecipes returns all recipes matching a paginated request
func PaginatedRecipes(paginatedRequest models.PaginatedRecipeRequest) ([]models.Recipe, error) {
	// Get the first page
	var emptyResults []models.Recipe
	var results []models.Recipe
	var nextID primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex("0")
	firstRecipeBatch, firstBatchErr := getMongoPage(paginatedRequest.PageSize, id)
	if firstBatchErr != nil {
		return emptyResults, firstBatchErr
	}
	results = firstRecipeBatch
	nextID = firstRecipeBatch[len(results)-1].RecipeID

	// continue to get a batch of recipes until we reach the page number
	if paginatedRequest.PageCount > 0 {
		for i := 1; i <= paginatedRequest.PageCount; i++ {
			newresults, err := getMongoPage(paginatedRequest.PageSize, nextID)
			if err != nil {
				return emptyResults, err
			}
			if len(newresults) == 0 {
				return results, nil
			}
			results = newresults
			nextID = newresults[len(results)-1].RecipeID
		}
	}
	return results, nil
}

//PostPaginateRecipes - gets all recipes
func PostPaginatedRecipes(paginatedRequest models.PaginatedRecipeRequest) (models.PaginatedRecipeResponse, error) {
	itemCount, err := RecipeCollection.CountDocuments(context.Background(), bson.D{{}})
	if err != nil {
		return models.PaginatedRecipeResponse{}, err
	}

	recipes, getErr := PaginatedRecipes(paginatedRequest)

	if getErr != nil {
		return models.PaginatedRecipeResponse{}, err
	}

	return models.PaginatedRecipeResponse{
		Recipes:         recipes,
		NumberOfRecipes: itemCount,
		PageCount:       paginatedRequest.PageCount,
		PageSize:        paginatedRequest.PageSize,
	}, nil
}

//GetRecipe - gets recipes by its ID
func GetRecipe(recipeID string) (models.Recipe, error) {
	result := models.Recipe{}
	id, _ := primitive.ObjectIDFromHex(recipeID)
	filter := bson.M{"_id": id}
	err := RecipeCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func contains(recipeNumbers []int, recipeNumber int) bool {
	for _, num := range recipeNumbers {
		if num == recipeNumber {
			return true
		}
	}
	return false
}

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
			PageCount: pageCount,
			PageSize:  int64(pageSize),
		}
		recipes, getErr := PaginatedRecipes(paginatedRequest)
		if getErr != nil {
			return emptyResults, err
		}
		randomRecipes = append(randomRecipes, recipes[selectedRecipe])
	}
	return randomRecipes, nil
}

//SearchRecipeByName - searches for a recipe by name
func QueryRecipe(recipe models.Recipe) ([]models.Recipe, error) {
	emptyResults := []models.Recipe{}
	filterArray := bson.A{}
	if recipe.RecipeName != "" {
		regex := `(?i).*` + recipe.RecipeName + `.*`
		nameFilter := bson.M{"recipename": bson.M{"$regex": regex}}
		filterArray = append(filterArray, nameFilter)
	}
	if len(recipe.Tags) > 0 {
		tagFilter := bson.M{"tags": bson.M{"$all": recipe.Tags}}
		filterArray = append(filterArray, tagFilter)
	}
	cur, err := RecipeCollection.Find(
		context.Background(),
		bson.M{"$and": filterArray})
	if err != nil {
		return emptyResults, err
	}
	return decodeCurToRecipes(cur)
}

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

//UpdateRecipe - updates an existing recipe by its id
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

func isValidRecipe(recipe models.Recipe) (valid bool, invalidFields []string) {

	if recipe.RecipeName == "" {
		invalidFields = append(invalidFields, "recipeName")
	}

	if len(invalidFields) > 0 {
		return false, invalidFields
	}
	return true, invalidFields

}

func decodeCurToRecipes(cur *mongo.Cursor) ([]models.Recipe, error) {
	emptyResults := []models.Recipe{}
	var results []models.Recipe
	for cur.Next(context.Background()) {
		result := models.Recipe{}
		e := cur.Decode(&result)
		if e != nil {
			return emptyResults, e
		}
		results = append(results, result)

	}

	if err := cur.Err(); err != nil || len(results) == 0 {
		return emptyResults, err
	}

	cur.Close(context.Background())
	return results, nil
}
