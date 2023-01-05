package test

import (
	"errors"
	"server/controller"
	"server/models"
	"testing"
)

// Good examples of why the new architecture is better
// By breaking database calls into their own structure, and having that structure implement a composite interface,
// we can simplify our mocking to just the methods we want to test instead of the whole db interface
type mockUserCreator struct{}

type failedUserCreator struct{}

func (f failedUserCreator) CreateUser(userInformation models.RequestedUser) (models.User, error) {
	return models.User{}, errors.New("user was empty")
}

func (m mockUserCreator) CreateUser(userInformation models.RequestedUser) (models.User, error) {
	return models.User{UserName: userInformation.UserName}, nil
}

func TestCreateUser(t *testing.T) {
	c := controller.NewUserController()
	user := models.RequestedUser{UserName: "TEST"}
	output, err := c.CreateUser(user, mockUserCreator{})
	if err != nil {
		t.Fatal("Error was returned")
	}

	if output.UserName != "TEST" {
		t.Fatal("User was not created")
	}
}

func TestCreateUserFailure(t *testing.T) {
	c := controller.NewUserController()
	user := models.RequestedUser{UserName: "TEST"}
	output, err := c.CreateUser(user, failedUserCreator{})
	if err == nil {
		t.Fatal("Fail did not return non nil error")
	}

	if output.UserName != "" {
		t.Fatal("Fail did not return empty user")
	}
}
