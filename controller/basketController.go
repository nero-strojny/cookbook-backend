package controller

import (
	"net/smtp"
	"server/models"
)

func SendEmail(basket models.Basket) error {
	from := "tasty.boi.shopping.list@gmail.com"
	pass := Config.EmailPassword
	
	to := []string{"jakestrojny@gmail.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"


	shoppingList := "Subject: Grocery List  \r\n\r\n"
	shoppingList += buildCategoryString("Produce", basket.Produce)
	shoppingList += buildCategoryString("Pantry", basket.Pantry)
	shoppingList += buildCategoryString("Protein", basket.Protein)
	shoppingList += buildCategoryString("Dairy", basket.Dairy)
	shoppingList += buildCategoryString("Alcohol", basket.Alcohol)

	message := []byte(shoppingList)

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}

func buildCategoryString(category string, ingredients []string) string {
	if len(ingredients) == 0 {
		return ""
	}
	output := "\n" + category + "\n"
	for _, ingredient := range ingredients {
		output += ingredient +"\n"
	}

	return output
}