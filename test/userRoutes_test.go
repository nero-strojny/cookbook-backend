package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"server/controller"
	"server/models"
	"server/router"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var userRouter = router.Router()

var defaultNonAdminUser = models.RequestedUser{
	UserName:      "testUser",
	Password:      "password123",
	UserType:      "nonAdmin",
	AgreedToTerms: true,
	Email:         "testnonadminemail@email.com",
}

var defaultNonAdminUser2 = models.RequestedUser{
	UserName:      "testUser2",
	Password:      "password1234",
	UserType:      "nonAdmin",
	AgreedToTerms: true,
	Email:         "testemail@email.com",
}

var defaultNonAdminAuthData = models.AuthData{
	UserName: "testUser",
	Password: "password123",
}

// This user is always present in the testing database, do not delete
var defaultAdminUser = models.RequestedUser{
	UserName:      "testAdminUser",
	Password:      "adminPassword123",
	UserType:      "admin",
	AgreedToTerms: true,
	Email:         "testemail@email.com",
}

var defaultAdminAuthData = models.AuthData{
	UserName: "testAdminUser",
	Password: "adminPassword123",
}

var adminToken string

func generateUserToken(authData models.AuthData) *httptest.ResponseRecorder {
	jsonAuth, _ := json.Marshal(authData)
	tokenRequest, _ := http.NewRequest("POST", "/api/userToken", bytes.NewBuffer(jsonAuth))
	tokenResponse := httptest.NewRecorder()
	userRouter.ServeHTTP(tokenResponse, tokenRequest)
	return tokenResponse
}

func createUser(user models.RequestedUser) *httptest.ResponseRecorder {
	jsonUser, _ := json.Marshal(user)
	createUserRequest, _ := http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonUser))
	createUserResponse := httptest.NewRecorder()
	userRouter.ServeHTTP(createUserResponse, createUserRequest)
	return createUserResponse
}

func deleteUser(userName string, token string) *httptest.ResponseRecorder {
	deleteUserRequest, _ := http.NewRequest("DELETE", "/api/user/"+userName, nil)
	deleteUserResponse := httptest.NewRecorder()
	deleteUserRequest.Header.Set("Authorization", "Bearer "+token)
	userRouter.ServeHTTP(deleteUserResponse, deleteUserRequest)
	return deleteUserResponse
}

func TestGenerateUserToken(t *testing.T) {
	// generate a token with admin user data
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)

	// assert the correct status code and body
	assert.Equal(t, 200, response.Code, "OK Response is expected")
	assert.NotNilf(t, accessTokenObject.AccessToken, "AccessToken should be set")
	adminToken = accessTokenObject.AccessToken
}

func TestCreateUser(t *testing.T) {
	// create a user
	var userName string
	response := createUser(defaultNonAdminUser)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &userName)

	// assert the correct status code and body
	assert.Equal(t, 201, response.Code, "OK Response is expected")
	assert.Equal(t, userName, defaultNonAdminUser.UserName, "Inputted UserName is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestDeleteUser(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser)

	// delete the same user
	response := deleteUser(defaultNonAdminUser.UserName, adminToken)

	// assert the correct status code
	assert.Equal(t, 204, response.Code, "OK Response is expected")
}

