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

var calorieLogRouter = router.Router()

var defaultCalorieLog = models.CalorieLog{
	Calories:    500,
	Description: "test",
	UserName:    "testUser",
}

func createCalorieLog(calorieLog models.CalorieLog) *httptest.ResponseRecorder {
	jsonCalorieLog, _ := json.Marshal(calorieLog)
	request, _ := http.NewRequest("POST", "/api/calorieLog", bytes.NewBuffer(jsonCalorieLog))
	response := httptest.NewRecorder()
	calorieLogRouter.ServeHTTP(response, request)
	return response
}

func deleteCalorieLog(calorieLogID string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("DELETE", "/api/calorieLog/"+calorieLogID, nil)
	response := httptest.NewRecorder()
	calorieLogRouter.ServeHTTP(response, request)
	return response
}

func TestCreateCalorieLog(t *testing.T) {
	// create a calorieLog
	calorieLog := models.CalorieLog{}
	response := createCalorieLog(defaultCalorieLog)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &calorieLog)

	// assert the correct status code and body
	assert.Equal(t, 201, response.Code, "OK response is expected")
	calorieLogFieldsAreExpected(t, calorieLog, defaultCalorieLog)

	// cleanup
	deleteCalorieLog(calorieLog.CalorieLogID.Hex())
}

func TestGetCalorieLog(t *testing.T) {
	// setUp, create a calorieLog
	createdCalorieLog := models.CalorieLog{}
	createResponse := createCalorieLog(defaultCalorieLog)
	createBody, _ := ioutil.ReadAll(createResponse.Body)
	json.Unmarshal(createBody, &createdCalorieLog)

	// make a getCalorieLog for the calorieLog we just created
	getCalorieLog := models.CalorieLog{}
	getRequest, _ := http.NewRequest("GET", "/api/calorieLog/"+createdCalorieLog.CalorieLogID.Hex(), nil)
	getResponse := httptest.NewRecorder()
	calorieLogRouter.ServeHTTP(getResponse, getRequest)
	getBody, _ := ioutil.ReadAll(getResponse.Body)
	json.Unmarshal(getBody, &getCalorieLog)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	calorieLogFieldsAreExpected(t, getCalorieLog, defaultCalorieLog)

	// cleanup
	deleteCalorieLog(createdCalorieLog.CalorieLogID.Hex())
}

func TestGetCalorieLogWithUnknownCalorieLog(t *testing.T) {
	// try to get a calorieLog with an unknown id
	getCalorieLog := models.CalorieLog{}
	getRequest, _ := http.NewRequest("GET", "/api/calorieLog/unknownCalorieLogId", nil)
	getResponse := httptest.NewRecorder()
	calorieLogRouter.ServeHTTP(getResponse, getRequest)
	getBody, _ := ioutil.ReadAll(getResponse.Body)
	json.Unmarshal(getBody, &getCalorieLog)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")
}

func TestUpdateCalorieLog(t *testing.T) {
	// setUp, create a calorieLog
	createdCalorieLog := models.CalorieLog{}
	createResponse := createCalorieLog(defaultCalorieLog)
	createBody, _ := ioutil.ReadAll(createResponse.Body)
	json.Unmarshal(createBody, &createdCalorieLog)

	updatedCalorieLog := createdCalorieLog
	updatedCalorieLog.Calories = 200
	jsonCalorieLog, _ := json.Marshal(updatedCalorieLog)
	updatedRequest, _ := http.NewRequest("PUT", "/api/calorieLog/"+createdCalorieLog.CalorieLogID.Hex(), bytes.NewBuffer(jsonCalorieLog))
	updatedResponse := httptest.NewRecorder()
	calorieLogRouter.ServeHTTP(updatedResponse, updatedRequest)
	updatedBody, _ := ioutil.ReadAll(updatedResponse.Body)
	json.Unmarshal(updatedBody, &updatedCalorieLog)

	// make a getCalorieLog for the calorieLog we just updated
	getCalorieLog := models.CalorieLog{}
	getRequest, _ := http.NewRequest("GET", "/api/calorieLog/"+createdCalorieLog.CalorieLogID.Hex(), nil)
	getResponse := httptest.NewRecorder()
	calorieLogRouter.ServeHTTP(getResponse, getRequest)
	getBody, _ := ioutil.ReadAll(getResponse.Body)
	json.Unmarshal(getBody, &getCalorieLog)

	// assert the correct status code and body
	assert.Equal(t, 200, getResponse.Code, "OK response is expected")
	assert.Equal(t, 200, getCalorieLog.Calories)

	// cleanup
	deleteCalorieLog(createdCalorieLog.CalorieLogID.Hex())
}

func TestDeleteCalorieLog(t *testing.T) {
	// set up, create a calorieLog
	calorieLog := models.CalorieLog{}
	createResponse := createCalorieLog(defaultCalorieLog)
	body, _ := ioutil.ReadAll(createResponse.Body)
	json.Unmarshal(body, &calorieLog)

	// delete a calorieLog
	deleteResponse := deleteCalorieLog(calorieLog.CalorieLogID.Hex())

	// assert the correct status code and body
	assert.Equal(t, 204, deleteResponse.Code, "OK response is expected")

	//try to get the calorieLog again to ensure it is deleted
	getRequest, _ := http.NewRequest("GET", "/api/calorieLog/unknownCalorieLogId", nil)
	getResponse := httptest.NewRecorder()
	calorieLogRouter.ServeHTTP(getResponse, getRequest)

	// assert the correct status code and body
	assert.Equal(t, 404, getResponse.Code, "Not Found response is expected")

}

func TestDeleteCalorieLogWithUnknownID(t *testing.T) {
	// delete a calorieLog
	deleteResponse := deleteCalorieLog("unknownCalorieLogID")

	// assert the correct status code and body
	assert.Equal(t, 404, deleteResponse.Code, "Not Found response is expected")
}

func calorieLogFieldsAreExpected(t *testing.T, calorieLog1 models.CalorieLog, calorieLog2 models.CalorieLog) {
	assert.NotNilf(t, calorieLog1.EnteredDate, "EnteredDate should be set")
	assert.Equal(t, calorieLog1.UserName, calorieLog2.UserName, "Inputted UserName value expected")
	assert.Equal(t, calorieLog1.Calories, calorieLog2.Calories, "Inputted Calories is expected")
	assert.Equal(t, calorieLog1.Description, calorieLog2.Description, "Inputted Description is expected")
}
