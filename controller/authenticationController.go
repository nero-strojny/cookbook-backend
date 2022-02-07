package controller

import (
	"errors"
	"server/db"
	"time"
)

type AuthController struct {
}

type AuthControl interface {
	ValidateUser(accessToken string, restrictAdmin bool, repository db.UserGetter) error
	ValidateSpecificUser(accessToken string, userName string, repository db.UserGetter) error
}

func NewAuthenticationController() AuthController {
	return AuthController{}
}

//ValidateUser
func (ac AuthController) ValidateUser(accessToken string, restrictAdmin bool, repository db.UserGetter) error {
	if len(accessToken) == 0 {
		return errors.New("no token in request")
	}
	user, err := repository.GetUserByAccessToken(accessToken)
	currentTime := time.Now().Format("2006.01.02 15:04:05")
	if err != nil {
		return err
	} else if len(user.ExpiryDate) == 0 || user.ExpiryDate < currentTime {
		return errors.New("expired Token")
	} else if restrictAdmin && user.UserType != "admin" {
		return errors.New("user does not have admin permissions")
	} else {
		return nil
	}
}

//ValidateSpecificUser
func (ac AuthController) ValidateSpecificUser(accessToken string, userName string, repository db.UserGetter) error {
	if len(accessToken) == 0 {
		return errors.New("No token in request")
	}
	currentTime := time.Now().Format("2006.01.02 15:04:05")
	user, err := repository.GetUserByAccessToken(accessToken)
	if err != nil || len(user.ExpiryDate) == 0 || user.ExpiryDate < currentTime {
		return err
	} else if user.UserName != userName {
		return errors.New("Invalid permissions")
	} else {
		return nil
	}
}
