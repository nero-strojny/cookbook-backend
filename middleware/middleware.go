package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

func writeCommonHeaders(w http.ResponseWriter) {
	acceptedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", acceptedHeaders)
}

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

// PostPaginateRecipes controller POST request
func PostPaginatedRecipes(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var paginatedRequest models.PaginatedRecipeRequest
		_ = json.NewDecoder(r.Body).Decode(&paginatedRequest)
		payload, err := controller.PostPaginatedRecipes(paginatedRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetRecipe by ID controller GET request
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		payload, err := controller.GetRecipe(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetRandomRecipes gets a random number of recipes
func GetRandomRecipes(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		i, convertErr := strconv.Atoi(params["numberOfRecipes"])
		if convertErr != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			payload, err := controller.GetRandomRecipes(i)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(payload)
			}
		}
	}
}

//QueryRecipe searches for a recipe using the recipename
func QueryRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var recipe models.Recipe
		json.NewDecoder(r.Body).Decode(&recipe)
		payload, err := controller.QueryRecipe(recipe)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// CreateRecipe controller POST request
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var recipe models.Recipe
		_ = json.NewDecoder(r.Body).Decode(&recipe)
		payload, invalidFields, err := controller.CreateRecipe(recipe)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if len(invalidFields) > 0 {
				json.NewEncoder(w).Encode(invalidFields)
			}
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// UpdateRecipe controller PUT request
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	getPayload, getErr := controller.GetRecipe(params["id"])
	if getErr != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	userErr := authenticateSpecificUser(w, r, getPayload.UserName)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var recipe models.Recipe
		json.NewDecoder(r.Body).Decode(&recipe)
		payload, err := controller.UpdateRecipe(params["id"], recipe)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// DeleteRecipe controller DELETE request
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	getPayload, getErr := controller.GetRecipe(params["id"])
	if getErr != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	userErr := authenticateSpecificUser(w, r, getPayload.UserName)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		err := controller.DeleteRecipe(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
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

//CreateUser creates a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := authenticateUser(w, r, true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var requestedUser models.RequestedUser
		_ = json.NewDecoder(r.Body).Decode(&requestedUser)
		payload, err := controller.CreateUser(requestedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload.UserName)
		}
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

//SingleRecipeOptions eats options requests
func SingleRecipeOptions(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT, POST")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}

//CreateRecipeOptions handles preflight CORS for creating a recipe
func CreateRecipeOptions(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}

//SingleUserOptions eats options requests
func SingleUserOptions(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}

//GenerateUserTokenOptions handles preflight CORS
func GenerateUserTokenOptions(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}
