package controller

import (
	"context"
	"errors"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//GetCalorieLog - gets calorie log s by its ID
func GetCalorieLog(calorieLogID string) (models.Recipe, error) {
	result := models.Recipe{}
	id, _ := primitive.ObjectIDFromHex(calorieLogID)
	filter := bson.M{"_id": id}
	err := CalorieLogCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//DeleteCalorieLog - deletes a CalorieLog by its ID.
func DeleteCalorieLog(calorieLogID string) error {
	id, _ := primitive.ObjectIDFromHex(calorieLogID)
	filter := bson.M{"_id": id}
	result, err := CalorieLogCollection.DeleteOne(context.Background(), filter)
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
	result, err := CalorieLogCollection.InsertOne(context.Background(), calorieLog)

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
	result, err := CalorieLogCollection.ReplaceOne(context.Background(), filter, updatedCalorieLog, opts)

	if err != nil {
		return models.CalorieLog{}, err
	}
	if result.UpsertedID != nil {
		updatedCalorieLog.CalorieLogID = result.UpsertedID.(primitive.ObjectID)
	}
	return updatedCalorieLog, nil
}
