package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"server/controller"
	"server/models"
	"server/util"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	dbConnectionString := util.OpenFile()
	clientOptions := options.Client().ApplyURI(dbConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	controller.SetRecipeClient(client)
	controller.SetCalorieClient(client)
	controller.SetUserClient(client)
}

// GetAllRecipes controller GET request
func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	payload, err := controller.GetAllRecipes()
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	payload, err := controller.GetRecipe(params["id"])
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	payload, err := controller.SearchRecipe(recipe.RecipeName)
	if err != nil {
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

// UpdateRecipe controller PUT request
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	payload, err := controller.UpdateRecipe(params["id"], recipe)
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
	err := controller.DeleteRecipe(params["id"])
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

// GetAllCalorieLogs controller GET request
func GetAllCalorieLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	payload, err := controller.GetAllCalorieLogs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}
}

// GetCalorieLog by ID controller GET request
func GetCalorieLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	payload, err := controller.GetCalorieLog(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}
}

// CreateCalorieLog controller POST request
func CreateCalorieLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var calorieLog models.CalorieLog
	_ = json.NewDecoder(r.Body).Decode(&calorieLog)
	payload, err := controller.CreateCalorieLog(calorieLog)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(payload)
	}
}

// UpdateCalorieLog controller PUT request
func UpdateCalorieLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	var calorieLog models.CalorieLog
	json.NewDecoder(r.Body).Decode(&calorieLog)
	payload, err := controller.UpdateCalorieLog(params["id"], calorieLog)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}
}

// DeleteCalorieLog controller DELETE request
func DeleteCalorieLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	err := controller.DeleteCalorieLog(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

//SingleCalorieLogOptions eats options requests
func SingleCalorieLogOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}

//CreateCalorieLogOptions handles preflight CORS for creating a calorie log
func CreateCalorieLogOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}

//GenerateUserToken refreshes a token
func GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var authData models.AuthData
	_ = json.NewDecoder(r.Body).Decode(&authData)
	token, err := controller.GenerateUserToken(authData)
	if err != nil && token == "failed authentication, unknown user or password" {
		json.NewEncoder(w).Encode(token)
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(token)
	}
}

//CreateUser creates a new user in the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	bearerToken := r.Header.Get("Authorization")
	userErr := controller.ValidateUser(strings.ReplaceAll(bearerToken, "Bearer ", ""), true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
		if userErr.Error() == "User does not have admin permissions" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
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

// DeleteUser controller DELETE request
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	bearerToken := r.Header.Get("Authorization")
	userErr := controller.ValidateUser(strings.ReplaceAll(bearerToken, "Bearer ", ""), true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
		if userErr.Error() == "User does not have admin permissions" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		params := mux.Vars(r)
		err := controller.DeleteUser(params["userId"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

// GetUsers controller GET request
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	bearerToken := r.Header.Get("Authorization")
	userErr := controller.ValidateUser(strings.ReplaceAll(bearerToken, "Bearer ", ""), true)
	if userErr != nil {
		json.NewEncoder(w).Encode(userErr.Error())
		if userErr.Error() == "User does not have admin permissions" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		payload, err := controller.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
		}
	}
}

//GenerateUserTokenOptions handles preflight CORS for creating a calorie log
func GenerateUserTokenOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
