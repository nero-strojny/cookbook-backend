package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"server/models"
	"server/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

var recipeRouter = router.Router()

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

func createRecipe(recipe models.Recipe) *httptest.ResponseRecorder {
	jsonRecipe, _ := json.Marshal(recipe)
	request, _ := http.NewRequest("POST", "/api/recipe", bytes.NewBuffer(jsonRecipe))
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	return response
}

func deleteRecipe(recipeID string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("DELETE", "/api/recipe/"+recipeID, nil)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	return response
}

func TestCreateRecipe(t *testing.T) {
	// create a recipe
	recipe := models.Recipe{}
	response := createRecipe(defaultRecipe)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &recipe)

	// assert the correct status code and body
	assert.Equal(t, 201, response.Code, "OK response is expected")
	recipeFieldsAreExpected(t, recipe, defaultRecipe)

	// cleanup
	deleteRecipe(recipe.RecipeID.Hex())
}

func TestCreateRecipeWithInvalidFields(t *testing.T) {
	// create a recipe with an empty recipe Name
	invalidRecipe := models.Recipe{}
	invalidRecipe = defaultRecipe
	invalidRecipe.RecipeName = ""
	response := createRecipe(invalidRecipe)

	// assert the correct status code
	assert.Equal(t, 400, response.Code, "Invalid input response is expected")
}

func TestGetRecipe(t *testing.T) {
	// setUp, create a recipe
	createdRecipe := models.Recipe{}
	createResponse := createRecipe(defaultRecipe)
	createBody, _ := ioutil.ReadAll(createResponse.Body)
	json.Unmarshal(createBody, &createdRecipe)

	// make a getRecipe for the recipe we just created
	getRecipe := models.Recipe{}
	getRequest, _ := http.NewRequest("GET", "/api/recipe/"+createdRecipe.RecipeID.Hex(), nil)
	getResponse := httptest.NewRecorder()
	recipeRouter.ServeHTTP(getResponse, getRequest)
	getBody, _ := ioutil.ReadAll(getResponse.Body)
	json.Unmarshal(getBody, &getRecipe)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	recipeFieldsAreExpected(t, getRecipe, defaultRecipe)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex())
}

func TestGetRecipeWithUnknownRecipe(t *testing.T) {
	// try to get a recipe with an unknown id
	getRecipe := models.Recipe{}
	getRequest, _ := http.NewRequest("GET", "/api/recipe/unknownRecipeId", nil)
	getResponse := httptest.NewRecorder()
	recipeRouter.ServeHTTP(getResponse, getRequest)
	getBody, _ := ioutil.ReadAll(getResponse.Body)
	json.Unmarshal(getBody, &getRecipe)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")
}

func TestGetAll(t *testing.T) {
	// setUp, create some recipes
	createdRecipe1 := models.Recipe{}
	createResponse1 := createRecipe(defaultRecipe)
	createBody1, _ := ioutil.ReadAll(createResponse1.Body)
	json.Unmarshal(createBody1, &createdRecipe1)
	createdRecipe2 := models.Recipe{}
	createResponse2 := createRecipe(defaultRecipe)
	createBody2, _ := ioutil.ReadAll(createResponse2.Body)
	json.Unmarshal(createBody2, &createdRecipe2)

	// make a getRecipe for the recipe we just created
	recipes := []models.Recipe{}
	getRequest, _ := http.NewRequest("GET", "/api/recipes", nil)
	getResponse := httptest.NewRecorder()
	recipeRouter.ServeHTTP(getResponse, getRequest)
	getBody, _ := ioutil.ReadAll(getResponse.Body)
	json.Unmarshal(getBody, &recipes)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Contains(t, recipes, createdRecipe1, "The first created recipe should be in the results")
	assert.Contains(t, recipes, createdRecipe2, "The second created recipe should be in the results")

	// cleanup
	deleteRecipe(createdRecipe1.RecipeID.Hex())
	deleteRecipe(createdRecipe2.RecipeID.Hex())
}

