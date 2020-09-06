package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/models"
	"server/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

var calorieRouter = router.Router()

var defaultCalorieLog = models.CalorieLog{
	Calories:    500,
	Description: "test",
	UserID:      "testUser",
}

func createCalorieLog(calorieLog models.CalorieLog) *httptest.ResponseRecorder {
	jsonCalorieLog, _ := json.Marshal(calorieLog)
	request, _ := http.NewRequest("POST", "/api/calorieLog", bytes.NewBuffer(jsonCalorieLog))
	response := httptest.NewRecorder()
	calorieRouter.ServeHTTP(response, request)
	return response
}

func deleteCalorieLog(calorieLogID string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("DELETE", "/api/calorieLog/"+calorieLogID, nil)
	response := httptest.NewRecorder()
	calorieRouter.ServeHTTP(response, request)
	return response
}

func calorieLogFieldsAreExpected(t *testing.T, calorieLog1 models.CalorieLog, calorieLog2 models.CalorieLog) {
	assert.NotNilf(t, calorieLog1.EnteredDate, "EnteredDate should be set")
	assert.Equal(t, calorieLog1.UserID, calorieLog2.UserID, "Inputted UserID value expected")
	assert.Equal(t, calorieLog1.Calories, calorieLog2.Calories, "Inputted Calories is expected")
	assert.Equal(t, calorieLog1.Description, calorieLog2.Description, "Inputted Description is expected")

}
