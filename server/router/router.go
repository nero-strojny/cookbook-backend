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

	router.HandleFunc("/api/recipe", middleware.CreateRecipe).Methods("POST")
	router.HandleFunc("/api/recipe", middleware.CreateRecipeOptions).Methods("OPTIONS")

	return router
}
