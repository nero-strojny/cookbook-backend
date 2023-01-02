package db

import (
	"context"
	"errors"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserDB interface {
	UserGetter
	UserDeleter
	UserUpdater
}

type UserGetterUpdater interface {
	UserGetter
	UserUpdater
}

type UserGetter interface {
	GetUser(username string, email string) (models.User, error)
	GetUserByAccessToken(token string) (models.User, error)
	GetAllUsers() ([]models.User, error)
}

type UserDeleter interface {
	DeleteUser(username string) error
}

type UserUpdater interface {
	UpdatePassword(username string, oldPassword string, newPassword string) error
	UpdateToken(user models.User) error
	UpdateUser(user models.User) (models.User, error)
}

type UserRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{
		userCollection: client.Database("tastyBoiDatabase").Collection("userCollection"),
	}
}

func (ur UserRepository) GetUser(username string, email string) (models.User, error) {
	user := models.User{}
	filter := bson.M{"$or": []bson.M{{"username": username}, {"email": email}}}
	err := ur.userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return models.User{}, errors.New("no user with that name or email")
	}
	return user, nil
}

func (ur UserRepository) GetUserByAccessToken(token string) (models.User, error) {
	user := models.User{}
	filter := bson.M{"accesstoken": token}
	err := ur.userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return models.User{}, errors.New("no user with that accses token")
	}
	return user, nil
}

func (ur UserRepository) UpdatePassword(username string, oldPassword string, newPassword string) error {
	user := models.User{}
	filter := bson.M{"username": username}
	getErr := ur.userCollection.FindOne(context.Background(), filter).Decode(&user)
	if getErr != nil {
		return errors.New("username or password is not correct")
	}
	hashErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if hashErr != nil {
		return errors.New("username or password is not correct")
	} else {
		bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)

		if err != nil {
			return err
		}

		user.PasswordHash = string(bytes)
		user.ExpiryDate = ""
		user.AccessToken = ""
		opts := options.Replace().SetUpsert(false)
		updateResult, updateErr := ur.userCollection.ReplaceOne(context.Background(), filter, user, opts)

		if updateErr != nil {
			return updateErr
		}

		if updateResult.UpsertedID != nil {
			return errors.New("username or password is not correct")
		}
		return nil
	}
}

func (ur UserRepository) DeleteUser(username string) error {
	filter := bson.M{"username": username}
	result, err := ur.userCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("nothing was deleted")
	}
	return nil
}

func (ur UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	cur, err := ur.userCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return users, err
	}

	// individually decode mongo results
	for cur.Next(context.Background()) {
		result := models.User{}
		e := cur.Decode(&result)
		if e != nil {
			return users, e
		}
		users = append(users, result)

	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	cur.Close(context.Background())
	return users, nil
}

func (ur UserRepository) UpdateToken(user models.User) error {
	updateFilter := bson.M{"_id": user.UserID}
	setOperation := bson.D{{"$set", bson.D{{"accesstoken", user.AccessToken}, {"expirydate", user.ExpiryDate}}}}
	updateResult, updateErr := ur.userCollection.UpdateOne(context.Background(), updateFilter, setOperation)

	if updateErr != nil || updateResult.ModifiedCount != 1 {
		return updateErr
	}
	return nil
}

func (ur UserRepository) UpdateUser(user models.User) (models.User, error) {
	filter := bson.M{"_id": user.UserID}

	opts := options.Replace().SetUpsert(true)
	_, err := ur.userCollection.ReplaceOne(context.Background(), filter, user, opts)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
