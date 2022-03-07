package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"server/models"
)

type HouseholdDB interface {
	HouseholdGetter
	HouseholdCreator
	HouseholdDeleter
}

type HouseholdGetter interface {
	GetHousehold(householdID string) (models.Household, error)
}

type HouseholdCreator interface {
	CreateHousehold(household models.Household, user models.User) (models.Household, error)
}

type HouseholdDeleter interface {
	DeleteHousehold(householdID string) error
}

type HouseholdRepository struct {
	householdCollection *mongo.Collection
}

func NewHouseholdRepository(client *mongo.Client) *HouseholdRepository {
	return &HouseholdRepository{
		householdCollection: client.Database("tastyBoiDatabase").Collection("householdCollection"),
	}
}

func (h HouseholdRepository) GetHousehold(householdID string) (models.Household, error) {
	result := models.Household{}
	id, _ := primitive.ObjectIDFromHex(householdID)
	filter := bson.M{"_id": id}
	err := h.householdCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (h HouseholdRepository) CreateHousehold(household models.Household, user models.User) (models.Household, error) {
	household.HeadOfHousehold = user.UserID.Hex()
	result, err := h.householdCollection.InsertOne(context.Background(), household)

	if err != nil {
		return models.Household{}, err
	}

	household.HouseholdID = result.InsertedID.(primitive.ObjectID)
	return household, nil
}

func (h HouseholdRepository) DeleteHousehold(householdID string) error {
	id, _ := primitive.ObjectIDFromHex(householdID)
	filter := bson.M{"_id": id}
	result, err := h.householdCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("nothing was deleted")
	}
	return nil
}
