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

var userRouter = router.Router()

var defaultNonAdminUser = models.RequestedUser{
	UserName: "testUser",
	Password: "password123",
	UserType: "nonAdmin",
}

var defaultNonAdminUser2 = models.RequestedUser{
	UserName: "testUser2",
	Password: "password1234",
	UserType: "nonAdmin",
}

var defaultNonAdminAuthData = models.AuthData{
	UserName: "testUser",
	Password: "password123",
}

// This user is always present in the testing database, do not delete
var defaultAdminUser = models.RequestedUser{
	UserName: "testAdminUser",
	Password: "adminPassword123",
	UserType: "admin",
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

func createUser(user models.RequestedUser, token string) *httptest.ResponseRecorder {
	jsonUser, _ := json.Marshal(user)
	createUserRequest, _ := http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonUser))
	createUserRequest.Header.Set("Authorization", "Bearer "+token)
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
	response := createUser(defaultNonAdminUser, adminToken)
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
	createUser(defaultNonAdminUser, adminToken)

	// delete the same user
	response := deleteUser(defaultNonAdminUser.UserName, adminToken)

	// assert the correct status code
	assert.Equal(t, 204, response.Code, "OK Response is expected")
}

func TestGetUsers(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser, adminToken)

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

func TestOnlyAdminCanCreateUsers(t *testing.T) {
	// generate a token with non admin user data
	createUser(defaultNonAdminUser, adminToken)
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultNonAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)

	//try using the nonAdmin's token to create another user
	createResponse := createUser(defaultNonAdminUser2, accessTokenObject.AccessToken)

	// assert the correct status code
	assert.Equal(t, 403, createResponse.Code, "Forbidden Response is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}

func TestOnlyAdminCanDeleteUsers(t *testing.T) {
	// generate a token with non admin user data, create a second user
	createUser(defaultNonAdminUser, adminToken)
	accessTokenObject := models.AccessToken{}
	response := generateUserToken(defaultNonAdminAuthData)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessTokenObject)
	createUser(defaultNonAdminUser2, adminToken)

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
	createUser(defaultNonAdminUser, adminToken)
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

func TestCannotUseTheSameUserName(t *testing.T) {
	// setup, create a user
	createUser(defaultNonAdminUser, adminToken)

	// try to create the same user
	createResponse := createUser(defaultNonAdminUser, adminToken)

	// assert the correct status code
	assert.Equal(t, 400, createResponse.Code, "Invalid Input Response is expected")

	// cleanup
	deleteUser(defaultNonAdminUser.UserName, adminToken)
}
