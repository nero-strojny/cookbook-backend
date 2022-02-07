package controller

const (
	subject = "Subject: Grocery List\r\n\r\n"
)

//func SendEmail(basket models.Basket, token string) error {
//	shoppingList := buildCategoryString("Produce", basket.Produce)
//	shoppingList += buildCategoryString("Pantry", basket.Pantry)
//	shoppingList += buildCategoryString("Protein", basket.Protein)
//	shoppingList += buildCategoryString("Dairy", basket.Dairy)
//	shoppingList += buildCategoryString("Alcohol", basket.Alcohol)
//	userEmail, err := lookUpUserEmail(token)
//
//	if err != nil {
//		return err
//	}
//	to := []string{userEmail}
//	return email.Send(subject, "", shoppingList, to)
//}
//
//func buildCategoryString(category string, ingredients []string) string {
//	if len(ingredients) == 0 {
//		return ""
//	}
//	output := category + "\n"
//	for _, ingredient := range ingredients {
//		output += ingredient + "\n"
//	}
//	output += "\n"
//	return output
//}
//
//func lookUpUserEmail(token string) (email string, error error) {
//	getFilter := bson.M{"accesstoken": token}
//	user := models.User{}
//	err := UserCollection.FindOne(context.Background(), getFilter).Decode(&user)
//	if err != nil {
//		return "", err
//	}
//	return user.Email, nil
//}
