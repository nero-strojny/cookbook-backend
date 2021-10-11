package controller

import (
	"server/email"
	"server/models"
)

const (
	mime    = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: Grocery List\r\n\r\n"
)

func SendEmail(basket models.Basket) error {
	shoppingList := buildCategoryString("Produce", basket.Produce)
	shoppingList += buildCategoryString("Pantry", basket.Pantry)
	shoppingList += buildCategoryString("Protein", basket.Protein)
	shoppingList += buildCategoryString("Dairy", basket.Dairy)
	shoppingList += buildCategoryString("Alcohol", basket.Alcohol)
	to := []string{"jakestrojny@gmail.com"}
	return email.Send(subject, mime, shoppingList, to)
}

func buildCategoryString(category string, ingredients []string) string {
	if len(ingredients) == 0 {
		return ""
	}
	output := "<b>" + category + "</b>" + "\n"
	output += "<ul>"
	for _, ingredient := range ingredients {
		output += "<li>" + ingredient + "</lu>" + "\n"
	}
	output += "</ul>"
	return output
}
