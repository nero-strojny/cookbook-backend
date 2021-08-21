package router

import (
	"server/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/recipes", middleware.PostPaginatedRecipes).Methods("POST")
	router.HandleFunc("/api/recipes", middleware.SingleRecipeOptions).Methods("OPTIONS")
	router.HandleFunc("/api/recipe/{id}", middleware.GetRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/{id}", middleware.DeleteRecipe).Methods("DELETE")
	router.HandleFunc("/api/recipe/{id}", middleware.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/api/recipe/{id}", middleware.SingleRecipeOptions).Methods("OPTIONS")

	router.HandleFunc("/api/recipe", middleware.CreateRecipe).Methods("POST")
	router.HandleFunc("/api/recipe", middleware.CreateRecipeOptions).Methods("OPTIONS")

	router.HandleFunc("/api/randomRecipe/{numberOfRecipes}", middleware.GetRandomRecipes).Methods("GET")
	router.HandleFunc("/api/randomRecipe/{numberOfRecipes}", middleware.SingleRecipeOptions).Methods("OPTIONS")

	router.HandleFunc("/api/user", middleware.CreateUser).Methods("POST")
	router.HandleFunc("/api/user", middleware.UpdateUserPassword).Methods("PUT")
	router.HandleFunc("/api/user/{userName}", middleware.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/users", middleware.GetUsers).Methods("GET")
	router.HandleFunc("/api/user", middleware.SingleUserOptions).Methods("OPTIONS")

	router.HandleFunc("/api/userToken", middleware.GenerateUserToken).Methods("POST")
	router.HandleFunc("/api/userToken", middleware.GenerateUserTokenOptions).Methods("OPTIONS")

	return router
}
