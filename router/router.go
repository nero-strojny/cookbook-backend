package router

import (
	"server/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/recipes", middleware.GetAllRecipes).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/recipe/{id}", middleware.GetRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/{id}", middleware.DeleteRecipe).Methods("DELETE")
	router.HandleFunc("/api/recipe/{id}", middleware.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/api/recipe/{id}", middleware.SingleRecipeOptions).Methods("OPTIONS")

	router.HandleFunc("/api/recipe/search", middleware.GetRecipeByName).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/recipe", middleware.CreateRecipe).Methods("POST")
	router.HandleFunc("/api/recipe", middleware.CreateRecipeOptions).Methods("OPTIONS")

	router.HandleFunc("/api/calorieLog/{id}", middleware.GetCalorieLog).Methods("GET")
	router.HandleFunc("/api/calorieLog/{id}", middleware.DeleteCalorieLog).Methods("DELETE")
	router.HandleFunc("/api/calorieLog/{id}", middleware.UpdateCalorieLog).Methods("PUT")
	router.HandleFunc("/api/calorieLog/{id}", middleware.SingleCalorieLogOptions).Methods("OPTIONS")

	router.HandleFunc("/api/calorieLog", middleware.CreateCalorieLog).Methods("POST")
	router.HandleFunc("/api/calorieLog", middleware.CreateCalorieLogOptions).Methods("OPTIONS")

	router.HandleFunc("/api/user", middleware.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/{userId}", middleware.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/users", middleware.GetUsers).Methods("GET")

	router.HandleFunc("/api/userToken", middleware.GenerateUserToken).Methods("POST")
	router.HandleFunc("/api/userToken", middleware.GenerateUserTokenOptions).Methods("OPTIONS")

	return router
}
