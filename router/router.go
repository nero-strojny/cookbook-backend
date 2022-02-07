package router

import (
	"server/middleware"

	"github.com/gorilla/mux"
)

type TastyBoiRouter struct {
	um middleware.UserMiddleware
	rm middleware.RecipeMiddleware
	im middleware.IngredientMiddleware
	sm middleware.ServerMiddleware
	hm middleware.HouseholdMiddleware
}

func NewTastyBoiRouter(um middleware.UserMiddleware,
	rm middleware.RecipeMiddleware,
	im middleware.IngredientMiddleware,
	sm middleware.ServerMiddleware,
	hm middleware.HouseholdMiddleware) TastyBoiRouter {
	return TastyBoiRouter{um, rm, im, sm, hm}
}

// Route is exported and used in main.go
func (r TastyBoiRouter) Route() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/recipes", r.rm.PostPaginatedRecipes).Methods("POST")
	router.HandleFunc("/api/recipes", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/recipe/{id}", r.rm.GetRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/{id}", r.rm.DeleteRecipe).Methods("DELETE")
	router.HandleFunc("/api/recipe/{id}", r.rm.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/api/recipe/{id}", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/recipe", r.rm.CreateRecipe).Methods("POST")
	router.HandleFunc("/api/recipe", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/randomRecipe/{numberOfRecipes}", r.rm.GetRandomRecipes).Methods("GET")
	router.HandleFunc("/api/randomRecipe/{numberOfRecipes}", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/ingredients", r.im.QueryIngredient).Queries(
		"prefixIngredient", "{prefixIngredient}",
	).Methods("GET")
	router.HandleFunc("/api/ingredients", middleware.Options).Queries(
		"prefixIngredient", "{prefixIngredient}",
	).Methods("OPTIONS")
	router.HandleFunc("/api/ingredients", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/ingredient/{id}", r.im.GetIngredient).Methods("GET")
	router.HandleFunc("/api/ingredient/{id}", r.im.DeleteIngredient).Methods("DELETE")
	router.HandleFunc("/api/ingredient/{id}", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/ingredient", r.im.CreateIngredient).Methods("POST")
	router.HandleFunc("/api/ingredient", middleware.Options).Methods("OPTIONS")



	router.HandleFunc("/api/health", r.sm.HealthCheck).Methods("GET")
	router.HandleFunc("/api/health", middleware.Options).Methods("OPTIONS")

	router.HandleFunc("/api/household", r.hm.CreateHousehold).Methods("POST")
	router.HandleFunc("/api/household", middleware.Options).Methods("OPTIONS")
	router.HandleFunc("/api/household/{id}", r.hm.GetHousehold).Methods("GET")
	router.HandleFunc("/api/household/{id}", middleware.Options).Methods("OPTIONS")
	router.HandleFunc("/api/household/{id}/user", r.hm.AddUserToHousehold).Methods("PUT")
	router.HandleFunc("/api/household/{id}/user", middleware.Options).Methods("OPTIONS")

	return router
}
