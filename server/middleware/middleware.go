package middleware

import (
	"encoding/json"
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
	payload := controller.GetAll()
	json.NewEncoder(w).Encode(payload)
}

// GetRecipe by ID controller GET request
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	payload := controller.Get(params["id"])
	json.NewEncoder(w).Encode(payload)
}

// CreateRecipe controller POST request
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var recipe models.Recipe
	_ = json.NewDecoder(r.Body).Decode(&recipe)
	_, code := controller.Create(recipe)
	json.NewEncoder(w).Encode(code)
}

// UpdateRecipe controller PUT request
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//Might need to change this to take the new recipe
	params := mux.Vars(r)
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	controller.Update(params["id"], recipe)
	json.NewEncoder(w).Encode(params["id"])
}

// DeleteRecipe controller DELETE request
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	controller.Delete(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}
