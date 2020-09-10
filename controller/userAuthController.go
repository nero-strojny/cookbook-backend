package controller

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

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

//CreateUser creates a new user
func CreateUser(requestedUser models.RequestedUser) (models.User, error) {
	// ensure there isn't another user with the same username
	result := models.User{}
	getFilter := bson.M{"username": requestedUser.UserName}
	getErr := UserCollection.FindOne(context.Background(), getFilter).Decode(&result)
	if getErr != nil {
		insertedUser := models.User{}
		insertedUser.UserName = requestedUser.UserName
		insertedUser.UserType = requestedUser.UserType
		bytes, err := bcrypt.GenerateFromPassword([]byte(requestedUser.Password), 14)
		if err != nil {
			return models.User{}, err
		}
		insertedUser.PasswordHash = string(bytes)
		result, insertErr := UserCollection.InsertOne(context.Background(), insertedUser)

		if insertErr != nil {
			return models.User{}, insertErr
		}

		insertedUser.UserID = result.InsertedID.(primitive.ObjectID)
		return insertedUser, nil
	}
	return models.User{}, errors.New("Username already taken")
}

//UpdateUserPassword updates a password
func UpdateUserPassword(updatedPassword models.UpdatedPassword) error {
	getResult := models.User{}
	getFilter := bson.M{"username": updatedPassword.UserName}
	getErr := UserCollection.FindOne(context.Background(), getFilter).Decode(&getResult)
	if getErr != nil {
		return errors.New("Username or password is not correct")
	}
	hashErr := bcrypt.CompareHashAndPassword([]byte(getResult.PasswordHash), []byte(updatedPassword.CurrentPassword))
	if hashErr != nil {
		return errors.New("Username or password is not correct")
	} else {
		bytes, err := bcrypt.GenerateFromPassword([]byte(updatedPassword.NewPassword), 14)

		if err != nil {
			return err
		}

		getResult.PasswordHash = string(bytes)
		getResult.ExpiryDate = ""
		getResult.AccessToken = ""
		opts := options.Replace().SetUpsert(false)
		updateResult, updateErr := UserCollection.ReplaceOne(context.Background(), getFilter, getResult, opts)

		if updateErr != nil {
			return updateErr
		}

		if updateResult.UpsertedID != nil {
			return errors.New("Username or password is not correct")
		}

		return nil
	}
}

//GetUsers - gets users
func GetUsers() ([]models.User, error) {
	var emptyResults []models.User
	cur, err := UserCollection.Find(context.Background(), bson.D{{}})
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
func DeleteUser(userName string) error {
	filter := bson.M{"username": userName}
	result, err := UserCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}

//GenerateUserToken generates a new token
func GenerateUserToken(authData models.AuthData) (string, error) {
	// token expires in an hour
	expiryTime := time.Now().AddDate(0, 0, 1)
	result := models.User{}
	getFilter := bson.M{"username": authData.UserName}
	getErr := UserCollection.FindOne(context.Background(), getFilter).Decode(&result)
	if getErr != nil {
		return "", errors.New("failed authentication, unknown user or password")
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(authData.Password))
	if hashErr != nil {
		result.AccessToken = ""
	} else {
		result.AccessToken = StringWithCharset(32)
		result.ExpiryDate = expiryTime.Format("2006.01.02 15:04:05")
	}

	updateFilter := bson.M{"_id": result.UserID}
	updateResult, updateErr := UserCollection.ReplaceOne(context.Background(), updateFilter, result)

	if hashErr != nil {
		return "", errors.New("failed authentication, unknown user or password")
	}

	if updateErr != nil || updateResult.ModifiedCount != 1 {
		return "", updateErr
	}

	return result.AccessToken, nil
}

//ValidateUser
func ValidateUser(accessToken string, restrictAdmin bool) error {
	if len(accessToken) == 0 {
		return errors.New("No token in request")
	}
	currentTime := time.Now().Format("2006.01.02 15:04:05")
	result := models.User{}
	filter := bson.M{"accesstoken": accessToken}
	err := UserCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return err
	} else if len(result.ExpiryDate) == 0 || result.ExpiryDate < currentTime {
		return errors.New("Expired Token")
	} else if restrictAdmin && result.UserType != "admin" {
		return errors.New("User does not have admin permissions")
	} else {
		return nil
	}
}

//ValidateSpecificUser
func ValidateSpecificUser(accessToken string, userName string) error {
	if len(accessToken) == 0 {
		return errors.New("No token in request")
	}
	currentTime := time.Now().Format("2006.01.02 15:04:05")
	result := models.User{}
	filter := bson.M{"accesstoken": accessToken}
	err := UserCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil || len(result.ExpiryDate) == 0 || result.ExpiryDate < currentTime {
		return err
	} else if result.UserName != userName {
		return errors.New("Invalid permissions")
	} else {
		return nil
	}
}
