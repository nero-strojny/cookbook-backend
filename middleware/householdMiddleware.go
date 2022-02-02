package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

//CreateHousehold creates a new household in the database
func CreateHousehold(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := authenticateUser(w, r, false)
	bearerToken := r.Header.Get("Authorization")
	currentUser, _ := controller.RetrieveUserFromToken(strings.ReplaceAll(bearerToken, "Bearer ", ""))
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var requestedHousehold models.Household
		_ = json.NewDecoder(r.Body).Decode(&requestedHousehold)
		payload, err := controller.CreateHousehold(requestedHousehold, currentUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			_, updateErr := controller.AddUserToHousehold(payload.HouseholdID.Hex(), currentUser.UserID.Hex())
			if updateErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetHousehold by ID controller GET request
func GetHousehold(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := authenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		payload, err := controller.GetHousehold(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

func AddUserToHousehold(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	getPayload, getErr := controller.GetHousehold(params["id"])
	if getErr != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	userErr := authenticateSpecificUser(w, r, getPayload.HeadOfHousehold)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var requestedHouseholdUpdate models.RequestedHouseholdUpdate
		json.NewDecoder(r.Body).Decode(&requestedHouseholdUpdate)
		payload, err := controller.AddUserToHousehold(params["id"], requestedHouseholdUpdate.UserIdToAdd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(payload)
		}
	}
}
