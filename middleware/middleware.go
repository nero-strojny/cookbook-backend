package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

// GetAllRecipes controller GET request
func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload, err := controller.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}
}

// GetRecipe by ID controller GET request
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	payload, err := controller.Get(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}
}

//GetRecipeByName looks up a recipe by its exact name
func GetRecipeByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	payload, err := controller.Search(recipe.RecipeName)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}
}

// CreateRecipe controller POST request
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var recipe models.Recipe
	_ = json.NewDecoder(r.Body).Decode(&recipe)
	payload, invalidFields, err := controller.Create(recipe)
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

// UpdateRecipe controller PUT request
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	payload, err := controller.Update(params["id"], recipe)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}
}

// DeleteRecipe controller DELETE request
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	err := controller.Delete(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

//SingleRecipeOptions eats options requests
func SingleRecipeOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}

//CreateRecipeOptions handles preflight CORS for creating a recipe
func CreateRecipeOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}
