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

var ingredientRouter = router.Router()
var ingredientAdminToken string

var defaultIngredient = models.Ingredient{
	Name: "Ingredient",
}

func createIngredient(inputIngredient models.Ingredient, accessToken string) (*httptest.ResponseRecorder, models.Ingredient) {
	outputIngredient := models.Ingredient{}
	jsonIngredient, _ := json.Marshal(inputIngredient)
	request, _ := http.NewRequest("POST", "/api/ingredient", bytes.NewBuffer(jsonIngredient))
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	ingredientRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &outputIngredient)
	return response, outputIngredient
}

func getIngredientByID(IngredientID string, accessToken string) (*httptest.ResponseRecorder, models.Ingredient) {
	Ingredient := models.Ingredient{}
	request, _ := http.NewRequest("GET", "/api/ingredient/"+IngredientID, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	ingredientRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &Ingredient)
	return response, Ingredient
}

func queryIngredients(prefixIngredient string, accessToken string) (*httptest.ResponseRecorder, []models.Ingredient) {
	ingredientsResult := []models.Ingredient{}
	request, _ := http.NewRequest("GET", "/api/ingredients?prefixIngredient="+prefixIngredient, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	ingredientRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &ingredientsResult)
	return response, ingredientsResult
}

func deleteIngredient(ingredientID string, accessToken string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("DELETE", "/api/ingredient/"+ingredientID, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	recipeRouter.ServeHTTP(response, request)
	return response
}

func TestIngredientSetUp(t *testing.T) {
	// generate a token with admin user data
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)
	ingredientAdminToken = accessTokenObject.AccessToken
}

func TestCreateIngredient(t *testing.T) {
	// create an ingredient
	response, Ingredient := createIngredient(defaultIngredient, ingredientAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 201, response.Code, "OK response is expected")
	IngredientFieldsAreExpected(t, Ingredient, defaultIngredient)

	// cleanup
	deleteIngredient(Ingredient.IngredientID.Hex(), ingredientAdminToken)
}

func TestGetIngredient(t *testing.T) {
	// setUp, create an ingredient
	_, createdIngredient := createIngredient(defaultIngredient, ingredientAdminToken)

	// make a getIngredient for the Ingredient we just created
	getResponse, getIngredient := getIngredientByID(createdIngredient.IngredientID.Hex(), ingredientAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	IngredientFieldsAreExpected(t, getIngredient, defaultIngredient)

	// cleanup
	deleteIngredient(createdIngredient.IngredientID.Hex(), ingredientAdminToken)
}

func TestGetIngredientWithUnknownIngredient(t *testing.T) {
	// try to get an ingredient with an unknown id
	getResponse, _ := getIngredientByID("unknownIngredientId", ingredientAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")
}

func TestSearchIngredientByPrefix(t *testing.T) {
	// set up, create an ingredient with a unique name
	specificNamedIngredient := defaultIngredient
	specificNamedIngredient.Name = "Specific Ingredient Name"
	_, createdIngredient := createIngredient(specificNamedIngredient, ingredientAdminToken)

	// make a search for the Ingredient we just created
	searchResponse, result := queryIngredients("Spec", ingredientAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	IngredientFieldsAreExpected(t, result[0], createdIngredient)

	// cleanup
	deleteIngredient(createdIngredient.IngredientID.Hex(), ingredientAdminToken)
}

func TestNoFoundSearchIngredient(t *testing.T) {
	// set up, create an ingredient with a unique name
	specificNamedIngredient1 := defaultIngredient
	specificNamedIngredient1.Name = "Specific Ingredient Name"
	_, createdIngredient1 := createIngredient(specificNamedIngredient1, ingredientAdminToken)
	specificNamedIngredient2 := defaultIngredient
	specificNamedIngredient2.Name = "Another specific Ingredient name"
	_, createdIngredient2 := createIngredient(specificNamedIngredient2, ingredientAdminToken)

	// make a search that will not return anything
	searchResponse, result := queryIngredients("SomethingElse", ingredientAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, searchResponse.Code, "OK response is expected")
	assert.Equal(t, 0, len(result))

	// cleanup
	deleteIngredient(createdIngredient1.IngredientID.Hex(), ingredientAdminToken)
	deleteIngredient(createdIngredient2.IngredientID.Hex(), ingredientAdminToken)
}

func TestDeleteIngredient(t *testing.T) {
	// set up, create an ingredient
	_, createdIngredient := createIngredient(defaultIngredient, ingredientAdminToken)

	// delete an ingredient
	deleteResponse := deleteIngredient(createdIngredient.IngredientID.Hex(), ingredientAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 204, deleteResponse.Code, "OK response is expected")

	//try to get the Ingredient again to ensure it is deleted
	getResponse, _ := getIngredientByID(createdIngredient.IngredientID.Hex(), ingredientAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")

}

func TestCreateIngredientWithInvalidToken(t *testing.T) {
	// create an ingredient
	response, Ingredient := createIngredient(defaultIngredient, "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, response.Code, "Unauthorized response is expected")

	// cleanup
	deleteIngredient(Ingredient.IngredientID.Hex(), ingredientAdminToken)

}

func TestGetIngredientWithInvalidToken(t *testing.T) {
	// setUp, create an ingredient
	_, createdIngredient := createIngredient(defaultIngredient, ingredientAdminToken)

	// make a getIngredient for the ingredient we just created
	getResponse, _ := getIngredientByID(createdIngredient.IngredientID.Hex(), "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, getResponse.Code, "Unauthorized response is expected")

	// cleanup
	deleteIngredient(createdIngredient.IngredientID.Hex(), ingredientAdminToken)

}

func TestGetAllIngredientsWithInvalidToken(t *testing.T) {
	// setUp, create some Ingredients
	_, createdIngredient1 := createIngredient(defaultIngredient, ingredientAdminToken)
	_, createdIngredient2 := createIngredient(defaultIngredient, ingredientAdminToken)

	// make a getIngredient for the Ingredient we just created
	getResponse, _ := queryIngredients("Ingred", "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, getResponse.Code, "Unauthorized response is expected")

	// cleanup
	deleteIngredient(createdIngredient1.IngredientID.Hex(), ingredientAdminToken)
	deleteIngredient(createdIngredient2.IngredientID.Hex(), ingredientAdminToken)

}

func IngredientFieldsAreExpected(t *testing.T, Ingredient1 models.Ingredient, Ingredient2 models.Ingredient) {
	assert.Equal(t, Ingredient1.Name, Ingredient2.Name, "Inputted author value expected")
}
