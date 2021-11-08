package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

func authenticateUser(response http.ResponseWriter, request *http.Request, isAdmin bool) error {
	bearerToken := request.Header.Get("Authorization")
	userErr := controller.ValidateUser(strings.ReplaceAll(bearerToken, "Bearer ", ""), isAdmin)
	if userErr != nil {
		if userErr.Error() == "User does not have admin permissions" {
			response.WriteHeader(http.StatusForbidden)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}
	}
	return userErr
}

func authenticateSpecificUser(response http.ResponseWriter, request *http.Request, userName string) error {
	bearerToken := request.Header.Get("Authorization")
	userErr := controller.ValidateSpecificUser(strings.ReplaceAll(bearerToken, "Bearer ", ""), userName)
	if userErr != nil {
		response.WriteHeader(http.StatusUnauthorized)
	}
	return userErr
}

//CreateUser creates a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var requestedUser models.RequestedUser
	_ = json.NewDecoder(r.Body).Decode(&requestedUser)
	if requestedUser.Email == "" || requestedUser.Password == "" || requestedUser.UserName == "" {
		http.Error(w, "Required fields not included", http.StatusBadRequest)
	} else if !requestedUser.AgreedToTerms {
		http.Error(w, "User has not agreed to terms", http.StatusBadRequest)
	}
	payload, err := controller.CreateUser(requestedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(payload.UserName)
	}
}

//UpdateUserPassword updates a user's password in the database
func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	var updatedPassword models.UpdatedPassword
	_ = json.NewDecoder(r.Body).Decode(&updatedPassword)
	err := controller.UpdateUserPassword(updatedPassword)
	if err != nil {
		if err.Error() == "Username or password is not correct" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// DeleteUser controller DELETE request
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := authenticateUser(w, r, true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		err := controller.DeleteUser(params["userName"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

// GetUsers controller GET request
func GetUsers(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := authenticateUser(w, r, true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		payload, err := controller.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(payload)
		}
	}
}

//GenerateUserToken refreshes a token
func GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var authData models.AuthData
	accessTokenObject := models.AccessToken{}
	_ = json.NewDecoder(r.Body).Decode(&authData)
	token, err := controller.GenerateUserToken(authData)
	if err != nil && err.Error() == "failed authentication, unknown user or password" {
		w.WriteHeader(http.StatusBadRequest)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		accessTokenObject.AccessToken = token
		json.NewEncoder(w).Encode(accessTokenObject)
	}
}

func EmailUser(w http.ResponseWriter, r *http.Request) {
	bearerToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var basket models.Basket
	_ = json.NewDecoder(r.Body).Decode(&basket)
	err := controller.SendEmail(basket, bearerToken)

	if err != nil {
		fmt.Println("Error Sending Email")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
