package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"server/controller"
	"server/models"
	"server/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

var recipeRouter = router.Router()
var recipeAdminToken string

var defaultRecipe = models.Recipe{
	RecipeName:      "recipe",
	Private:         false,
	CreatedDate:     "2020.08.16 22:20:53",
	LastUpdatedDate: "2020.08.16 22:20:53",
	Author:          "someAuthor",
	Rating:          5,
	Servings:        2,
	Calories:        500,
	PrepTime:        5,
	CookTime:        5,
	UserName:        "testAdminUser",
}

func createRecipe(inputRecipe models.Recipe, accessToken string) (*httptest.ResponseRecorder, models.Recipe) {
	outputRecipe := models.Recipe{}
	jsonRecipe, _ := json.Marshal(inputRecipe)
	request, _ := http.NewRequest("POST", "/api/recipe", bytes.NewBuffer(jsonRecipe))
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &outputRecipe)
	return response, outputRecipe
}

func getRecipeByID(recipeID string, accessToken string) (*httptest.ResponseRecorder, models.Recipe) {
	recipe := models.Recipe{}
	request, _ := http.NewRequest("GET", "/api/recipe/"+recipeID, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &recipe)
	return response, recipe
}

func getAllRecipes(accessToken string) (*httptest.ResponseRecorder, []models.Recipe) {
	recipes := []models.Recipe{}
	request, _ := http.NewRequest("GET", "/api/recipes", nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &recipes)
	return response, recipes
}

func searchRecipeByName(recipeName string, accessToken string) (*httptest.ResponseRecorder, models.Recipe) {
	var searchRecipe = models.Recipe{
		RecipeName: recipeName,
	}
	recipe := models.Recipe{}
	jsonRecipe, _ := json.Marshal(searchRecipe)
	request, _ := http.NewRequest("POST", "/api/recipe/search", bytes.NewBuffer(jsonRecipe))
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &recipe)
	return response, recipe
}

func updateRecipe(inputRecipe models.Recipe, accessToken string) (*httptest.ResponseRecorder, models.Recipe) {
	outputRecipe := models.Recipe{}
	jsonRecipe, _ := json.Marshal(inputRecipe)
	request, _ := http.NewRequest("PUT", "/api/recipe/"+inputRecipe.RecipeID.Hex(), bytes.NewBuffer(jsonRecipe))
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &outputRecipe)
	return response, outputRecipe
}

func deleteRecipe(recipeID string, accessToken string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("DELETE", "/api/recipe/"+recipeID, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	return response
}

func TestRecipeSetUp(t *testing.T) {
	// generate a token with admin user data
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)
	recipeAdminToken = accessTokenObject.AccessToken
}

func TestCreateRecipe(t *testing.T) {
	// create a recipe
	response, recipe := createRecipe(defaultRecipe, recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 201, response.Code, "OK response is expected")
	recipeFieldsAreExpected(t, recipe, defaultRecipe)

	// cleanup
	deleteRecipe(recipe.RecipeID.Hex(), recipeAdminToken)
}

func TestCreateRecipeWithInvalidFields(t *testing.T) {
	// create a recipe with an empty recipe Name
	invalidRecipe := models.Recipe{}
	invalidRecipe = defaultRecipe
	invalidRecipe.RecipeName = ""
	response, _ := createRecipe(invalidRecipe, recipeAdminToken)

	// assert the correct status code
	assert.Equal(t, 400, response.Code, "Invalid input response is expected")
}

func TestGetRecipe(t *testing.T) {
	// setUp, create a recipe
	_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)

	// make a getRecipe for the recipe we just created
	getResponse, getRecipe := getRecipeByID(createdRecipe.RecipeID.Hex(), recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	recipeFieldsAreExpected(t, getRecipe, defaultRecipe)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)
}

func TestGetRecipeWithUnknownRecipe(t *testing.T) {
	// try to get a recipe with an unknown id
	getResponse, _ := getRecipeByID("unknownRecipeId", recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")
}

func TestGetAll(t *testing.T) {
	// setUp, create some recipes
	_, createdRecipe1 := createRecipe(defaultRecipe, recipeAdminToken)
	_, createdRecipe2 := createRecipe(defaultRecipe, recipeAdminToken)

	// make a getRecipe for the recipe we just created
	getResponse, recipes := getAllRecipes(recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Contains(t, recipes, createdRecipe1, "The first created recipe should be in the results")
	assert.Contains(t, recipes, createdRecipe2, "The second created recipe should be in the results")

	// cleanup
	deleteRecipe(createdRecipe1.RecipeID.Hex(), recipeAdminToken)
	deleteRecipe(createdRecipe2.RecipeID.Hex(), recipeAdminToken)
}

func TestPaginatedRecipes(t *testing.T) {
	// setUp, create some recipes
	var createdRecipes []models.Recipe
	for i := 0; i < 10; i++ {
		defaultRecipe.RecipeName = fmt.Sprint("recipe", i)
		_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)
		createdRecipes = append(createdRecipes, createdRecipe)
	}

	var paginatedRequest = models.PaginatedRequest{
		PageSize:  2,
		PageCount: 3,
	}

	// make a getRecipe for the recipe we just created
	recipes, _ := controller.PaginatedRecipes(paginatedRequest)

	// assert the correct status code and body
	assert.Equal(t, 2, len(recipes), "Correct Length is expected")
	assert.Equal(t, "recipe4",	recipes[0].RecipeName, "Correct Length is expected")
	assert.Equal(t, "recipe5",	recipes[0x.
	?].RecipeName, "Correct Length is expected")

	// cleanup
	for i := 0; i < 10; i++ {
		deleteRecipe(createdRecipes[i].RecipeID.Hex(), recipeAdminToken)
	}
}

func TestSearchRecipe(t *testing.T) {
	// set up, create a recipe with a unique name
	specificNamedRecipe := defaultRecipe
	specificNamedRecipe.RecipeName = "Specific Recipe Name"
	_, createdRecipe := createRecipe(specificNamedRecipe, recipeAdminToken)

	// make a search for the recipe we just created
	searchResponse, resultRecipe := searchRecipeByName("Specific Recipe Name", recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	recipeFieldsAreExpected(t, resultRecipe, createdRecipe)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)
}

func TestUpdateRecipe(t *testing.T) {
	// setUp, create a recipe
	_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)

	createdRecipe.RecipeName = "updatedRecipe"
	updatedResponse, updatedRecipe := updateRecipe(createdRecipe, recipeAdminToken)

	// make a getRecipe for the recipe we just updated
	getResponse, getRecipe := getRecipeByID(createdRecipe.RecipeID.Hex(), recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Equal(t, 200, updatedResponse.Code, "OK response is expected")
	assert.Equal(t, "updatedRecipe", updatedRecipe.RecipeName)
	assert.Equal(t, "updatedRecipe", getRecipe.RecipeName)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)
}

