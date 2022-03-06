package test

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type mockCalendarGetterUpdater struct{}

func (m mockCalendarGetterUpdater) GetCalendar(householdID string, startDate string) (models.Calendar, error) {
	originalID, err1 := primitive.ObjectIDFromHex("111111111111111111111111")
	originalMonday, err2 := primitive.ObjectIDFromHex("222222222222222222222222")
	if err1 != nil {

	}

	if err2 != nil {

	}
	calendar := models.Calendar{Monday: originalMonday, CalendarID: originalID}
	return calendar, nil
}

func (m mockCalendarGetterUpdater) UpdateCalendar(calendarID string, updatedCalendar models.Calendar) (models.Calendar, error) {

	return models.Calendar{CalendarID: updatedCalendar.CalendarID, Monday: updatedCalendar.Monday}, nil
}

func TestUpdateCalendar(t *testing.T) {
	hc := controller.NewHouseholdController()
	newMonday, _ := primitive.ObjectIDFromHex("333333333333333333333333")
	newCalendar := models.Calendar{Monday: newMonday}
	calendar, _ := hc.UpdateCalendar("testHousehold", newCalendar, mockCalendarGetterUpdater{})
	if calendar.CalendarID.Hex() != "111111111111111111111111" {
		t.Fatalf("Update changed calendarID, expected originalID but got %s", calendar.CalendarID.Hex())
	}

	if calendar.Monday.Hex() == "222222222222222222222222" {
		t.Fatal("Monday ID was not updated from new calendar")
	}

	if calendar.Monday.Hex() != "333333333333333333333333" {
		t.Fatalf("Expected NewMonday but got %s", calendar.Monday.Hex())
	}
}
