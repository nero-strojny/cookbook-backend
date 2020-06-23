package router

import (
	"server/middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/recipes", middleware.GetAllRecipes).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/recipe/{id}", middleware.GetRecipe).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/recipe", middleware.CreateRecipe).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/recipe/{id}", middleware.DeleteRecipe).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/recipe/{id}", middleware.UpdateRecipe).Methods("PUT", "OPTIONS")
	return router
}
