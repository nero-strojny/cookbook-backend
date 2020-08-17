package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const calorieLogDBName = "calorieLogTable"
const calorieLogCollectionName = "calorieLogCollection"

// collection object/instance
var collection *mongo.Collection
var calorieClient *mongo.Client

//SetCalorieClient
func SetCalorieClient(c *mongo.Client) {
	calorieClient = c
	collection = calorieClient.Database(calorieLogDBName).Collection(calorieLogCollectionName)
	fmt.Println("Collection instance created!")
}

//GetAllCalorieLogs - gets all CalorieLog
func GetAllCalorieLogs() ([]primitive.M, error) {
	var emptyResults []primitive.M
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return emptyResults, err
	}

	// individually decode mongo results
	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			return emptyResults, e
		}
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		return emptyResults, err
	}

	cur.Close(context.Background())
	return results, nil
}

//GetCalorieLog - gets calorie log s by its ID
func GetCalorieLog(calorieLogID string) (models.Recipe, error) {
	result := models.Recipe{}
	id, _ := primitive.ObjectIDFromHex(calorieLogID)
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//DeleteCalorieLog - deletes a CalorieLog by its ID.
func DeleteCalorieLog(calorieLogID string) error {
	id, _ := primitive.ObjectIDFromHex(calorieLogID)
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}

//CreateCalorieLog - creates a new CalorieLog
func CreateCalorieLog(calorieLog models.CalorieLog) (models.CalorieLog, error) {
	currentTime := time.Now()
	calorieLog.EnteredDate = currentTime.Format("2006.01.02 15:04:05")
	result, err := collection.InsertOne(context.Background(), calorieLog)

	if err != nil {
		return models.CalorieLog{}, err
	}
	calorieLog.CalorieLogID = result.InsertedID.(primitive.ObjectID)
	return calorieLog, nil
}

//UpdateCalorieLog - updates an existing CalorieLog by its id
func UpdateCalorieLog(calorieLogID string, updatedCalorieLog models.CalorieLog) (models.CalorieLog, error) {
	id, _ := primitive.ObjectIDFromHex(calorieLogID)
	filter := bson.M{"_id": id}
	//Could do this as an update but that requires checking what fields are different between calorie log s
	//Could be a hassle with a long list of ingredients or measurements. Easier to just replace the entire calorie log  with the new update
	opts := options.Replace().SetUpsert(true)
	result, err := collection.ReplaceOne(context.Background(), filter, updatedCalorieLog, opts)

	if err != nil {
		return models.CalorieLog{}, err
	}
	if result.UpsertedID != nil {
		updatedCalorieLog.CalorieLogID = result.UpsertedID.(primitive.ObjectID)
	}
	return updatedCalorieLog, nil
}
