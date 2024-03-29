package middleware

import (
	"encoding/json"
	"net/http"
	"server/db"
	"strings"

	"server/controller"
	"server/models"

	"github.com/gorilla/mux"
)

type HouseholdMiddleware struct {
	auth         AuthMiddleware
	um           UserMiddleware
	controller   controller.HouseholdControl
	repository   db.HouseholdDB
	calendarRepo db.CalendarDB
}

func NewHouseholdMiddleware(auth AuthMiddleware, um UserMiddleware, controller controller.HouseholdControl, r db.HouseholdDB, c db.CalendarDB) HouseholdMiddleware {
	return HouseholdMiddleware{auth, um, controller, r, c}
}

//CreateHousehold creates a new household in the database
func (hm HouseholdMiddleware) CreateHousehold(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	userErr := hm.auth.AuthenticateUser(w, r, false)
	bearerToken := r.Header.Get("Authorization")
	currentUser, _ := hm.um.repository.GetUserByAccessToken(strings.ReplaceAll(bearerToken, "Bearer ", ""))
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var requestedHousehold models.Household
		_ = json.NewDecoder(r.Body).Decode(&requestedHousehold)
		payload, err := hm.controller.CreateHousehold(requestedHousehold, currentUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			_, updateErr := hm.controller.AddUserToHousehold(payload.HouseholdID.Hex(), currentUser.UserName, hm.um.repository)
			if updateErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

// GetHousehold by ID controller GET request
func (hm HouseholdMiddleware) GetHousehold(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := hm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		params := mux.Vars(r)
		payload, err := hm.controller.GetHousehold(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

func (hm HouseholdMiddleware) AddUserToHousehold(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	getPayload, getErr := hm.controller.GetHousehold(params["id"])
	if getErr != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	userErr := hm.auth.AuthenticateSpecificUser(w, r, getPayload.HeadOfHousehold)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var requestedHouseholdUpdate models.RequestedHouseholdUpdate
		json.NewDecoder(r.Body).Decode(&requestedHouseholdUpdate)
		payload, err := hm.controller.AddUserToHousehold(params["id"], requestedHouseholdUpdate.UserIdToAdd, hm.um.repository)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(payload)
		}
	}
}

func (hm HouseholdMiddleware) GetCalendar(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := hm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		startDate := r.URL.Query().Get("startDate")
		bearerToken := r.Header.Get("Authorization")
		currentUser, _ := hm.um.repository.GetUserByAccessToken(strings.ReplaceAll(bearerToken, "Bearer ", ""))
		payload, err := hm.controller.GetCalendar(currentUser.HouseholdId, startDate)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

func (hm HouseholdMiddleware) UpdateCalendar(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := hm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		var updatedCalendar models.Calendar
		_ = json.NewDecoder(r.Body).Decode(&updatedCalendar)
		bearerToken := r.Header.Get("Authorization")
		currentUser, _ := hm.um.repository.GetUserByAccessToken(strings.ReplaceAll(bearerToken, "Bearer ", ""))
		payload, err := hm.controller.UpdateCalendar(currentUser.HouseholdId, updatedCalendar)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(payload)
		}
	}
}

func (hm HouseholdMiddleware) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	userErr := hm.auth.AuthenticateUser(w, r, false)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
	} else {
		bearerToken := r.Header.Get("Authorization")
		user, _ := hm.um.repository.GetUserByAccessToken(strings.ReplaceAll(bearerToken, "Bearer ", ""))
		// Should this be a separate model, something like CalendarCreateRequest that only has the StartDate?
		var calendar models.Calendar
		json.NewDecoder(r.Body).Decode(&calendar)
		payload, err := hm.controller.CreateCalendar(calendar.StartDate, user.HouseholdId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload)
		}
	}
}
