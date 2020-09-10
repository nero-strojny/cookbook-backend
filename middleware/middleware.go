package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

func writeCommonHeaders(w http.ResponseWriter) http.ResponseWriter {
	acceptedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", acceptedHeaders)
	return w
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
		if userErr.Error() == "User does not have admin permissions" {
			response.WriteHeader(http.StatusForbidden)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}
	}
	return userErr
}

// GetAllRecipes controller GET request
func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "GET")
	payload, err := controller.GetAllRecipes()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(payload)
	}
}

// GetRecipe by ID controller GET request
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	payload, err := controller.GetRecipe(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(payload)
	}
}

//GetRecipeByName looks up a recipe by its exact name
func GetRecipeByName(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "GET")
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	payload, err := controller.SearchRecipe(recipe.RecipeName)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
	} else {
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(payload)
	}
}

// CreateRecipe controller POST request
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	var recipe models.Recipe
	_ = json.NewDecoder(r.Body).Decode(&recipe)
	payload, invalidFields, err := controller.CreateRecipe(recipe)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		if len(invalidFields) > 0 {
			json.NewEncoder(response).Encode(invalidFields)
		}
	} else {
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(payload)
	}
}

// UpdateRecipe controller PUT request
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	payload, err := controller.UpdateRecipe(params["id"], recipe)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
	} else {
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(payload)
	}
}

// DeleteRecipe controller DELETE request
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	err := controller.DeleteRecipe(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

//SingleRecipeOptions eats options requests
func SingleRecipeOptions(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode("")
}

//CreateRecipeOptions handles preflight CORS for creating a recipe
func CreateRecipeOptions(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode("")
}

// GetCalorieLog by ID controller GET request
func GetCalorieLog(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	payload, err := controller.GetCalorieLog(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(payload)
	}
}

// CreateCalorieLog controller POST request
func CreateCalorieLog(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	var calorieLog models.CalorieLog
	_ = json.NewDecoder(r.Body).Decode(&calorieLog)
	payload, err := controller.CreateCalorieLog(calorieLog)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
	} else {
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(payload)
	}
}

// UpdateCalorieLog controller PUT request
func UpdateCalorieLog(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	var calorieLog models.CalorieLog
	json.NewDecoder(r.Body).Decode(&calorieLog)
	payload, err := controller.UpdateCalorieLog(params["id"], calorieLog)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
	} else {
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(payload)
	}
}

// DeleteCalorieLog controller DELETE request
func DeleteCalorieLog(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	err := controller.DeleteCalorieLog(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

//SingleCalorieLogOptions eats options requests
func SingleCalorieLogOptions(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode("")
}

//CreateCalorieLogOptions handles preflight CORS for creating a calorie log
func CreateCalorieLogOptions(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode("")
}

//GenerateUserToken refreshes a token
func GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	var authData models.AuthData
	accessTokenObject := models.AccessToken{}
	_ = json.NewDecoder(r.Body).Decode(&authData)
	token, err := controller.GenerateUserToken(authData)
	if err != nil && err.Error() == "failed authentication, unknown user or password" {
		response.WriteHeader(http.StatusBadRequest)
	} else if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		accessTokenObject.AccessToken = token
		json.NewEncoder(response).Encode(accessTokenObject)
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

//SingleUserOptions eats options requests
func SingleUserOptions(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "PUT")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode("")
}

//GenerateUserTokenOptions handles preflight CORS for creating a calorie log
func GenerateUserTokenOptions(w http.ResponseWriter, r *http.Request) {
	response := writeCommonHeaders(w)
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode("")
}
