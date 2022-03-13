package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"server/db"
	"server/middleware"
	"server/models"
	"strings"
	"testing"
)

type mockAuthHandler struct{}

func (m mockAuthHandler) AuthenticateUser(response http.ResponseWriter, request *http.Request, isAdmin bool) error {
	return nil
}

func (m mockAuthHandler) AuthenticateSpecificUser(response http.ResponseWriter, request *http.Request, userInfo string) error {
	return nil
}

type mockIngredientController struct{}

func (m mockIngredientController) CreateIngredient(ingredient models.Ingredient, repository db.IngredientCreator) (models.Ingredient, error) {
	return ingredient, nil
}

func (m mockIngredientController) DeleteIngredient(ingredientID string, repository db.IngredientDeleter) error {
	return nil
}

func (m mockIngredientController) GetIngredient(ingredientID string, repository db.IngredientGetter) (models.Ingredient, error) {
	return models.Ingredient{Name: "Test Ingredient"}, nil
}

func (m mockIngredientController) QueryIngredient(prefixIngredient string, repository db.IngredientGetter) ([]models.Ingredient, error) {
	return []models.Ingredient{}, nil
}

type invalidIngredientController struct{}

func (i invalidIngredientController) CreateIngredient(ingredient models.Ingredient, repository db.IngredientCreator) (models.Ingredient, error) {
	return models.Ingredient{}, errors.New("")
}

func (i invalidIngredientController) DeleteIngredient(ingredientID string, repository db.IngredientDeleter) error {
	return errors.New("")
}

func (i invalidIngredientController) GetIngredient(ingredientID string, repository db.IngredientGetter) (models.Ingredient, error) {
	return models.Ingredient{}, errors.New("")
}

func (i invalidIngredientController) QueryIngredient(prefixIngredient string, repository db.IngredientGetter) ([]models.Ingredient, error) {
	return []models.Ingredient{}, errors.New("")
}

var validMiddleware = middleware.NewIngredientMiddleware(mockAuthHandler{}, mockIngredientController{}, nil)
var errorMiddleware = middleware.NewIngredientMiddleware(mockAuthHandler{}, invalidIngredientController{}, nil)

func TestCreateIngredient(t *testing.T) {
	body := strings.NewReader("{\"name\": \"TestIngredient123\"}")
	request, _ := http.NewRequest("GET", "", body)
	response := httptest.NewRecorder()
	validMiddleware.CreateIngredient(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("Received status code %d expected 201", response.Code)
	}

	if !strings.Contains(response.Body.String(), "TestIngredient123") {
		t.Fatalf("Did not return created ingredient")
	}

}

func TestCreateIngredientBadRequest(t *testing.T) {
	body := strings.NewReader("{\"BadKey\": \"TestIngredient123\"}")
	request, _ := http.NewRequest("GET", "", body)
	response := httptest.NewRecorder()
	errorMiddleware.CreateIngredient(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Received status code %d expected 201", response.Code)
	}
}
