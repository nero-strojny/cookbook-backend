package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Recipe is used to marshal and unmarshal from a document database
type Recipe struct {
	RecipeID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RecipeName      string             `json:"recipeName,omitempty"`
	CreatedDate     string             `json:"createdDate,omitempty"`
	LastUpdatedDate string             `json:"lastUpdatedDate,omitempty"`
	Ingredients     []ingredient       `json:"ingredients,omitempty"`
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
type ingredient struct {
	Name        string  `json:"name,omitempty"`
	Amount      float32 `json:"amount,omitempty"`
	Measurement string  `json:"measurement,omitempty"`
}

// Step is what to do in order for a recipe
type step struct {
	Number int    `json:"number,omitempty"`
	Text   string `json:"text,omitEmpty"`
}

// UserData is data pertaining to the user's physical characteristics
type UserData struct {
	UserID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"userName,omitempty"`
}

// User is the data representation of a user
type User struct {
	UserID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName     string             `json:"userName,omitempty"`
	PasswordHash string             `json:"passwordHash,omitEmpty"`
	AccessToken  string             `json:"accessToken,omitEmpty"`
	ExpiryDate   string             `json:"expiryDate,omitempty"`
	UserType     string             `json:"userType,omitempty"`
}

// RequestedUser is what is needed to create a user
type RequestedUser struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitEmpty"`
	UserType string `json:"userType,omitempty"`
}

// AuthData is the authentication information for a user
type AuthData struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// AccessToken is the authentication information for a user
type AccessToken struct {
	AccessToken string `json:"accessToken,omitEmpty"`
}

// UpdatedPassword contains the updated password information
type UpdatedPassword struct {
	UserName        string `json:"userName,omitEmpty"`
	CurrentPassword string `json:"currentPassword,omitEmpty"`
	NewPassword     string `json:"newPassword,omitEmpty"`
}

// PaginatedRequest
type PaginatedRecipeRequest struct {
	PageSize  int64 `json:"pageSize,omitEmpty"`
	PageCount int   `json:"pageCount,omitEmpty"`
}

// PaginatedResponse
type PaginatedRecipeResponse struct {
	PageSize        int64    `json:"pageSize,omitEmpty"`
	PageCount       int      `json:"pageCount,omitEmpty"`
	NumberOfRecipes int64    `json:"numberOfRecipes,omitEmpty"`
	Recipes         []Recipe `json:"recipes,omitEmpty"`
}
