package models
import "go.mongodb.org/mongo-driver/bson/primitive"

// Recipe is used to marshal and unmarshal from a document database
type Recipe struct {
  RecipeID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  CreatedDate   string             `json:"createdDate,omitempty"`
  LastUpdatedDate string               `json:"lastUpdatedDate,omitempty"`
  Ingredients []ingredient `json:"ingredients,omitempty"`
  Author string `json:"author,omitempty"`
  PrepTime int `json:"prepTime,omitempty"`
  CookTime int `json:"cookTime,omitempty"`
  Steps []step `json:"steps,omitempty"`

}

// Ingredient is a component of a recipe consisting of the name, amount, and the measurement for that amount (cups, tbsp, lbs, etc)
type ingredient struct {
  Name string `json:"name,omitempty"`
  Amount float32 `json:"amount,omitempty"`
  Measurement string `json:"measurement,omitempty"`
}

// Step is what to do in order for a recipe
type step struct {
  Number int `json:"number,omitempty"`
  Text string `json:"text,omitEmpty"`
}