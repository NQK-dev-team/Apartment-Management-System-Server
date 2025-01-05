package services

import (
	"api/config"
	"api/models"
	"api/repositories"
	"api/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type ImportServiceDev struct {
	userRepository *repositories.UserRepository
}

func NewImportServiceDev() *ImportServiceDev {
	userRepository := repositories.NewUserRepository()
	return &ImportServiceDev{userRepository: userRepository}
}

func (s *ImportServiceDev) ImportFileDev(table int) error {
	if table == 1 {
		// Import into user table
		// Get the list of files from root/assets/imports/ that has this pattern: imp_user_*.*
		files, err := getImportedFiles("import_user_*.*")
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return errors.New(config.GetMessageCode("NO_IMPORT_FILE"))
		}

		// Loop through the files and import them
		for _, file := range files {
			err := config.DB.Transaction(func(tx *gorm.DB) error {
				var err error
				var data []models.UserModel
				if filepath.Ext(file) == ".csv" {
					err = readUserCSVFile(file, &data)
				} else {
					err = readUserXLSXFile(file, &data)
				}

				if err != nil {
					return err
				}

				// Insert or update data in the database
				for _, user := range data {
					user.EmailVerifiedAt.Time = time.Now()
					user.EmailVerifiedAt.Valid = true
					findUser := models.UserModel{}
					s.userRepository.GetBySSN(nil, &findUser, user.SSN)
					if findUser.ID != 0 {
						user.ID = findUser.ID
						user.EmailVerifiedAt = findUser.EmailVerifiedAt
						if !utils.CompareHashPassword(findUser.Password, user.Password) {
							if user.Password, err = utils.HashPassword(user.Password); err != nil {
								return err
							}
						} else {
							user.Password = findUser.Password
						}
						if err := s.userRepository.Update(nil, &user); err != nil {
							return err
						}
					} else {
						if user.Password, err = utils.HashPassword(user.Password); err != nil {
							return err
						}
						if err := s.userRepository.Create(nil, &user); err != nil {
							return err
						}
					}
				}

				// Remove the imported file
				os.Remove(file)

				return nil
			})

			if err != nil {
				fmt.Printf("Error while importing file %s: %v\n", file, err)
			}
		}

	} else if table == 2 {
	} else if table == 3 {
	} else if table == 4 {
	}

	return nil
}
