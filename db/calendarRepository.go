package db

import (
	"context"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CalendarDB interface {
	CalendarGetter
	CalendarCreator
	CalendarDeleter
	CalendarUpdater
}

type CalendarGetterUpdater interface {
	CalendarGetter
	CalendarUpdater
}

type CalendarGetter interface {
	GetCalendar(householdID string, startDate string) (models.Calendar, error)
}

type CalendarCreator interface {
	CreateCalendar(calendar models.Calendar) (models.Calendar, error)
}

type CalendarDeleter interface {
	DeleteCalendar(householdID string) error
}

type CalendarUpdater interface {
	UpdateCalendar(updatedCalendar models.Calendar) (models.Calendar, error)
}

type CalendarRepository struct {
	calendarCollection *mongo.Collection
}

func NewCalendarRepository(client *mongo.Client) *CalendarRepository {
	return &CalendarRepository{
		calendarCollection: client.Database("tastyBoiDatabase").Collection("calendarCollection"),
	}
}

func (c CalendarRepository) GetCalendar(householdID string, startDate string) (models.Calendar, error) {
	result := models.Calendar{}
	householdIDObject, _ := primitive.ObjectIDFromHex(householdID)
	filter := bson.M{"householdID": householdIDObject, "startdate": startDate}

	err := c.calendarCollection.FindOne(
		context.Background(),
		filter,
	).Decode(&result)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (c CalendarRepository) CreateCalendar(calendar models.Calendar) (models.Calendar, error) {
	result, err := c.calendarCollection.InsertOne(context.Background(), calendar)

	if err != nil {
		return models.Calendar{}, err
	}

	calendar.CalendarID = result.InsertedID.(primitive.ObjectID)
	return calendar, nil
}

func (c CalendarRepository) UpdateCalendar(updatedCalendar models.Calendar) (models.Calendar, error) {
	filter := bson.M{"_id": updatedCalendar.CalendarID}
	opts := options.Replace().SetUpsert(true)
	result, err := c.calendarCollection.ReplaceOne(context.Background(), filter, updatedCalendar, opts)
	if err != nil {
		return models.Calendar{}, err
	}
	if result.UpsertedID != nil {
		updatedCalendar.CalendarID = result.UpsertedID.(primitive.ObjectID)
	}
	return updatedCalendar, nil
}

func (c CalendarRepository) DeleteCalendar(householdID string) error {
	panic("implement me")
}