func TestDeleteRecipe(t *testing.T) {
	// set up, create a recipe
	_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)

	// delete a recipe
	deleteResponse := deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 204, deleteResponse.Code, "OK response is expected")

	//try to get the recipe again to ensure it is deleted
	getResponse, _ := getRecipeByID(createdRecipe.RecipeID.Hex(), recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")

}

func TestCreateRecipeWithInvalidToken(t *testing.T) {
	// create a recipe
	response, recipe := createRecipe(defaultRecipe, "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, response.Code, "Unauthorized response is expected")

	// cleanup
	deleteRecipe(recipe.RecipeID.Hex(), recipeAdminToken)

}

func TestGetRecipeWithInvalidToken(t *testing.T) {
	// setUp, create a recipe
	_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)

	// make a getRecipe for the recipe we just created
	getResponse, _ := getRecipeByID(createdRecipe.RecipeID.Hex(), "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, getResponse.Code, "Unauthorized response is expected")

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)

}

func TestGetAllRecipesWithInvalidToken(t *testing.T) {
	// setUp, create some recipes
	_, createdRecipe1 := createRecipe(defaultRecipe, recipeAdminToken)
	_, createdRecipe2 := createRecipe(defaultRecipe, recipeAdminToken)

	// make a getRecipe for the recipe we just created
	getResponse, _ := getAllRecipes("invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, getResponse.Code, "Unauthorized response is expected")

	// cleanup
	deleteRecipe(createdRecipe1.RecipeID.Hex(), recipeAdminToken)
	deleteRecipe(createdRecipe2.RecipeID.Hex(), recipeAdminToken)

}

func TestSearchRecipeWithInvalidToken(t *testing.T) {
	// set up, create a recipe with a unique name
	specificNamedRecipe := defaultRecipe
	specificNamedRecipe.RecipeName = "Specific Recipe Name"
	_, createdRecipe := createRecipe(specificNamedRecipe, recipeAdminToken)

	// make a search for the recipe we just created
	searchResponse, _ := searchRecipeByName("Specific Recipe Name", "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, searchResponse.Code, "Unauthorized response is expected")

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)

}

func TestDeleteRecipeWithUnknownID(t *testing.T) {
	// delete a recipe
	deleteResponse := deleteRecipe("unknownRecipeID", recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 404, deleteResponse.Code, "Unauthorized response is expected")
}

func TestUpdateRecipeWithOtherUser(t *testing.T) {
	//setUp create a recipe and a user
	createUser(defaultNonAdminUser, recipeAdminToken)
	accessTokenObject := models.AccessToken{}
	userResponse := generateUserToken(defaultNonAdminAuthData)
	body, _ := ioutil.ReadAll(userResponse.Body)
	json.Unmarshal(body, &accessTokenObject)
	_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)

	// update the recipe
	createdRecipe.RecipeName = "updatedRecipe"
	updatedResponse, _ := updateRecipe(createdRecipe, accessTokenObject.AccessToken)

	// assert the correct status code and body
	assert.Equal(t, 401, updatedResponse.Code, "Not Found response is expected")

	//cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)
	deleteUser(defaultNonAdminUser.UserName, recipeAdminToken)
}

func TestDeleteRecipeWithOtherUser(t *testing.T) {
	//setUp create a recipe and a user
	createUser(defaultNonAdminUser, recipeAdminToken)
	accessTokenObject := models.AccessToken{}
	userResponse := generateUserToken(defaultNonAdminAuthData)
	body, _ := ioutil.ReadAll(userResponse.Body)
	json.Unmarshal(body, &accessTokenObject)
	_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)

	// delete the recipe
	deleteResponse := deleteRecipe(createdRecipe.RecipeID.Hex(), accessTokenObject.AccessToken)

	// assert the correct status code and body
	assert.Equal(t, 401, deleteResponse.Code, "Not Found response is expected")

	//cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)
	deleteUser(defaultNonAdminUser.UserName, recipeAdminToken)
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
