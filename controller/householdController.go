package controller

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/db"
	"server/models"
)

type HouseholdControl interface {
	CreateHousehold(household models.Household, user models.User, repository db.HouseholdCreator) (models.Household, error)
	GetHousehold(householdID string, repository db.HouseholdGetter) (models.Household, error)
	AddUserToHousehold(householdID string, username string, ur db.UserGetterUpdater) (models.User, error)
	DeleteHousehold(householdID string, repository db.HouseholdDeleter) error
	GetCalendar(householdID string, startDate string, repository db.CalendarGetter) (models.Calendar, error)
	UpdateCalendar(householdID string, calendar models.Calendar, updater db.CalendarGetterUpdater) (models.Calendar, error)
	CreateCalendar(calendar models.Calendar, creator db.CalendarCreator) (models.Calendar, error)
}

type HouseholdController struct {
}

func NewHouseholdController() HouseholdController {
	return HouseholdController{}
}

//CreateHousehold creates a new household
func (hc HouseholdController) CreateHousehold(household models.Household, user models.User, repository db.HouseholdCreator) (models.Household, error) {
	return repository.CreateHousehold(household, user)
}

//GetHousehold - gets households by its ID
func (hc HouseholdController) GetHousehold(householdID string, repository db.HouseholdGetter) (models.Household, error) {
	return repository.GetHousehold(householdID)
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
func (hc HouseholdController) DeleteHousehold(householdID string, repository db.HouseholdDeleter) error {
	return repository.DeleteHousehold(householdID)
}

func (hc HouseholdController) GetCalendar(householdID string, startDate string, repository db.CalendarGetter) (models.Calendar, error) {
	return repository.GetCalendar(householdID, startDate)
}

func (hc HouseholdController) UpdateCalendar(householdID string, calendar models.Calendar, updater db.CalendarGetterUpdater) (models.Calendar, error) {
	calendar.HouseholdID, _ = primitive.ObjectIDFromHex(householdID)
	updatedCalendar, err := updater.UpdateCalendar(calendar)
	return updatedCalendar, err
}

func (hc HouseholdController) CreateCalendar(calendar models.Calendar, creator db.CalendarCreator) (models.Calendar, error) {
	return creator.CreateCalendar(calendar)
}
