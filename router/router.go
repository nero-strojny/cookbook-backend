package router

import (
	"server/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/recipes", middleware.PostPaginatedRecipes).Methods("POST")
	router.HandleFunc("/api/recipes", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/recipe/{id}", middleware.GetRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/{id}", middleware.DeleteRecipe).Methods("DELETE")
	router.HandleFunc("/api/recipe/{id}", middleware.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/api/recipe/{id}", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/recipe", middleware.CreateRecipe).Methods("POST")
	router.HandleFunc("/api/recipe", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/randomRecipe/{numberOfRecipes}", middleware.GetRandomRecipes).Methods("GET")
	router.HandleFunc("/api/randomRecipe/{numberOfRecipes}", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/ingredients", middleware.QueryIngredient).Queries(
		"prefixIngredient", "{prefixIngredient}",
	).Methods("POST")
	router.HandleFunc("/api/ingredients", middleware.Options).Queries(
		"prefixIngredient", "{prefixIngredient}",
	).Methods("OPTIONS")
	router.HandleFunc("/api/ingredients", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/ingredient/{id}", middleware.GetIngredient).Methods("GET")
	router.HandleFunc("/api/ingredient/{id}", middleware.DeleteIngredient).Methods("DELETE")
	router.HandleFunc("/api/ingredient/{id}", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/ingredient", middleware.CreateIngredient).Methods("POST")
	router.HandleFunc("/api/ingredient", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/user", middleware.CreateUser).Methods("POST")
	router.HandleFunc("/api/user", middleware.UpdateUserPassword).Methods("PUT")
	router.HandleFunc("/api/user", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/users", middleware.GetUsers).Methods("GET")

	router.HandleFunc("/api/user/{userName}", middleware.DeleteUser).Methods("DELETE")

	router.HandleFunc("/api/userToken", middleware.GenerateUserToken).Methods("POST")
	router.HandleFunc("/api/userToken", middleware.Options).Methods("OPTIONS")

	return router
}
