package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/models"
	"server/router"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	r := router.Router()
	return r
}

func TestCreate(t *testing.T) {
	recipe := &models.Recipe{
		RecipeName: "recipe",
	}
	jsonRecipe, _ := json.Marshal(recipe)
	request, _ := http.NewRequest("POST", "/api/recipe", bytes.NewBuffer(jsonRecipe))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 201, response.Code, "OK response is expected")
}
