package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/db"
	"strings"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

type UserMiddleware struct {
	auth AuthMiddleware
	Controller controller.UserControl
	repository db.UserDB
}

func NewUserMiddleware(auth AuthMiddleware, controller controller.UserController, db db.UserDB) UserMiddleware {
	return UserMiddleware{auth, controller, db}
}

//CreateUser creates a new user in the database
func (um UserMiddleware) CreateUser(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var requestedUser models.RequestedUser
	_ = json.NewDecoder(r.Body).Decode(&requestedUser)
	if requestedUser.Email == "" || requestedUser.Password == "" || requestedUser.UserName == "" {
		http.Error(w, "Required fields not included", http.StatusBadRequest)
	} else if !requestedUser.AgreedToTerms {
		http.Error(w, "User has not agreed to terms", http.StatusBadRequest)
	} else {
		payload, err := um.Controller.CreateUser(requestedUser, um.repository)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload.UserName)
		}
	}
}

//UpdateUserPassword updates a user's password in the database
func (um UserMiddleware) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	var updatedPassword models.UpdatedPassword
	_ = json.NewDecoder(r.Body).Decode(&updatedPassword)
	err := um.Controller.UpdateUserPassword(updatedPassword, um.repository)
	if err != nil {
		if err.Error() == "Username or password is not correct" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// DeleteUser controller DELETE request
func (um UserMiddleware) DeleteUser(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := um.auth.authenticateUser(w, r, true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		err := um.Controller.DeleteUser(params["userName"], um.repository)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

// GetUsers controller GET request
func (um UserMiddleware) GetUsers(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := um.auth.authenticateUser(w, r, true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		payload, err := um.Controller.GetUsers(um.repository)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(payload)
		}
	}
}

//GenerateUserToken refreshes a token
func (um UserMiddleware) GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var authData models.AuthData
	accessTokenObject := models.AccessToken{}
	_ = json.NewDecoder(r.Body).Decode(&authData)
	token, err := um.Controller.GenerateUserToken(authData, um.repository)
	if err != nil && err.Error() == "failed authentication, unknown user or password" {
		w.WriteHeader(http.StatusBadRequest)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		accessTokenObject.AccessToken = token
		json.NewEncoder(w).Encode(accessTokenObject)
	}
}

func (um UserMiddleware) EmailUser(w http.ResponseWriter, r *http.Request) {
	bearerToken := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var basket models.Basket
	_ = json.NewDecoder(r.Body).Decode(&basket)
	err := um.Controller.EmailUser(basket, bearerToken, um.repository)

	if err != nil {
		fmt.Println("Error Sending Email")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
