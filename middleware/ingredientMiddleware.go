package middleware

import (
	"encoding/json"
	"net/http"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

//CreateIngredient creates a new ingredient in the database
func CreateIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var requestedIngredient models.Ingredient
		_ = json.NewDecoder(r.Body).Decode(&requestedIngredient)
		payload, err := controller.CreateIngredient(requestedIngredient)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetIngredient by ID controller GET request
func GetIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		payload, err := controller.GetIngredient(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// QueryIngredient controller POST request
func QueryIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		prefixIngredient := r.URL.Query().Get("prefixIngredient")
		payload, err := controller.QueryIngredient(prefixIngredient)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// DeleteIngredient controller DELETE request
func DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	err := controller.DeleteIngredient(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
