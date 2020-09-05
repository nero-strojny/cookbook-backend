package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"server/controller"
	"server/models"
	"server/router"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dbPointer = flag.String("DB_STRING", "", "Database connection string")
var envPointer = flag.String("ENV", "", "Environment string")

var defaultRecipe = models.Recipe{
	RecipeName:      "recipe",
	Private:         false,
	CreatedDate:     "2020.08.16 22:20:53",
	LastUpdatedDate: "2020.08.16 22:20:53",
	Author:          "testUser",
	Rating:          5,
	Servings:        2,
	Calories:        500,
	PrepTime:        5,
	CookTime:        5,
}

func Router() *mux.Router {
	r := router.Router()
	controller.SetClients(*dbPointer, *envPointer)
	return r
}

func createRecipe(recipe models.Recipe) *httptest.ResponseRecorder {
	jsonRecipe, _ := json.Marshal(recipe)
	request, _ := http.NewRequest("POST", "/api/recipe", bytes.NewBuffer(jsonRecipe))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	return response
}

func deleteRecipe(recipeID string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("DELETE", "/api/recipe/"+recipeID, nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	return response
}

func TestCreate(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("1")
	defaultRecipe.RecipeID = id
	response := createRecipe(defaultRecipe)
	assert.Equal(t, 201, response.Code, "OK response is expected")
	deleteRecipe("1")
}

// func TestCreateWithInvalidFields(t *testing.T) {
// 	id, _ := primitive.ObjectIDFromHex("2")
// 	defaultRecipe.RecipeID = id
// 	response := createRecipe(defaultRecipe)
// 	assert.Equal(t, 201, response.Code, "OK response is expected")
// 	deleteRecipe("1")
// }

// func TestGet(t *testing.T) {

// }

// func TestGetWithUnknownRecipe(t *testing.T) {

// }

// func TestGetAll(t *testing.T) {

// }

// func TestSearch(t *testing.T) {

// }

// func TestUpdate(t *testing.T) {

// }

// func TestDelete(t *testing.T) {

// }

// func TestDeleteWithUnknownRecipe(t *testing.T) {

// }
