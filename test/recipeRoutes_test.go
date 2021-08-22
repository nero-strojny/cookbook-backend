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
	Servings:        2,
	Tags:            []string{"testTag"},
	Calories:        500,
	PrepTime:        5,
	CookTime:        5,
	UserName:        "testAdminUser",
}

var defaultPaginatedRequest = models.PaginatedRecipeRequest{
	PageCount:   0,
	PageSize:    10,
	QueryRecipe: models.Recipe{},
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

func postPaginatedRecipes(paginatedRequest models.PaginatedRecipeRequest, accessToken string) (*httptest.ResponseRecorder, models.PaginatedRecipeResponse) {
	jsonRequest, _ := json.Marshal(paginatedRequest)
	paginatedResult := models.PaginatedRecipeResponse{}
	request, _ := http.NewRequest("POST", "/api/recipes", bytes.NewBuffer(jsonRequest))
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &paginatedResult)
	return response, paginatedResult
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
	getResponse, paginatedResult := postPaginatedRecipes(defaultPaginatedRequest, recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Equal(t, int64(2), paginatedResult.NumberOfRecipes, "Correct number of recipes expected")
	assert.Contains(t, paginatedResult.Recipes, createdRecipe1, "The first created recipe should be in the results")
	assert.Contains(t, paginatedResult.Recipes, createdRecipe2, "The second created recipe should be in the results")

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

	var paginatedRequest = models.PaginatedRecipeRequest{
		PageSize:  2,
		PageCount: 3,
	}

	// make a getRecipe for the recipe we just created
	getResponse, paginatedResult := postPaginatedRecipes(paginatedRequest, recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Equal(t, 2, len(paginatedResult.Recipes), "Correct Length is expected")
	assert.Equal(t, int64(10), paginatedResult.NumberOfRecipes, "Correct number of recipes expected")
	assert.Equal(t, "recipe6", paginatedResult.Recipes[0].RecipeName, "Correct Length is expected")
	assert.Equal(t, "recipe7", paginatedResult.Recipes[1].RecipeName, "Correct Length is expected")

	// cleanup
	for i := 0; i < 10; i++ {
		deleteRecipe(createdRecipes[i].RecipeID.Hex(), recipeAdminToken)
	}
}

func TestPaginatedOutOfBoundsRecipes(t *testing.T) {
	// setUp, create some recipes
	var createdRecipes []models.Recipe
	for i := 0; i < 10; i++ {
		defaultRecipe.RecipeName = fmt.Sprint("recipe", i)
		_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)
		createdRecipes = append(createdRecipes, createdRecipe)
	}

	var paginatedRequest = models.PaginatedRecipeRequest{
		PageSize:  5,
		PageCount: 3,
	}

	// make a getRecipe for the recipe we just created
	getResponse, paginatedResult := postPaginatedRecipes(paginatedRequest, recipeAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Equal(t, 5, len(paginatedResult.Recipes), "Correct Length is expected")
	assert.Equal(t, int64(10), paginatedResult.NumberOfRecipes, "Correct number of recipes expected")
	assert.Equal(t, "recipe5", paginatedResult.Recipes[0].RecipeName, "Correct Length is expected")
	assert.Equal(t, "recipe6", paginatedResult.Recipes[1].RecipeName, "Correct Length is expected")
	assert.Equal(t, "recipe7", paginatedResult.Recipes[2].RecipeName, "Correct Length is expected")
	assert.Equal(t, "recipe8", paginatedResult.Recipes[3].RecipeName, "Correct Length is expected")
	assert.Equal(t, "recipe9", paginatedResult.Recipes[4].RecipeName, "Correct Length is expected")

	// cleanup
	for i := 0; i < 10; i++ {
		deleteRecipe(createdRecipes[i].RecipeID.Hex(), recipeAdminToken)
	}
}

func TestRandomRecipes(t *testing.T) {
	// setUp, create some recipes
	var createdRecipes []models.Recipe
	for i := 0; i < 20; i++ {
		defaultRecipe.RecipeName = fmt.Sprint("recipe", i)
		_, createdRecipe := createRecipe(defaultRecipe, recipeAdminToken)
		createdRecipes = append(createdRecipes, createdRecipe)
	}

	randomRecipes, _ := controller.GetRandomRecipes(5)
	assert.Equal(t, 5, len(randomRecipes), "Inputted Length Expected")

	// cleanup
	for i := 0; i < 20; i++ {
		deleteRecipe(createdRecipes[i].RecipeID.Hex(), recipeAdminToken)
	}

}

func TestSearchRecipeByName(t *testing.T) {
	// set up, create a recipe with a unique name
	specificNamedRecipe := defaultRecipe
	specificNamedRecipe.RecipeName = "Specific Recipe Name"
	_, createdRecipe := createRecipe(specificNamedRecipe, recipeAdminToken)

	// make a search for the recipe we just created
	var paginatedRequest = models.PaginatedRecipeRequest{
		PageSize:  10,
		PageCount: 0,
		QueryRecipe: models.Recipe{
			RecipeName: "Specific Recipe Name",
		},
	}
	searchResponse, paginatedResult := postPaginatedRecipes(paginatedRequest, recipeAdminToken)

	resultRecipes := paginatedResult.Recipes

	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	recipeFieldsAreExpected(t, resultRecipes[0], createdRecipe)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)
}

func TestSearchRecipeByTag(t *testing.T) {
	// set up, create a recipe
	specificNamedRecipe := defaultRecipe
	specificNamedRecipe.Tags = []string{"tag1", "tag2"}
	_, createdRecipe := createRecipe(specificNamedRecipe, recipeAdminToken)

	// make a search for the recipe we just created
	var paginatedRequest = models.PaginatedRecipeRequest{
		PageSize:  10,
		PageCount: 0,
		QueryRecipe: models.Recipe{
			Tags: []string{"tag1"},
		},
	}
	searchResponse, paginatedResult := postPaginatedRecipes(paginatedRequest, recipeAdminToken)

	resultRecipes := paginatedResult.Recipes
	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	recipeFieldsAreExpected(t, resultRecipes[0], createdRecipe)

	// cleanup
	deleteRecipe(createdRecipe.RecipeID.Hex(), recipeAdminToken)
}

func TestNoFoundSearchRecipe(t *testing.T) {
	// set up, create a recipe with a unique name
	specificNamedRecipe1 := defaultRecipe
	specificNamedRecipe1.RecipeName = "Specific Recipe Name"
	_, createdRecipe1 := createRecipe(specificNamedRecipe1, recipeAdminToken)
	specificNamedRecipe2 := defaultRecipe
	specificNamedRecipe2.RecipeName = "Another specific recipe name"
	_, createdRecipe2 := createRecipe(specificNamedRecipe2, recipeAdminToken)

	// make a search that will not return anything
	var paginatedRequest = models.PaginatedRecipeRequest{
		PageSize:  10,
		PageCount: 0,
		QueryRecipe: models.Recipe{
			RecipeName: "Something else",
		},
	}

	searchResponse, paginatedResult := postPaginatedRecipes(paginatedRequest, recipeAdminToken)
	resultRecipes := paginatedResult.Recipes

	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	assert.Equal(t, int64(0), paginatedResult.NumberOfRecipes, "Correct number of recipes expected")
	assert.Equal(t, 0, len(resultRecipes))

	// cleanup
	deleteRecipe(createdRecipe1.RecipeID.Hex(), recipeAdminToken)
	deleteRecipe(createdRecipe2.RecipeID.Hex(), recipeAdminToken)
}

func TestPartialNameSearchRecipe(t *testing.T) {
	// set up, create a recipe with a unique name
	specificNamedRecipe1 := defaultRecipe
	specificNamedRecipe1.RecipeName = "Specific Recipe Name"
	_, createdRecipe1 := createRecipe(specificNamedRecipe1, recipeAdminToken)
	specificNamedRecipe2 := defaultRecipe
	specificNamedRecipe2.RecipeName = "Another specific recipe name"
	_, createdRecipe2 := createRecipe(specificNamedRecipe2, recipeAdminToken)
	specificNamedRecipe3 := defaultRecipe
	specificNamedRecipe3.RecipeName = "Mismatching name"
	_, createdRecipe3 := createRecipe(specificNamedRecipe3, recipeAdminToken)

	// make a search for the recipe we just created
	var paginatedRequest = models.PaginatedRecipeRequest{
		PageSize:  10,
		PageCount: 0,
		QueryRecipe: models.Recipe{
			RecipeName: "Specific Recipe Name",
		},
	}
	searchResponse, paginatedResult := postPaginatedRecipes(paginatedRequest, recipeAdminToken)

	resultRecipes := paginatedResult.Recipes

	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	assert.Equal(t, int64(2), paginatedResult.NumberOfRecipes, "Correct number of recipes expected")
	recipeFieldsAreExpected(t, resultRecipes[0], createdRecipe1)
	recipeFieldsAreExpected(t, resultRecipes[1], createdRecipe2)

	// cleanup
	deleteRecipe(createdRecipe1.RecipeID.Hex(), recipeAdminToken)
	deleteRecipe(createdRecipe2.RecipeID.Hex(), recipeAdminToken)
	deleteRecipe(createdRecipe3.RecipeID.Hex(), recipeAdminToken)
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
	getResponse, _ := postPaginatedRecipes(defaultPaginatedRequest, "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, getResponse.Code, "Unauthorized response is expected")

	// cleanup
	deleteRecipe(createdRecipe1.RecipeID.Hex(), recipeAdminToken)
	deleteRecipe(createdRecipe2.RecipeID.Hex(), recipeAdminToken)

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
	assert.Equal(t, recipe1.Servings, recipe2.Servings, "Inputted servings is expected")
	assert.Equal(t, recipe1.Tags, recipe2.Tags, "Inputted tags is expected")
	assert.Equal(t, recipe1.Calories, recipe2.Calories, "Inputted calories is expected")
	assert.Equal(t, recipe1.RecipeName, recipe2.RecipeName, "Inputted recipe name is expected")
	assert.Equal(t, recipe1.Ingredients, recipe2.Ingredients, "Inputted ingredients are expected")
	assert.Equal(t, recipe1.Steps, recipe2.Steps, "Inputted steps are expected")
	assert.Equal(t, recipe1.Private, recipe2.Private, "Inputted private setting is expected")

}
