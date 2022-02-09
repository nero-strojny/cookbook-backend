package test

import (
	"server/controller"
	"server/models"
	"testing"
	"time"
)

type expiredUserGetter struct{}

func (m expiredUserGetter) GetUser(username string, email string) (models.User, error) {
	panic("implement me")
}

func (m expiredUserGetter) GetUserByAccessToken(token string) (models.User, error) {
	return models.User{
		UserName: "TestUser",
		ExpiryDate: time.Date(2021, time.April, 13, 3, 3, 3, 3, time.UTC).Format("2006.01.02 15:04:05"),
		AccessToken: token,
	}, nil
}

func (m expiredUserGetter) GetAllUsers() ([]models.User, error) {
	panic("implement me")
}

func TestExpiredToken(t *testing.T) {
	ac := controller.NewAuthenticationController()
	err := ac.ValidateUser("FakeToken", false, expiredUserGetter{})
	if err.Error() != "expired token" {
		t.Fatal("Token was not expired")
	}
}

func TestNeedToken(t *testing.T) {
	ac := controller.NewAuthenticationController()
	err := ac.ValidateUser("", false, expiredUserGetter{})
	if err.Error() != "no token in request" {
		t.Fatal("Token was not expired")
	}
}

type nonAdminGetter struct {}

func (n nonAdminGetter) GetUser(username string, email string) (models.User, error) {
	panic("implement me")
}

func (n nonAdminGetter) GetUserByAccessToken(token string) (models.User, error) {
	return models.User{
		UserType: "normal",
		AccessToken: "token",
		ExpiryDate: time.Now().Add(5 * time.Hour).Format("2006.01.02 15:04:05"),
	}, nil
}

func (n nonAdminGetter) GetAllUsers() ([]models.User, error) {
	panic("implement me")
}

func TestNeedAdmin(t *testing.T) {
	ac := controller.NewAuthenticationController()
	err := ac.ValidateUser("token", true, nonAdminGetter{})
	if err.Error() != "user does not have admin permissions" {
		t.Fatal("User was not admin but needed to be")
	}
}

type validUserGetter struct{}

func (v validUserGetter) GetUser(username string, email string) (models.User, error) {
	panic("implement me")
}

func (v validUserGetter) GetUserByAccessToken(token string) (models.User, error) {
	return models.User{
		UserType: "normal",
		AccessToken: "token",
		ExpiryDate: time.Now().Add(5 * time.Hour).Format("2006.01.02 15:04:05"),
	}, nil
}

func (v validUserGetter) GetAllUsers() ([]models.User, error) {
	panic("implement me")
}

func TestSuccessfulAuth(t *testing.T) {
	ac := controller.NewAuthenticationController()
	err := ac.ValidateUser("token", false, validUserGetter{})
	if err != nil {
		t.Fatal("Error while authing when should be successful")
	}
}
