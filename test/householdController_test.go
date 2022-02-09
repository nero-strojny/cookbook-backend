package test

import (
	"server/controller"
	"server/models"
	"testing"
)

type mockUserUpdater struct {
}

func (m mockUserUpdater) GetUser(username string, email string) (models.User, error) {
	return models.User{UserName: "SuccessulUser", HouseholdId: "OriginalID"}, nil
}

func (m mockUserUpdater) GetUserByAccessToken(token string) (models.User, error) {
	panic("implement me")
}

func (m mockUserUpdater) GetAllUsers() ([]models.User, error) {
	panic("implement me")
}

func (m mockUserUpdater) UpdatePassword(username string, oldPassword string, newPassword string) error {
	panic("implement me")
}

func (m mockUserUpdater) UpdateToken(user models.User) error {
	panic("implement me")
}

func (m mockUserUpdater) UpdateUser(user models.User) (models.User, error) {
	return models.User{HouseholdId: user.HouseholdId}, nil
}

func TestUpdateUserHousehold(t *testing.T) {
	h := controller.NewHouseholdController()
	updater := mockUserUpdater{}
	currentUser, _ := updater.GetUser("SuccessfulUser", "")
	if currentUser.HouseholdId != "OriginalID" {
		t.Fatal("Failed to get current user")
	}
	user, _ := h.AddUserToHousehold("NewID", "SuccessfulUser", updater)
	if user.HouseholdId != "NewID" {
		t.Fatal("User's house hold ID was not updated'")
	}
}


