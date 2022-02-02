package controller

import (
	"context"
	"errors"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//CreateHousehold creates a new household
func CreateHousehold(household models.Household, user models.User) (models.Household, error) {
	household.HeadOfHousehold = user.UserID.Hex()
	result, err := HouseholdCollection.InsertOne(context.Background(), household)

	if err != nil {
		return models.Household{}, err
	}

	household.HouseholdID = result.InsertedID.(primitive.ObjectID)
	return household, nil
}

//GetHousehold - gets households by its ID
func GetHousehold(householdID string) (models.Household, error) {
	result := models.Household{}
	id, _ := primitive.ObjectIDFromHex(householdID)
	filter := bson.M{"_id": id}
	err := HouseholdCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//AddUserToHousehold - updates user's household id field
func AddUserToHousehold(householdID string, userID string) (models.User, error) {
	updatedUser, getUserErr := GetUser(userID)
	if getUserErr != nil {
		return models.User{}, getUserErr
	}

	updatedUser.HouseholdId = householdID
	id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": id}

	opts := options.Replace().SetUpsert(true)
	_, err := UserCollection.ReplaceOne(context.Background(), filter, updatedUser, opts)

	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

//DeleteHousehold- deletes a household by its ID.
func DeleteHousehold(householdID string) error {
	id, _ := primitive.ObjectIDFromHex(householdID)
	filter := bson.M{"_id": id}
	result, err := HouseholdCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}
