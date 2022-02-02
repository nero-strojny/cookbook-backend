package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Recipe is used to marshal and unmarshal from a document database
type Recipe struct {
	RecipeID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RecipeName      string             `json:"recipeName,omitempty"`
	CreatedDate     string             `json:"createdDate,omitempty"`
	LastUpdatedDate string             `json:"lastUpdatedDate,omitempty"`
	Ingredients     []Ingredient       `json:"ingredients,omitempty"`
	Author          string             `json:"author,omitempty"`
	PrepTime        int                `json:"prepTime"`
	CookTime        int                `json:"cookTime"`
	Steps           []step             `json:"steps,omitempty"`
	Tags            []string           `json:"tags,omitempty"`
	Servings        int                `json:"servings"`
	Calories        int                `json:"calories"`
	UserName        string             `json:"userName,omitempty"`
	Private         bool               `json:"private,omitempty"`
}

// Ingredient is a component of a recipe consisting of the name, amount, and the measurement for that amount (cups, tbsp, lbs, etc)
type Ingredient struct {
	IngredientID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Amount       float32            `json:"amount,omitempty"`
	Measurement  string             `json:"measurement,omitempty"`
	Category     string             `json:"category,omitempty"`
}

// Step is what to do in order for a recipe
type step struct {
	Number int    `json:"number,omitempty"`
	Text   string `json:"text,omitempty"`
}

// User is the data representation of a user
type User struct {
	UserID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName     string             `json:"userName,omitempty"`
	PasswordHash string             `json:"passwordHash,omitempty"`
	AccessToken  string             `json:"accessToken,omitempty"`
	ExpiryDate   string             `json:"expiryDate,omitempty"`
	UserType     string             `json:"userType,omitempty"`
	Email        string             `json:"email,omitempty"`
	HouseholdId  string             `json:"householdId,omitempty"`
}

// RequestedUser is what is needed to create a user
type RequestedUser struct {
	UserName      string `json:"userName,omitempty"`
	Password      string `json:"password,omitempty"`
	UserType      string `json:"userType,omitempty"`
	Email         string `json:"email,omitempty"`
	AgreedToTerms bool   `json:"agreedToTerms,omitempty"`
}

// AuthData is the authentication information for a user
type AuthData struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// AccessToken is the authentication information for a user
type AccessToken struct {
	AccessToken string `json:"accessToken,omitempty"`
}

// UpdatedPassword contains the updated password information
type UpdatedPassword struct {
	UserName        string `json:"userName,omitempty"`
	CurrentPassword string `json:"currentPassword,omitempty"`
	NewPassword     string `json:"newPassword,omitempty"`
}

// PaginatedRequest
type PaginatedRecipeRequest struct {
	PageSize    int64  `json:"pageSize,omitempty"`
	PageCount   int    `json:"pageCount,omitempty"`
	QueryRecipe Recipe `json:"queryRecipe,omitempty"`
}

// PaginatedResponse
type PaginatedRecipeResponse struct {
	PageSize        int64    `json:"pageSize,omitempty"`
	PageCount       int      `json:"pageCount,omitempty"`
	NumberOfRecipes int64    `json:"numberOfRecipes,omitempty"`
	Recipes         []Recipe `json:"recipes,omitempty"`
}

type Basket struct {
	Produce  []string `json:"produce,omitempty"`
	Protein  []string `json:"protein,omitempty"`
	Pantry   []string `json:"pantry,omitempty"`
	Dairy    []string `json:"dairy,omitempty"`
	Alcohol  []string `json:"alcohol,omitempty"`
	UserName string   `json:"userName"`
}

type Config struct {
	ConnectionString string `json:"connectionString"`
	EmailPassword    string `json:"emailPassword"`
}

type HealthStatus struct {
	DB string `json:"database"`
}

type Household struct {
	HouseholdID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	HouseholdName   string             `json:"householdName,omitempty"`
	HeadOfHousehold string             `json:"headOfHousehold,omitempty"`
}

type RequestedHouseholdUpdate struct {
	UserIdToAdd string `json:"userIdToAdd,omitempty"`
}
