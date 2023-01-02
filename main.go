package main

import (
	"context"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"server/config"
	"server/controller"
	"server/db"
	"server/middleware"
	"server/router"
)

func main() {
	var dBFlag string
	var envFlag string

	flag.StringVar(&dBFlag, "DB_STRING", "", "Database connection string")
	flag.StringVar(&envFlag, "ENV", "", "Environment string")
	flag.Parse()

	// If the dBString is empty, then we need to fall back on a file if one is present
	if dBFlag == "" {
		fmt.Println("Environment variable not found")
		dBFlag = config.GetConfig().ConnectionString
	} else {
		fmt.Println("Environment variable found")
	}
	clientOptions := options.Client().ApplyURI(dBFlag)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	pingErr := mongoClient.Ping(context.TODO(), nil)

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to MongoDB!")

	// Get controllers with their associated DB connections
	var userController = controller.NewUserController()
	var recipeController = controller.NewRecipeController(db.NewRecipeRepository(mongoClient))
	var ingredientController = controller.NewIngredientController()
	// Check this one since it calls NewUserRepository a second time
	var authController = controller.NewAuthenticationController()
	var serverController = controller.NewServerController(mongoClient)
	var householdController = controller.NewHouseholdController(db.NewCalendarRepository(mongoClient), db.NewHouseholdRepository(mongoClient), recipeController)

	// Get middleware wrapping their controllers
	var authMiddleware = middleware.NewAuthMiddleware(authController, db.NewUserRepository(mongoClient))
	var userMiddleware = middleware.NewUserMiddleware(authMiddleware, userController, db.NewUserRepository(mongoClient))
	var recipeMiddleware = middleware.NewRecipeMiddleware(authMiddleware, recipeController)
	var ingredientMiddleware = middleware.NewIngredientMiddleware(authMiddleware, ingredientController, db.NewIngredientRepository(mongoClient))
	var serverMiddleware = middleware.NewServerMiddleware(serverController)
	var householdMiddleware = middleware.NewHouseholdMiddleware(authMiddleware,
		userMiddleware,
		householdController,
		db.NewHouseholdRepository(mongoClient),
		db.NewCalendarRepository(mongoClient))

	// If the above dependency setup starts getting much bigger we might want to look into a DI package like dig or wire
	// to more cleanly manage it

	// Build router from middleware
	var tastyRouter = router.NewTastyBoiRouter(userMiddleware, recipeMiddleware, ingredientMiddleware, serverMiddleware, householdMiddleware)
	if err != nil {
		log.Fatal(err)
	}

	// Use router to build routes with middleware
	r := tastyRouter.Route()
	fmt.Println("Starting server on the port 8080...")
	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
