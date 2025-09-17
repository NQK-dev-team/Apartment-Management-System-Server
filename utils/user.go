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

func GetUserRole(user *models.UserModel) string {
	var str = ""
	if user.IsOwner {
		str += "1"
	} else {
		str += "0"
	}

	if user.IsManager {
		str += "1"
	} else {
		str += "0"
	}

	if user.IsCustomer {
		str += "1"
	} else {
		str += "0"
	}

	return str
}
