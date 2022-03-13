package middleware

import (
	"encoding/json"
	"net/http"
	"server/db"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

type IngredientHandler interface {
	CreateIngredient(w http.ResponseWriter, r *http.Request)
	GetIngredient(w http.ResponseWriter, r *http.Request)
	QueryIngredient(w http.ResponseWriter, r *http.Request)
	DeleteIngredient(w http.ResponseWriter, r *http.Request)
}

type IngredientMiddleware struct {
	auth       AuthHandler
	controller controller.IngredientControl
	repository db.IngredientDB
}

func NewIngredientMiddleware(auth AuthHandler, controller controller.IngredientControl, r db.IngredientDB) IngredientMiddleware {
	return IngredientMiddleware{auth, controller, r}
}

//CreateIngredient creates a new ingredient in the database
func (im IngredientMiddleware) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := im.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var requestedIngredient models.Ingredient
		_ = json.NewDecoder(r.Body).Decode(&requestedIngredient)
		payload, err := im.controller.CreateIngredient(requestedIngredient, im.repository)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetIngredient by ID controller GET request
func (im IngredientMiddleware) GetIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := im.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		payload, err := im.controller.GetIngredient(params["id"], im.repository)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// QueryIngredient controller POST request
func (im IngredientMiddleware) QueryIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := im.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		prefixIngredient := r.URL.Query().Get("prefixIngredient")
		payload, err := im.controller.QueryIngredient(prefixIngredient, im.repository)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// DeleteIngredient controller DELETE request
func (im IngredientMiddleware) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	err := im.controller.DeleteIngredient(params["id"], im.repository)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
