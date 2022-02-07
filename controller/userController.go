package controller

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/smtp"
	"server/config"
	"server/db"
	"server/models"
	"time"
)

const (
	from     = "tasty.boi.shopping.list@gmail.com"
	subject = "Subject: Grocery List\r\n\r\n"
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

type UserControl interface {
	CreateUser(requestedUser models.RequestedUser, repository db.UserCreator) (models.User, error)
	UpdateUserPassword(updatedPassword models.UpdatedPassword, repository db.UserUpdater) error
	GetUsers(repository db.UserGetter) ([]models.User, error)
	DeleteUser(userName string, repository db.UserDeleter) error
	GenerateUserToken(authData models.AuthData, repository db.UserGetterUpdater) (string, error)
	EmailUser(basket models.Basket, token string, repository db.UserGetter) error
}

type UserController struct {
}

func NewUserController() UserController {
	return UserController{}
}

//CreateUser creates a new user
func (uc UserController) CreateUser(requestedUser models.RequestedUser, repository db.UserCreator) (models.User, error) {
	user, err := repository.CreateUser(requestedUser)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

//UpdateUserPassword updates a password
func (uc UserController) UpdateUserPassword(updatedPassword models.UpdatedPassword, repository db.UserUpdater) error {
	err := repository.UpdatePassword(updatedPassword.UserName, updatedPassword.CurrentPassword, updatedPassword.NewPassword)
	if err != nil {
		return err
	}
	return nil
}

//GetUsers - gets users
func (uc UserController) GetUsers(repository db.UserGetter) ([]models.User, error) {
	users, err := repository.GetAllUsers()
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

//DeleteUser - deletes a User by its ID.
func (uc UserController) DeleteUser(userName string, repository db.UserDeleter) error {
	err := repository.DeleteUser(userName)
	if err != nil {
		return err
	}
	return nil
}

//GenerateUserToken generates a new token
func (uc UserController) GenerateUserToken(authData models.AuthData, repository db.UserGetterUpdater) (string, error) {
	// token expires in an hour
	expiryTime := time.Now().AddDate(0, 0, 1)

	user, err := repository.GetUser(authData.UserName, "")
	if err != nil {
		return "", errors.New("failed authentication, unknown user or password")
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(authData.Password))
	if hashErr != nil {
		user.AccessToken = ""
	} else {
		user.AccessToken = stringWithCharset(32)
		user.ExpiryDate = expiryTime.Format("2006.01.02 15:04:05")
	}

	repository.UpdateToken(user)

	return user.AccessToken, nil
}



func (uc UserController) EmailUser(basket models.Basket, token string, repository db.UserGetter) error {
	userEmail, err := uc.lookUpUserEmail(token, repository)

	if err != nil {
		return err
	}
	shoppingList := buildCategoryString("Produce", basket.Produce)
	shoppingList += buildCategoryString("Pantry", basket.Pantry)
	shoppingList += buildCategoryString("Protein", basket.Protein)
	shoppingList += buildCategoryString("Dairy", basket.Dairy)
	shoppingList += buildCategoryString("Alcohol", basket.Alcohol)

	to := []string{userEmail}
	return sendEmail(subject, "", shoppingList, to)
}

func buildCategoryString(category string, ingredients []string) string {
	if len(ingredients) == 0 {
		return ""
	}
	output := category + "\n"
	for _, ingredient := range ingredients {
		output += ingredient + "\n"
	}
	output += "\n"
	return output
}

func (uc UserController) lookUpUserEmail(token string, repository db.UserGetter) (email string, error error) {
	user, err := repository.GetUserByAccessToken(token)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func sendEmail(subject string, mime string, body string, recipients []string) error {
	auth := smtp.PlainAuth("", from, config.GetConfig().EmailPassword, smtpHost)
	email := []byte(subject + mime + body)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, email)
}

//StringWithCharset generates a random string
func stringWithCharset(length int) string {
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
