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
	filterArray := bson.A{}
	idFilter := bson.M{"householdID": householdIDObject}
	dateFilter := bson.M{"startDate": startDate}
	filterArray = append(filterArray, idFilter)
	filterArray = append(filterArray, dateFilter)

	err := c.calendarCollection.FindOne(
		context.Background(),
		bson.M{"$and": filterArray},
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

func (c CalendarRepository) UpdateCalendar(calendarID string, updatedCalendar models.Calendar) (models.Calendar, error) {
	id, _ := primitive.ObjectIDFromHex(calendarID)
	filter := bson.M{"_id": id}
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
