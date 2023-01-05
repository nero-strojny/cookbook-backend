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
	h := controller.NewHouseholdController(nil, nil, nil)
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

type mockCalendarDB struct {
}

func (m mockCalendarDB) CreateCalendar(calendar models.Calendar) (models.Calendar, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockCalendarDB) DeleteCalendar(householdID string) error {
	//TODO implement me
	panic("implement me")
}

func (m mockCalendarDB) GetHousehold(householdID string) (models.Household, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockCalendarDB) CreateHousehold(household models.Household, user models.User) (models.Household, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockCalendarDB) DeleteHousehold(householdID string) error {
	//TODO implement me
	panic("implement me")
}

func (m mockCalendarDB) GetCalendar(householdID string, startDate string) (models.Calendar, error) {
	originalID, err1 := primitive.ObjectIDFromHex("111111111111111111111111")
	originalMonday, err2 := primitive.ObjectIDFromHex("222222222222222222222222")
	if err1 != nil {

	}

	if err2 != nil {

	}
	calendar := models.Calendar{Monday: models.Recipe{RecipeID: originalMonday}, CalendarID: originalID}
	return calendar, nil
}

func (m mockCalendarDB) UpdateCalendar(updatedCalendar models.Calendar) (models.Calendar, error) {
	return models.Calendar{CalendarID: updatedCalendar.CalendarID, Monday: updatedCalendar.Monday}, nil
}

func TestUpdateCalendar(t *testing.T) {
	hc := controller.NewHouseholdController(mockCalendarDB{}, nil, nil)
	newMonday, _ := primitive.ObjectIDFromHex("333333333333333333333333")
	calendarID, _ := primitive.ObjectIDFromHex("111111111111111111111111")
	newCalendar := models.Calendar{CalendarID: calendarID, Monday: models.Recipe{RecipeID: newMonday}}
	calendar, _ := hc.UpdateCalendar("testHousehold", newCalendar)
	if calendar.CalendarID.Hex() != "111111111111111111111111" {
		t.Fatalf("Update changed calendarID, expected originalID but got %s", calendar.CalendarID.Hex())
	}

	if calendar.Monday.RecipeID.Hex() == "222222222222222222222222" {
		t.Fatal("Monday ID was not updated from new calendar")
	}

	if calendar.Monday.RecipeID.Hex() != "333333333333333333333333" {
		t.Fatalf("Expected NewMonday but got %s", calendar.Monday.RecipeID)
	}
}
