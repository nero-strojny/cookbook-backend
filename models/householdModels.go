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
	Monday      primitive.ObjectID `json:"monday,omitempty" bson:"monday,omitempty"`
	Tuesday     primitive.ObjectID `json:"tuesday,omitempty" bson:"tuesday,omitempty"`
	Wednesday   primitive.ObjectID `json:"wednesday,omitempty" bson:"wednesday,omitempty"`
	Thursday    primitive.ObjectID `json:"thursday,omitempty" bson:"thursday,omitempty"`
	Friday      primitive.ObjectID `json:"friday,omitempty" bson:"friday,omitempty"`
	Saturday    primitive.ObjectID `json:"saturday,omitempty" bson:"saturday,omitempty"`
	Sunday      primitive.ObjectID `json:"sunday,omitempty" bson:"sunday,omitempty"`
}
