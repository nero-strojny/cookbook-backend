package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

type RecipeMiddleware struct {
	auth       AuthMiddleware
	controller controller.RecipeControl
}

func NewRecipeMiddleware(auth AuthMiddleware, controller controller.RecipeController) RecipeMiddleware {
	return RecipeMiddleware{auth, controller}
}

// PostPaginateRecipes controller POST request
func (rm RecipeMiddleware) PostPaginatedRecipes(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := rm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var paginatedRequest models.PaginatedRecipeRequest
		_ = json.NewDecoder(r.Body).Decode(&paginatedRequest)
		payload, err := rm.controller.PostPaginatedRecipes(paginatedRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetRecipe by ID controller GET request
func (rm RecipeMiddleware) GetRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := rm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		payload, err := rm.controller.GetRecipe(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetRandomRecipes gets a random number of recipes
func (rm RecipeMiddleware) GetRandomRecipes(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := rm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		i, convertErr := strconv.Atoi(params["numberOfRecipes"])
		if convertErr != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			payload, err := rm.controller.GetRandomRecipes(i)
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
func (rm RecipeMiddleware) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := rm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var recipe models.Recipe
		_ = json.NewDecoder(r.Body).Decode(&recipe)
		payload, invalidFields, err := rm.controller.CreateRecipe(recipe)
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
func (rm RecipeMiddleware) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	getPayload, getErr := rm.controller.GetRecipe(params["id"])
	if getErr != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	userErr := rm.auth.AuthenticateSpecificUser(w, r, getPayload.UserName)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var recipe models.Recipe
		json.NewDecoder(r.Body).Decode(&recipe)
		payload, err := rm.controller.UpdateRecipe(params["id"], recipe)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// DeleteRecipe controller DELETE request
func (rm RecipeMiddleware) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	getPayload, getErr := rm.controller.GetRecipe(params["id"])
	if getErr != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	userErr := rm.auth.AuthenticateSpecificUser(w, r, getPayload.UserName)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		err := rm.controller.DeleteRecipe(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
