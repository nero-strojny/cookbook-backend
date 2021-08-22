package controller

import (
	"context"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// getMongoPage takes the max id and page size to find a selection of recipes
func getMongoPage(pageSize int64, idLimit primitive.ObjectID, queryRecipe models.Recipe) ([]models.Recipe, error) {
	var emptyResults []models.Recipe
	filterArray := bson.A{}
	if queryRecipe.RecipeName != "" {
		regex := `(?i).*` + queryRecipe.RecipeName + `.*`
		nameFilter := bson.M{"recipename": bson.M{"$regex": regex}}
		filterArray = append(filterArray, nameFilter)
	}
	if len(queryRecipe.Tags) > 0 {
		tagFilter := bson.M{"tags": bson.M{"$all": queryRecipe.Tags}}
		filterArray = append(filterArray, tagFilter)
	}
	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	filterArray = append(filterArray, bson.M{"_id": bson.M{"$gt": idLimit}})
	cur, err := RecipeCollection.Find(
		context.Background(),
		bson.M{"$and": filterArray},
		findOptions,
	)
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
	firstRecipeBatch, firstBatchErr := getMongoPage(paginatedRequest.PageSize, id, paginatedRequest.QueryRecipe)
	if firstBatchErr != nil {
		return emptyResults, firstBatchErr
	}
	results = firstRecipeBatch
	nextID = firstRecipeBatch[len(results)-1].RecipeID

	// continue to get a batch of recipes until we reach the page number
	if paginatedRequest.PageCount > 0 {
		for i := 1; i <= paginatedRequest.PageCount; i++ {
			newresults, err := getMongoPage(paginatedRequest.PageSize, nextID, paginatedRequest.QueryRecipe)
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
	filterArray := bson.A{}
	queryRecipe := paginatedRequest.QueryRecipe
	var itemNumber int64
	var countErr error
	if queryRecipe.RecipeName != "" {
		regex := `(?i).*` + queryRecipe.RecipeName + `.*`
		nameFilter := bson.M{"recipename": bson.M{"$regex": regex}}
		filterArray = append(filterArray, nameFilter)
	}
	if len(queryRecipe.Tags) > 0 {
		tagFilter := bson.M{"tags": bson.M{"$all": queryRecipe.Tags}}
		filterArray = append(filterArray, tagFilter)
	}
	if len(filterArray) > 0 {
		itemCount, err := RecipeCollection.CountDocuments(
			context.Background(),
			bson.M{"$and": filterArray},
		)
		itemNumber = itemCount
		countErr = err
	} else {
		itemCount, err := RecipeCollection.CountDocuments(
			context.Background(),
			bson.D{{}},
		)
		itemNumber = itemCount
		countErr = err
	}
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
	recipes, getErr := PaginatedRecipes(paginatedRequest)

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
