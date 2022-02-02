package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"server/controller"
	"server/models"
	"server/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

var householdRouter = router.Router()
var householdAdminToken string

var defaultHousehold = models.Household{
	HouseholdName: "Household",
}

func createHousehold(inputHousehold models.Household, accessToken string) (*httptest.ResponseRecorder, models.Household) {
	outputHousehold := models.Household{}
	jsonHousehold, _ := json.Marshal(inputHousehold)
	request, _ := http.NewRequest("POST", "/api/household", bytes.NewBuffer(jsonHousehold))
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	householdRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &outputHousehold)
	return response, outputHousehold
}

func getHouseholdByID(HouseholdID string, accessToken string) (*httptest.ResponseRecorder, models.Household) {
	Household := models.Household{}
	request, _ := http.NewRequest("GET", "/api/household/"+HouseholdID, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	response := httptest.NewRecorder()
	householdRouter.ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &Household)
	return response, Household
}

func TestHouseholdSetUp(t *testing.T) {
	// generate a token with admin user data
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)
	householdAdminToken = accessTokenObject.AccessToken
}

func TestCreateHousehold(t *testing.T) {
	// create an household
	response, Household := createHousehold(defaultHousehold, householdAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 201, response.Code, "OK response is expected")
	HouseholdFieldsAreExpected(t, Household, defaultHousehold)

	//cleanup
	controller.DeleteHousehold(Household.HouseholdID.Hex())
}

func TestGetHousehold(t *testing.T) {
	// setUp, create an household
	_, createdHousehold := createHousehold(defaultHousehold, householdAdminToken)

	// make a getHousehold for the Household we just created
	getResponse, getHousehold := getHouseholdByID(createdHousehold.HouseholdID.Hex(), householdAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	HouseholdFieldsAreExpected(t, getHousehold, defaultHousehold)
	//cleanup
	controller.DeleteHousehold(createdHousehold.HouseholdID.Hex())
}

func TestGetHouseholdWithUnknownHousehold(t *testing.T) {
	// try to get an household with an unknown id
	getResponse, _ := getHouseholdByID("unknownHouseholdId", householdAdminToken)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")
}

func TestCreateHouseholdWithInvalidToken(t *testing.T) {
	// create an household
	response, _ := createHousehold(defaultHousehold, "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, response.Code, "Unauthorized response is expected")
}

func TestGetHouseholdWithInvalidToken(t *testing.T) {
	// setUp, create an household
	_, createdHousehold := createHousehold(defaultHousehold, householdAdminToken)

	// make a getHousehold for the household we just created
	getResponse, _ := getHouseholdByID(createdHousehold.HouseholdID.Hex(), "invalidToken")

	// assert the correct status code and body
	assert.Equal(t, 401, getResponse.Code, "Unauthorized response is expected")
	//cleanup
	controller.DeleteHousehold(createdHousehold.HouseholdID.Hex())
}

func HouseholdFieldsAreExpected(t *testing.T, Household1 models.Household, Household2 models.Household) {
	assert.Equal(t, Household1.HouseholdName, Household2.HouseholdName, "Inputted name value expected")
}
