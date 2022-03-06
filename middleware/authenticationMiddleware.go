package middleware

import (
	"net/http"
	"server/controller"
	"server/db"
	"strings"
)

type AuthMiddleware struct {
	ac         controller.AuthControl
	repository db.UserDB
}

func NewAuthMiddleware(ac controller.AuthControl, db db.UserDB) AuthMiddleware {
	return AuthMiddleware{ac, db}
}

func (am AuthMiddleware) AuthenticateUser(response http.ResponseWriter, request *http.Request, isAdmin bool) error {
	bearerToken := request.Header.Get("Authorization")
	userErr := am.ac.ValidateUser(strings.ReplaceAll(bearerToken, "Bearer ", ""), isAdmin, am.repository)
	if userErr != nil {
		if userErr.Error() == "User does not have admin permissions" {
			response.WriteHeader(http.StatusForbidden)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		response.WriteHeader(http.StatusOK)
	}
	return userErr
}

func (am AuthMiddleware) AuthenticateSpecificUser(response http.ResponseWriter, request *http.Request, userInfo string) error {
	bearerToken := request.Header.Get("Authorization")
	userErr := am.ac.ValidateSpecificUser(strings.ReplaceAll(bearerToken, "Bearer ", ""), userInfo, am.repository)
	if userErr != nil {
		response.WriteHeader(http.StatusUnauthorized)
	}
	return userErr
}
