package controller

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/db"
	"server/models"
)

type HouseholdControl interface {
	CreateHousehold(household models.Household, user models.User) (models.Household, error)
	GetHousehold(householdID string) (models.Household, error)
	AddUserToHousehold(householdID string, username string, ur db.UserGetterUpdater) (models.User, error)
	DeleteHousehold(householdID string) error
	GetCalendar(householdID string, startDate string) (models.Calendar, error)
	UpdateCalendar(householdID string, calendar models.Calendar) (models.Calendar, error)
	CreateCalendar(startDate string, householdID string) (models.Calendar, error)
}

type HouseholdController struct {
	calendarRepo  db.CalendarDB
	householdRepo db.HouseholdDB
	rc            RecipeControl
}

func NewHouseholdController(cr db.CalendarDB, hr db.HouseholdDB, rc RecipeControl) HouseholdController {
	return HouseholdController{calendarRepo: cr, householdRepo: hr, rc: rc}
}

//CreateHousehold creates a new household
func (hc HouseholdController) CreateHousehold(household models.Household, user models.User) (models.Household, error) {
	return hc.householdRepo.CreateHousehold(household, user)
}

//GetHousehold - gets households by its ID
func (hc HouseholdController) GetHousehold(householdID string) (models.Household, error) {
	return hc.householdRepo.GetHousehold(householdID)
}

//AddUserToHousehold - updates user's household id field
func (hc HouseholdController) AddUserToHousehold(householdID string, username string, ur db.UserGetterUpdater) (models.User, error) {
	updatedUser, getUserErr := ur.GetUser(username, "")
	if getUserErr != nil {
		return models.User{}, getUserErr
	}

	updatedUser.HouseholdId = householdID
	return ur.UpdateUser(updatedUser)
}

// DeleteHousehold - deletes a household by its ID.
func (hc HouseholdController) DeleteHousehold(householdID string) error {
	return hc.householdRepo.DeleteHousehold(householdID)
}

func (hc HouseholdController) GetCalendar(householdID string, startDate string) (models.Calendar, error) {
	return hc.calendarRepo.GetCalendar(householdID, startDate)
}

func (hc HouseholdController) UpdateCalendar(householdID string, calendar models.Calendar) (models.Calendar, error) {
	calendar.HouseholdID, _ = primitive.ObjectIDFromHex(householdID)
	updatedCalendar, err := hc.calendarRepo.UpdateCalendar(calendar)
	return updatedCalendar, err
}

func (hc HouseholdController) CreateCalendar(startDate string, householdID string) (models.Calendar, error) {
	recipes, err := hc.rc.GetRandomRecipes(7)
	if err != nil {
		return models.Calendar{}, errors.New("could not generate new calendar")
	}

	calendar := models.Calendar{
		StartDate: startDate,
		Sunday:    recipes[0],
		Monday:    recipes[1],
		Tuesday:   recipes[2],
		Wednesday: recipes[3],
		Thursday:  recipes[4],
		Friday:    recipes[5],
		Saturday:  recipes[6],
	}
	calendar.HouseholdID, _ = primitive.ObjectIDFromHex(householdID)

	return hc.calendarRepo.CreateCalendar(calendar)
}
