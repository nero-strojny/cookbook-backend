package controller

import (
	"context"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//GetIngredient - gets ingredients by its ID
func GetIngredient(ingredientID string) (models.Ingredient, error) {
	result := models.Ingredient{}
	id, _ := primitive.ObjectIDFromHex(ingredientID)
	filter := bson.M{"_id": id}
	err := IngredientCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// QueryIngredient returns all ingredients matching a prefix
func QueryIngredient(prefixIngredient string) ([]models.Ingredient, error) {
	// Get the first page
	var emptyResults []models.Ingredient
	var emptyIngredients []models.Ingredient
	pageSize := int64(5)
	regex := `(?i)^` + prefixIngredient
	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	cur, err := IngredientCollection.Find(
		context.Background(),
		bson.M{"name": bson.M{"$regex": regex}},
		findOptions,
	)
	if err != nil {
		return emptyIngredients, err
	}

	ingredientBatch, batchErr := decodeCurToIngredients(cur)
	if batchErr != nil {
		return emptyResults, batchErr
	}
	return ingredientBatch, nil
}

func decodeCurToIngredients(cur *mongo.Cursor) ([]models.Ingredient, error) {
	emptyResults := []models.Ingredient{}
	var results []models.Ingredient
	for cur.Next(context.Background()) {
		result := models.Ingredient{}
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
