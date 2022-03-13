package test

import (
	"net/http"
	"net/http/httptest"
	"server/router"
	"testing"
)

// Authentication middleware mocking
type mockAuthMiddleware struct {
}

func (m mockAuthMiddleware) AuthenticateUser(response http.ResponseWriter, request *http.Request, isAdmin bool) error {
	return nil
}

func (m mockAuthMiddleware) AuthenticateSpecificUser(response http.ResponseWriter, request *http.Request, userInfo string) error {
	return nil
}

type mockHouseholdMiddleware struct {
}

func (m mockHouseholdMiddleware) CreateHousehold(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockHouseholdMiddleware) GetHousehold(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockHouseholdMiddleware) AddUserToHousehold(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockHouseholdMiddleware) GetCalendar(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockHouseholdMiddleware) UpdateCalendar(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockHouseholdMiddleware) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type mockIngredientMiddleware struct {
}

func (m mockIngredientMiddleware) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockIngredientMiddleware) GetIngredient(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockIngredientMiddleware) QueryIngredient(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockIngredientMiddleware) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type mockRecipeMiddleware struct {
}

func (m mockRecipeMiddleware) PostPaginatedRecipes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockRecipeMiddleware) GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockRecipeMiddleware) GetRandomRecipes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockRecipeMiddleware) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockRecipeMiddleware) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockRecipeMiddleware) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type mockServerMiddleware struct {
}

func (m mockServerMiddleware) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type mockUserMiddleware struct {
}

func (m mockUserMiddleware) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockUserMiddleware) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockUserMiddleware) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockUserMiddleware) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockUserMiddleware) GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (m mockUserMiddleware) EmailUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var r = router.NewTastyBoiRouter(mockUserMiddleware{},
	mockRecipeMiddleware{},
	mockIngredientMiddleware{},
	mockServerMiddleware{},
	mockHouseholdMiddleware{}).Route()

func TestRecipesPostRoute(t *testing.T) {
	request, _ := http.NewRequest("POST", "/api/recipes", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestRecipesOptionsRoute(t *testing.T) {
	request, _ := http.NewRequest("OPTIONS", "/api/recipes", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestSpecificRecipeGetRoute(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/recipe/recipeID1234", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestSpecificRecipeOptionRoute(t *testing.T) {
	request, _ := http.NewRequest("OPTION", "/api/recipe/recipeID1234", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestSpecificRecipeDeleteRoute(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/api/recipe/recipeID1234", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestSpecificRecipePutRoute(t *testing.T) {
	request, _ := http.NewRequest("PUT", "/api/recipe/recipeID1234", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestRecipePostRoute(t *testing.T) {
	request, _ := http.NewRequest("POST", "/api/recipe", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestRecipeOptionRoute(t *testing.T) {
	request, _ := http.NewRequest("OPTIONS", "/api/recipe", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestRandomRecipeGetRoute(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/randomRecipe/10", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}

func TestRandomRecipeOptionRoute(t *testing.T) {
	request, _ := http.NewRequest("OPTIONS", "/api/randomRecipe/10", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal("Route returned non-200, likely not registered correctly.")
	}
}
