package utils

import "api/models"

func GetUserFullName(user *models.UserModel) string {
	fullName := ""
	if user.MiddleName.Valid {
		fullName = user.LastName + " " + user.MiddleName.String + " " + user.FirstName
	} else {
		fullName = user.LastName + " " + user.FirstName
	}
	return fullName
}
