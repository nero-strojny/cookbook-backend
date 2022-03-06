package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"server/db"
	"server/middleware"
	"testing"
)

type mockAuthControl struct{}

func (m mockAuthControl) ValidateUser(accessToken string, restrictAdmin bool, repository db.UserGetter) error {
	return nil
}

func (m mockAuthControl) ValidateSpecificUser(accessToken string, userName string, repository db.UserGetter) error {
	return nil
}

type unauthorizedUserControl struct{}

func (u unauthorizedUserControl) ValidateUser(accessToken string, restrictAdmin bool, repository db.UserGetter) error {
	return errors.New("unauthorized user")
}

func (u unauthorizedUserControl) ValidateSpecificUser(accessToken string, userName string, repository db.UserGetter) error {
	return errors.New("unauthorized user")
}

type nonAdminUserControl struct{}

func (n nonAdminUserControl) ValidateUser(accessToken string, restrictAdmin bool, repository db.UserGetter) error {
	return errors.New("User does not have admin permissions")
}

func (n nonAdminUserControl) ValidateSpecificUser(accessToken string, userName string, repository db.UserGetter) error {
	return errors.New("User does not have admin permissions")
}

func TestValidAuth(t *testing.T) {
	am := middleware.NewAuthMiddleware(mockAuthControl{}, nil)
	req, _ := http.NewRequest("GET", "Test", nil)
	rr := httptest.NewRecorder()
	authErr := am.AuthenticateUser(rr, req, false)
	if authErr != nil {
		t.Fatal("Error occurred for auth when it shouldn't have")
	}
	if rr.Code != http.StatusOK {
		t.Fatal("Response did not have OK status")
	}
}

func TestUnauthorizedUser(t *testing.T) {
	am := middleware.NewAuthMiddleware(unauthorizedUserControl{}, nil)
	req, _ := http.NewRequest("GET", "Test", nil)
	rr := httptest.NewRecorder()
	authErr := am.AuthenticateUser(rr, req, false)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected StatusUnauthorized got %d", rr.Code)
	}

	if authErr == nil {
		t.Fatal("Error should have been returned for unauthorized user")
	}
}

func TestNonAdminUserForAdminFunction(t *testing.T) {
	am := middleware.NewAuthMiddleware(nonAdminUserControl{}, nil)
	req, _ := http.NewRequest("GET", "Test", nil)
	rr := httptest.NewRecorder()
	authErr := am.AuthenticateUser(rr, req, false)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("Expected StatusForbidden received %d", rr.Code)
	}

	if authErr == nil {
		t.Fatal("Error should have been returned for nonadmin using admin endpoint")
	}
}
