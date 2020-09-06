package test

import (
	"server/models"
	"server/router"
	"testing"
)

var userRouter = router.Router()

var defaultUser = models.Recipe{
	RecipeName: "recipe",
	Private:    false,
	Author:     "testUser",
	Rating:     5,
	Servings:   2,
	Calories:   500,
	PrepTime:   5,
	CookTime:   5,
}

func userFieldsAreExpected(t *testing.T, recipe1 models.Recipe, recipe2 models.Recipe) {

}
