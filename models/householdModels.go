package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Household struct {
	HouseholdID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	HouseholdName   string             `json:"householdName,omitempty"`
	HeadOfHousehold string             `json:"headOfHousehold,omitempty"`
}

type RequestedHouseholdUpdate struct {
	UserIdToAdd string `json:"userIdToAdd,omitempty"`
}

type Calendar struct {
	CalendarID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	HouseholdID primitive.ObjectID `json:"householdID,omitempty" bson:"householdID,omitempty"`
	StartDate   string             `json:"startDate,omitempty"`
	Monday      Recipe             `json:"monday,omitempty" bson:"monday,omitempty"`
	Tuesday     Recipe             `json:"tuesday,omitempty" bson:"tuesday,omitempty"`
	Wednesday   Recipe             `json:"wednesday,omitempty" bson:"wednesday,omitempty"`
	Thursday    Recipe             `json:"thursday,omitempty" bson:"thursday,omitempty"`
	Friday      Recipe             `json:"friday,omitempty" bson:"friday,omitempty"`
	Saturday    Recipe             `json:"saturday,omitempty" bson:"saturday,omitempty"`
	Sunday      Recipe             `json:"sunday,omitempty" bson:"sunday,omitempty"`
}