func TestGetUsers(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser)

	// get all the users
	users := []models.User{}
	getUserRequest, _ := http.NewRequest("GET", "/api/users", nil)
	getUserResponse := httptest.NewRecorder()
	getUserRequest.Header.Set("Authorization", "Bearer "+adminToken)
	userRouter.ServeHTTP(getUserResponse, getUserRequest)
	body, _ := ioutil.ReadAll(getUserResponse.Body)
	json.Unmarshal(body, &users)

	// assert the correct status code
	assert.Equal(t, 200, getUserResponse.Code, "OK Response is expected")
	assert.Equal(t, 2, len(users), "there should be 2 users in the database")
	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestUpdatedUserPassword(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser)

	var updatedPassword = models.UpdatedPassword{
		UserName:        "testUser",
		CurrentPassword: "password123",
		NewPassword:     "newPassword123",
	}

	jsonUserPassword, _ := json.Marshal(updatedPassword)
	updateUserRequest, _ := http.NewRequest("PUT", "/api/user", bytes.NewBuffer(jsonUserPassword))
	updateUserResponse := httptest.NewRecorder()
	userRouter.ServeHTTP(updateUserResponse, updateUserRequest)
	checkPastPasswordResponse := generateUserToken(defaultNonAdminAuthData)

	// assert the correct status code and body, check that the old password wouldn't work
	assert.Equal(t, 200, updateUserResponse.Code, "OK Response is expected")
	assert.Equal(t, 400, checkPastPasswordResponse.Code, "Invalid Input Response is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestNeedTokenToDeleteUsers(t *testing.T) {
	// generate a token with non admin user data, create a second user
	createUser(defaultNonAdminUser)

	// delete the same user
	deleteResponse := deleteUser(defaultNonAdminUser.UserName, "")

	// assert the correct status code
	assert.Equal(t, 401, deleteResponse.Code, "Unauthorized Response is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestNeedTokenToGetUsers(t *testing.T) {
	// get all the users
	getUserRequest, _ := http.NewRequest("GET", "/api/users", nil)
	getUserResponse := httptest.NewRecorder()
	getUserRequest.Header.Set("Authorization", "Bearer ")
	userRouter.ServeHTTP(getUserResponse, getUserRequest)

	// assert the correct status code
	assert.Equal(t, 401, getUserResponse.Code, "Unauthorized Response is expected")
	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestNeedAllFieldsToCreateUser(t *testing.T) {
	var incompleteUser = models.RequestedUser{
		UserName:      "testAdminUser",
		Password:      "adminPassword123",
		UserType:      "admin",
		AgreedToTerms: true,
	}

	//try using the nonAdmin's token to create another user
	createResponse := createUser(incompleteUser)

	// assert the correct status code
	assert.Equal(t, 400, createResponse.Code, "Bad Request Response is expected")
}

func TestNeedToAgreeToCreateUser(t *testing.T) {
	var badUser = models.RequestedUser{
		UserName:      "someUser",
		Password:      "adminPassword123",
		UserType:      "admin",
		AgreedToTerms: false,
		Email:         "someEmail.com",
	}

	//try using the nonAdmin's token to create another user
	createResponse := createUser(badUser)

	// assert the correct status code
	assert.Equal(t, 400, createResponse.Code, "Bad Request Response is expected")
}

func TestOnlyAdminCanDeleteUsers(t *testing.T) {
	// generate a token with non admin user data, create a second user
	createUser(defaultNonAdminUser)
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultNonAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)
	createUser(defaultNonAdminUser2)

	// delete the same user
	deleteResponse := deleteUser(defaultNonAdminUser.UserName, accessTokenObject.AccessToken)

	// assert the correct status code
	assert.Equal(t, 403, deleteResponse.Code, "Forbidden Response is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
	deleteUser(defaultNonAdminUser2.UserName, adminToken)
}

func TestOnlyAdminCanGetUsers(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser)
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultNonAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)

	// get all the users
	getUserRequest, _ := http.NewRequest("GET", "/api/users", nil)
	getUserResponse := httptest.NewRecorder()
	getUserRequest.Header.Set("Authorization", "Bearer "+accessTokenObject.AccessToken)
	userRouter.ServeHTTP(getUserResponse, getUserRequest)

	// assert the correct status code
	assert.Equal(t, 403, getUserResponse.Code, "Forbidden Response is expected")
	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestErrorOnWrongPassword(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser)

	var defaultNonAdminAuthData = models.AuthData{
		UserName: "testUser",
		Password: "wrongPassword",
	}

	response := generateUserToken(defaultNonAdminAuthData)

	// assert the correct status code
	assert.Equal(t, 400, response.Code, "Invalid Input Response is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestCannotUseTheSameUserName(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser)

	// try to create the same user
	createResponse := createUser(defaultNonAdminUser)

	// assert the correct status code
	assert.Equal(t, 400, createResponse.Code, "Invalid Input Response is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestExpiredPassword(t *testing.T) {
	// setup, change the user's password
	expiryTime := time.Now().AddDate(0, 0, -1)
	result := models.User{}
	getFilter := bson.M{"username": "testAdminUser"}
	controller.UserCollection.FindOne(context.Background(), getFilter).Decode(&result)
	result.ExpiryDate = expiryTime.Format("2006.01.02 15:04:05")
	controller.UserCollection.ReplaceOne(context.Background(), getFilter, result)

	response := deleteUser(defaultNonAdminUser.UserName, adminToken)

	// assert the correct status code
	assert.Equal(t, 401, response.Code, "Unauthorized Response is expected")
}
