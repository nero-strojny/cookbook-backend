package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/models"
)

type IngredientDB interface {
	IngredientGetter
	IngredientDeleter
	IngredientCreator
}

type IngredientGetter interface {
	GetIngredient(ingredientID string) (models.Ingredient, error)
	QueryIngredients(prefix string) ([]models.Ingredient, error)
}

type IngredientDeleter interface {
	DeleteIngredient(ingredientID string) error
}

type IngredientCreator interface {
	CreateIngredient(ingredient models.Ingredient) (models.Ingredient, error)
}

type IngredientRepository struct {
	ingredientCollection *mongo.Collection
}

func NewIngredientRepository(client *mongo.Client) *IngredientRepository {
	return &IngredientRepository{
		ingredientCollection: client.Database("tastyBoiDatabase").Collection("ingredientCollection"),
	}
}

func (i IngredientRepository) CreateIngredient(ingredient models.Ingredient) (models.Ingredient, error) {
	result, err := i.ingredientCollection.InsertOne(context.Background(), ingredient)

	if err != nil {
		return models.Ingredient{}, err
	}

	ingredient.IngredientID = result.InsertedID.(primitive.ObjectID)
	return ingredient, nil
}

func (i IngredientRepository) DeleteIngredient(ingredientID string) error {
	id, _ := primitive.ObjectIDFromHex(ingredientID)
	filter := bson.M{"_id": id}
	result, err := i.ingredientCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("nothing was deleted")
	}
	return nil
}

func (i IngredientRepository) GetIngredient(ingredientID string) (models.Ingredient, error) {
	result := models.Ingredient{}
	id, _ := primitive.ObjectIDFromHex(ingredientID)
	filter := bson.M{"_id": id}
	err := i.ingredientCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return models.Ingredient{}, err
	}
	return result, nil
}

func (i IngredientRepository) QueryIngredients(prefix string) ([]models.Ingredient, error) {
	// Get the first page
	var emptyResults []models.Ingredient
	var emptyIngredients []models.Ingredient
	pageSize := int64(5)
	regex := `(?i)^` + prefix
	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	cur, err := i.ingredientCollection.Find(
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
