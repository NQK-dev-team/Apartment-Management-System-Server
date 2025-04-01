package services

// import (
// 	"api/config"
// 	"api/models"
// 	"api/repositories"
// 	"api/utils"
// 	"encoding/csv"
// 	"errors"
// 	"os"
// 	"path/filepath"

// 	"github.com/gin-gonic/gin"
// 	"github.com/xuri/excelize/v2"
// 	"gorm.io/gorm"
// )

// type ImportService struct {
// 	userRepository *repositories.UserRepository
// }

// func CheckValidImport(ctx *gin.Context, user *models.UserModel) bool {
// 	return true
// }

// func NewImportService() *ImportService {
// 	userRepository := repositories.NewUserRepository()
// 	return &ImportService{userRepository: userRepository}
// }

// func readUserCSVFile(file string, users *[]models.UserModel) error {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	reader := csv.NewReader(f)
// 	rows, err := reader.ReadAll()
// 	if err != nil {
// 		return err
// 	}

// 	for i, row := range rows {
// 		if i == 0 {
// 			continue
// 		}
// 		*users = append(*users, models.UserModel{
// 			FirstName:        row[0],
// 			MiddleName:       row[1],
// 			LastName:         row[2],
// 			SSN:              row[3],
// 			OldSSN:           row[4],
// 			DOB:              row[5],
// 			POB:              row[6],
// 			Email:            row[7],
// 			Password:         row[8],
// 			Phone:            row[9],
// 			SSNFrontFilePath: row[10],
// 			SSNBackFilePath:  row[11],
// 			ProfileFilePath:  row[12],
// 			IsOwner:          row[13] == "1",
// 			IsManager:        row[14] == "1",
// 			IsCustomer:       row[15] == "1",
// 		})
// 	}

// 	return nil
// }

// func readUserXLSXFile(file string, users *[]models.UserModel) error {
// 	f, err := excelize.OpenFile(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	rows, err := f.GetRows(f.GetSheetName(0))
// 	if err != nil {
// 		return err
// 	}

// 	// Ignore the first row and read the rest of the file
// 	for i, row := range rows {
// 		if i == 0 {
// 			continue
// 		}
// 		*users = append(*users, models.UserModel{
// 			FirstName:        row[0],
// 			MiddleName:       row[1],
// 			LastName:         row[2],
// 			SSN:              row[3],
// 			OldSSN:           row[4],
// 			DOB:              row[5],
// 			POB:              row[6],
// 			Email:            row[7],
// 			Password:         row[8],
// 			Phone:            row[9],
// 			SSNFrontFilePath: row[10],
// 			SSNBackFilePath:  row[11],
// 			ProfileFilePath:  row[12],
// 			IsOwner:          row[13] == "1",
// 			IsManager:        row[14] == "1",
// 			IsCustomer:       row[15] == "1",
// 		})
// 	}

// 	return nil
// }

// func getImportedFiles(pattern string) ([]string, error) {
// 	pattern = filepath.Join("assets", "imports", pattern)
// 	files, err := filepath.Glob(pattern)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return files, nil
// }

// func (s *ImportService) ImportFile(ctx *gin.Context, table int) error {
// 	if table == 1 {
// 		// Import into user table
// 		// Get the list of files from root/assets/imports/ that has this pattern: imp_user_*.*
// 		files, err := getImportedFiles("import_user_*.*")
// 		if err != nil {
// 			return err
// 		}

// 		if len(files) == 0 {
// 			return errors.New(config.GetMessageCode("NO_IMPORT_FILE"))
// 		}

// 		// Loop through the files and import them
// 		for _, file := range files {
// 			err := config.DB.Transaction(func(tx *gorm.DB) error {
// 				var err error
// 				var data []models.UserModel
// 				if filepath.Ext(file) == ".csv" {
// 					err = readUserCSVFile(file, &data)
// 				} else {
// 					err = readUserXLSXFile(file, &data)
// 				}

// 				if err != nil {
// 					return err
// 				}

// 				for _, user := range data {
// 					if !CheckValidImport(ctx, &user) {
// 						return errors.New(config.GetMessageCode("INVALID_IMPORT"))
// 					}
// 				}

// 				// Insert or update data in the database
// 				for _, user := range data {
// 					findUser := models.UserModel{}
// 					s.userRepository.GetBySSN(nil, &findUser, user.SSN)
// 					if findUser.ID != 0 {
// 						user.ID = findUser.ID
// 						if !utils.CompareHashPassword(findUser.Password, user.Password) {
// 							if user.Password, err = utils.HashPassword(user.Password); err != nil {
// 								return err
// 							}
// 						} else {
// 							user.Password = findUser.Password
// 						}
// 						if err := s.userRepository.Update(nil, tx, &user); err != nil {
// 							return err
// 						}
// 					} else {
// 						if user.Password, err = utils.HashPassword(user.Password); err != nil {
// 							return err
// 						}
// 						if err := s.userRepository.Create(nil, tx, &user); err != nil {
// 							return err
// 						}
// 					}
// 				}

// 				// Remove the imported file
// 				os.Remove(file)

// 				// Send success notification to the client via web socket

// 				return nil
// 			})

// 			if err != nil {
// 				// Send error notification to the client via web socket
// 			}
// 		}

// 	} else if table == 2 {
// 	} else if table == 3 {
// 	} else if table == 4 {
// 	}

// 	return nil
// }
