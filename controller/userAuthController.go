package controller

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userDBName = "calorieLogTable"
const userCollectionName = "calorieLogCollection"

// collection object/instance
var userCollection *mongo.Collection
var userClient *mongo.Client

//StringWithCharset generates a random string
func StringWithCharset(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//SetUserClient sets the client to connect to user database
func SetUserClient(c *mongo.Client) {
	userClient = c
	userCollection = userClient.Database(userDBName).Collection(userCollectionName)
	fmt.Println("Collection instance created!")
}

//CreateUser creates a new user
func CreateUser(user models.User) (models.User, error) {
	result, err := userCollection.InsertOne(context.Background(), user)

	if err != nil {
		return models.User{}, err
	}

	user.UserID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

//GetUsers - gets users by name
func GetUsers() ([]models.User, error) {
	var emptyResults []models.User
	cur, err := userCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return emptyResults, err
	}

	// individually decode mongo results
	var results []models.User
	for cur.Next(context.Background()) {
		result := models.User{}
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

//DeleteUser - deletes a User by its ID.
func DeleteUser(userID string) error {
	id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": id}
	result, err := userCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}

//GenerateUserToken gernates a new token
func GenerateUserToken(authData models.AuthData) (string, error) {
	// token expires in an hour
	expiryTime := time.Now().AddDate(0, 0, 1)
	result := models.User{}
	getFilter := bson.M{"username": authData.UserName}
	getErr := userCollection.FindOne(context.Background(), getFilter).Decode(&result)
	if getErr != nil || result.PasswordHash != authData.PasswordHash {
		return "failed authentication, unknown user or password", getErr
	}

	result.AccessToken = StringWithCharset(32)
	result.ExpiryDate = expiryTime.Format("2006.01.02 15:04:05")

	updateFilter := bson.M{"_id": result.UserID}
	updateResult, updateErr := userCollection.ReplaceOne(context.Background(), updateFilter, result)

	if updateErr != nil || updateResult.ModifiedCount != 1 {
		return "failed authentication", updateErr
	}

	return result.AccessToken, nil
}