func TestSearchRecipe(t *testing.T) {
	// set up, create a recipe with a unique name
	createdRecipe := models.Recipe{}
	specificNamedRecipe := defaultRecipe
	specificNamedRecipe.RecipeName = "Specific Recipe Name"
	createResponse := createRecipe(specificNamedRecipe)
	createBody, _ := ioutil.ReadAll(createResponse.Body)
	json.Unmarshal(createBody, &createdRecipe)

	// make a search for the recipe we just created
	var searchRecipe = models.Recipe{
		RecipeName: "Specific Recipe Name",
	}
	resultRecipe := models.Recipe{}
	jsonRecipe, _ := json.Marshal(searchRecipe)
	searchRequest, _ := http.NewRequest("POST", "/api/recipe/search", bytes.NewBuffer(jsonRecipe))
	searchResponse := httptest.NewRecorder()
	recipeRouter.ServeHTTP(searchResponse, searchRequest)
	searchBody, _ := ioutil.ReadAll(searchResponse.Body)
	json.Unmarshal(searchBody, &resultRecipe)

	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	recipeFieldsAreExpected(t, resultRecipe, createdRecipe)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex())
}

func TestUpdateRecipe(t *testing.T) {
	// setUp, create a recipe
	createdRecipe := models.Recipe{}
	createResponse := createRecipe(defaultRecipe)
	createBody, _ := ioutil.ReadAll(createResponse.Body)
	json.Unmarshal(createBody, &createdRecipe)

	createdRecipe.RecipeName = "updatedRecipe"
	updatedRecipe := createdRecipe
	updatedRecipe.RecipeName = "updatedRecipe"
	jsonRecipe, _ := json.Marshal(updatedRecipe)
	updatedRequest, _ := http.NewRequest("PUT", "/api/recipe/"+createdRecipe.RecipeID.Hex(), bytes.NewBuffer(jsonRecipe))
	updatedResponse := httptest.NewRecorder()
	recipeRouter.ServeHTTP(updatedResponse, updatedRequest)
	updatedBody, _ := ioutil.ReadAll(updatedResponse.Body)
	json.Unmarshal(updatedBody, &updatedRecipe)

	// make a getRecipe for the recipe we just updated
	getRecipe := models.Recipe{}
	getRequest, _ := http.NewRequest("GET", "/api/recipe/"+createdRecipe.RecipeID.Hex(), nil)
	getResponse := httptest.NewRecorder()
	recipeRouter.ServeHTTP(getResponse, getRequest)
	getBody, _ := ioutil.ReadAll(getResponse.Body)
	json.Unmarshal(getBody, &getRecipe)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Equal(t, "updatedRecipe", getRecipe.RecipeName)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex())
}

func TestDeleteRecipe(t *testing.T) {
	// set up, create a recipe
	recipe := models.Recipe{}
	createResponse := createRecipe(defaultRecipe)
	body, _ := ioutil.ReadAll(createResponse.Body)
	json.Unmarshal(body, &recipe)

	// delete a recipe
	deleteResponse := deleteRecipe(recipe.RecipeID.Hex())

	// assert the correct status code and body
	assert.Equal(t, 204, deleteResponse.Code, "OK response is expected")

	//try to get the recipe again to ensure it is deleted
	getRequest, _ := http.NewRequest("GET", "/api/recipe/unknownRecipeId", nil)
	getResponse := httptest.NewRecorder()
	recipeRouter.ServeHTTP(getResponse, getRequest)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")

}

func TestDeleteRecipeWithUnknownID(t *testing.T) {
	// delete a recipe
	deleteResponse := deleteRecipe("unknownRecipeID")

	// assert the correct status code and body
	assert.Equal(t, 404, deleteResponse.Code, "Not Found response is expected")
}

func recipeFieldsAreExpected(t *testing.T, recipe1 models.Recipe, recipe2 models.Recipe) {
	assert.NotNilf(t, recipe1.CreatedDate, "CreatedDate should be set")
	assert.NotNilf(t, recipe1.LastUpdatedDate, "LastUpdatedDate should be set")
	assert.Equal(t, recipe1.Author, recipe2.Author, "Inputted author value expected")
	assert.Equal(t, recipe1.Rating, recipe2.Rating, "Inputted rating is expected")
	assert.Equal(t, recipe1.Servings, recipe2.Servings, "Inputted servings is expected")
	assert.Equal(t, recipe1.Calories, recipe2.Calories, "Inputted calories is expected")
	assert.Equal(t, recipe1.RecipeName, recipe2.RecipeName, "Inputted recipe name is expected")
	assert.Equal(t, recipe1.Ingredients, recipe2.Ingredients, "Inputted ingredients are expected")
	assert.Equal(t, recipe1.Steps, recipe2.Steps, "Inputted steps are expected")
	assert.Equal(t, recipe1.Private, recipe2.Private, "Inputted private setting is expected")

}
