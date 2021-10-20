package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"server/email"
	"server/models"
)

const (
	subject = "Subject: Grocery List\r\n\r\n"
)

func SendEmail(basket models.Basket) error {
	shoppingList := buildCategoryString("Produce", basket.Produce)
	shoppingList += buildCategoryString("Pantry", basket.Pantry)
	shoppingList += buildCategoryString("Protein", basket.Protein)
	shoppingList += buildCategoryString("Dairy", basket.Dairy)
	shoppingList += buildCategoryString("Alcohol", basket.Alcohol)
	userEmail, err := lookUpUserEmail(basket.UserName)
	if err != nil {
		return err
	}
	to := []string{userEmail}
	return email.Send(subject, "", shoppingList, to)
}

func buildCategoryString(category string, ingredients []string) string {
	if len(ingredients) == 0 {
		return ""
	}
	output := category + "\n"
	for _, ingredient := range ingredients {
		output += ingredient + "\n"
	}
	output += "\n"
	return output
}

func lookUpUserEmail(userName string) (email string, error error) {
	getFilter := bson.M{"username": userName}
	user := models.User{}
	err := UserCollection.FindOne(context.Background(), getFilter).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
