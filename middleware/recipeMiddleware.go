package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

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