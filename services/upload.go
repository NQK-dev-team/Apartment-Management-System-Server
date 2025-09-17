package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type UploadService struct {
	repository         *repositories.UploadRepository
	userRepository     *repositories.UserRepository
	buildingRepository *repositories.BuildingRepository
	contractRepository *repositories.ContractRepository
	billRepository     *repositories.BillRepository
	roomRepository     *repositories.RoomRepository
	emailService       *EmailService
}

func NewUploadService() *UploadService {
	return &UploadService{
		repository:         repositories.NewUploadRepository(),
		userRepository:     repositories.NewUserRepository(),
		buildingRepository: repositories.NewBuildingRepository(),
		contractRepository: repositories.NewContractRepository(),
		billRepository:     repositories.NewBillRepository(),
		roomRepository:     repositories.NewRoomRepository(),
		emailService:       NewEmailService(),
	}
}

func (s *UploadService) UploadFile(ctx *gin.Context, upload *structs.UploadStruct) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)

		uploadModel := &models.UploadFileModel{
			CreatorID: ctx.GetInt64("userID"),
			FileName:  upload.File.Filename,
			URLPath:   "",
			// StoragePath: "",
			Size:       upload.File.Size,
			UploadType: upload.UploadType,
		}

		if err := s.repository.Create(ctx, tx, uploadModel); err != nil {
			return err
		}

		uploadIDStr := strconv.FormatInt(uploadModel.ID, 10)

		filePath, err := utils.StoreFileSingleMedia(upload.File, constants.GetUploadFileURL("files", uploadIDStr, ""))
		if err != nil {
			utils.RemoveFile(filePath)
			return err
		}

		uploadModel.URLPath = filePath
		// uploadModel.StoragePath = strings.ReplaceAll(filePath, "/api/", "")

		if err := s.repository.Update(ctx, tx, uploadModel); err != nil {
			return err
		}

		return nil
	})
}

func (s *UploadService) GetUploads(ctx *gin.Context, uploads *[]models.UploadFileModel, uploadType int, isProcessed bool, date string) error {
	return s.repository.Get(ctx, uploads, uploadType, isProcessed, date)
}

func (s *UploadService) CronFileFail(tx *gorm.DB, upload *models.UploadFileModel) error {
	upload.ProcessResult = sql.NullInt64{
		Int64: constants.Common.CronUploadProcessResult.FAILED,
		Valid: true,
	}

	upload.ProcessDate = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	return s.repository.Update(nil, tx, upload)
}

func (s *UploadService) CronFileSuccess(tx *gorm.DB, upload *models.UploadFileModel) error {
	upload.ProcessResult = sql.NullInt64{
		Int64: constants.Common.CronUploadProcessResult.SUCCESS,
		Valid: true,
	}

	upload.ProcessDate = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	return s.repository.Update(nil, tx, upload)
}

func (s *UploadService) ProcessAddCustomer(f *excelize.File, tx *gorm.DB, upload *models.UploadFileModel) []error {
	failCell, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#FF0000", // Red
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1}, // Style 1 = thin
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	successCell, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#00FF00", // Green
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1}, // Style 1 = thin
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})

	listSheet := f.GetSheetName(0)
	if listSheet == "" {
		return []error{errors.New("first sheet not found")}
	}

	rows, err := f.GetRows(listSheet)
	if err != nil {
		return []error{err}
	}

	if len(rows) < 6 {
		return []error{errors.New("no data rows found")}
	}

	dataRows := rows[5:]

	fileError := []error{}
	type emailQueueStruct struct {
		Email    string
		FullName string
		Password string
	}
	emailQueue := []emailQueueStruct{}

	for index := range dataRows {
		processSuccess := true
		row := index + 6

		// Get row data
		lastName, _ := f.GetCellValue(listSheet, fmt.Sprintf("B%d", row))
		middleName, _ := f.GetCellValue(listSheet, fmt.Sprintf("C%d", row))
		firstName, _ := f.GetCellValue(listSheet, fmt.Sprintf("D%d", row))
		genderStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("E%d", row))
		ssn, _ := f.GetCellValue(listSheet, fmt.Sprintf("F%d", row))
		oldSSN, _ := f.GetCellValue(listSheet, fmt.Sprintf("G%d", row))
		dobStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("H%d", row))
		pob, _ := f.GetCellValue(listSheet, fmt.Sprintf("I%d", row))
		email, _ := f.GetCellValue(listSheet, fmt.Sprintf("J%d", row))
		phone, _ := f.GetCellValue(listSheet, fmt.Sprintf("K%d", row))
		permanentAddress, _ := f.GetCellValue(listSheet, fmt.Sprintf("L%d", row))
		tempAddress, _ := f.GetCellValue(listSheet, fmt.Sprintf("M%d", row))

		// Process customer gender string
		gender := 0
		if utils.CompareStringRaw("Nam", strings.TrimSpace(genderStr)) || utils.CompareStringRaw("Male", strings.TrimSpace(genderStr)) {
			gender = constants.Common.UserGender.MALE
		} else if utils.CompareStringRaw("Nữ", strings.TrimSpace(genderStr)) || utils.CompareStringRaw("Female", strings.TrimSpace(genderStr)) {
			gender = constants.Common.UserGender.FEMALE
		} else if utils.CompareStringRaw("Khác", strings.TrimSpace(genderStr)) || utils.CompareStringRaw("Other", strings.TrimSpace(genderStr)) {
			gender = constants.Common.UserGender.OTHER
		}

		// Generate customer's password
		newPassword, err := utils.GeneratePassword(constants.Common.NewPasswordLength)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to generate password", row-5))
			newPassword = "123456"
		}
		hashedPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to hash password", row-5))
			hashedPassword = "$12$xG0qlWDqXflwTqTBgFRnjuA1J5zZRSd6dbzOT353TAQS7ScjVfqXW"
		}

		// Create new customer's account
		newUser := &structs.NewUploadCustomer{
			LastName:         strings.TrimSpace(lastName),
			FirstName:        strings.TrimSpace(firstName),
			MiddleName:       strings.TrimSpace(middleName),
			Dob:              strings.TrimSpace(dobStr),
			Pob:              strings.TrimSpace(pob),
			Gender:           gender,
			SSN:              strings.TrimSpace(ssn),
			OldSSN:           strings.TrimSpace(oldSSN),
			Email:            strings.TrimSpace(email),
			Phone:            strings.TrimSpace(phone),
			PermanentAddress: strings.TrimSpace(permanentAddress),
			TemporaryAddress: strings.TrimSpace(tempAddress),
			ProfileImage:     "/image/placeholder_image.png",
			FrontSSNImage:    "/image/placeholder_image.png",
			BackSSNImage:     "/image/placeholder_image.png",
		}

		if err := constants.Validate.Struct(newUser); err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
			processSuccess = false
		}

		if processSuccess {
			dob, _ := time.Parse("2006-01-02", newUser.Dob)

			customerModel := &models.UserModel{
				Password:         hashedPassword,
				IsOwner:          false,
				IsManager:        false,
				IsCustomer:       true,
				LastName:         newUser.LastName,
				MiddleName:       sql.NullString{String: newUser.MiddleName, Valid: newUser.MiddleName != ""},
				FirstName:        newUser.FirstName,
				DOB:              dob,
				POB:              newUser.Pob,
				Gender:           newUser.Gender,
				SSN:              newUser.SSN,
				OldSSN:           sql.NullString{String: newUser.OldSSN, Valid: newUser.OldSSN != ""},
				Email:            newUser.Email,
				Phone:            newUser.Phone,
				PermanentAddress: newUser.PermanentAddress,
				TemporaryAddress: newUser.TemporaryAddress,
				ProfileFilePath:  newUser.ProfileImage,
				SSNFrontFilePath: newUser.FrontSSNImage,
				SSNBackFilePath:  newUser.BackSSNImage,
			}

			ctx := &gin.Context{}
			ctx.Set("userID", upload.CreatorID)
			tx.SavePoint(fmt.Sprintf("sp_customer_%d", row-5))
			if err := s.userRepository.Create(ctx, tx, customerModel); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: failed to create user: %v", row-5, err))
				processSuccess = false
				tx.RollbackTo(fmt.Sprintf("sp_customer_%d", row-5))
			}

			if processSuccess {
				emailQueue = append(emailQueue, emailQueueStruct{
					Email:    customerModel.Email,
					FullName: utils.GetUserFullName(customerModel),
					Password: newPassword,
				})
			}
		}

		if !processSuccess {
			f.SetCellValue(listSheet, fmt.Sprintf("N%d", row), "✘")
			f.SetCellStyle(listSheet, fmt.Sprintf("N%d", row), fmt.Sprintf("N%d", row), failCell)
		} else {
			f.SetCellValue(listSheet, fmt.Sprintf("N%d", row), "✔")
			f.SetCellStyle(listSheet, fmt.Sprintf("N%d", row), fmt.Sprintf("N%d", row), successCell)
		}
	}

	if len(fileError) == 0 {
		for _, elem := range emailQueue {
			s.emailService.SendAccountCreationEmail(config.GetEnv("APM_CLIENT_BASE_URL")+"/login", elem.Email, elem.FullName, elem.Password)
		}
	}

	return fileError
}

func (s *UploadService) CheckUserPermissionForBuilding(ctx *gin.Context, buildingID int64) bool {
	userID := ctx.GetInt64("userID")

	userModel := &models.UserModel{}
	if err := s.userRepository.GetByID(ctx, userModel, userID); err != nil {
		return false
	}

	role := utils.GetUserRole(userModel)

	if role != constants.Roles.Manager && role != constants.Roles.Owner {
		return false
	}

	if role == constants.Roles.Manager {
		buildings := []models.BuildingModel{}

		if err := s.buildingRepository.GetBuildingBaseOnSchedule(ctx, &buildings, userID); err != nil {
			return false
		}

		if len(buildings) == 0 {
			return false
		}

		var result = false

		for _, building := range buildings {
			if building.ID == buildingID {
				result = true
				break
			}
		}

		return result
	}

	return true
}

func (s *UploadService) ProcessAddContract(f *excelize.File, tx *gorm.DB, upload *models.UploadFileModel) []error {
	failCell, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#FF0000", // Red
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1}, // Style 1 = thin
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	successCell, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#00FF00", // Green
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1}, // Style 1 = thin
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})

	listSheet := f.GetSheetName(0)
	if listSheet == "" {
		return []error{errors.New("first sheet not found")}
	}

	rows, err := f.GetRows(listSheet)
	if err != nil {
		return []error{err}
	}

	if len(rows) < 6 {
		return []error{errors.New("no data rows found")}
	}

	ctx := &gin.Context{}
	ctx.Set("userID", upload.CreatorID)

	buildingInfoRow := rows[0]
	// buildingName := strings.TrimSpace(buildingInfoRow[1])
	buildingIDStr := strings.TrimSpace(buildingInfoRow[3])
	buildingID, err := strconv.ParseInt(buildingIDStr, 10, 64)
	if err != nil {
		return []error{err}
	}

	buildingModel := &models.BuildingModel{}
	if err := s.buildingRepository.GetById(nil, buildingModel, buildingID); err != nil {
		return []error{err}
	} else if buildingModel.ID == 0 { // || !utils.CompareStringRaw(strings.ReplaceAll(buildingModel.Name, " ", ""), strings.ReplaceAll(buildingName, " ", ""))
		return []error{fmt.Errorf("building with ID %d not found", buildingID)}
	}

	if !s.CheckUserPermissionForBuilding(ctx, buildingID) {
		return []error{fmt.Errorf("user: %d does not have permission for building ID: %d", upload.CreatorID, buildingID)}
	}

	dataRows := rows[5:]

	fileError := []error{}

	for index := range dataRows {
		processSuccess := true
		row := index + 6

		// floorStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("B%d", row))
		// // Get room floor
		// floor, err := strconv.ParseInt(strings.TrimSpace(floorStr), 10, 64)
		// if err != nil {
		// 	fileError = append(fileError, fmt.Errorf("row %d: failed to parse room floor: %v", row-5, err))
		// 	processSuccess = false
		// }

		// Get room no
		roomNoStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("C%d", row))
		roomNo, err := strconv.ParseInt(strings.TrimSpace(roomNoStr), 10, 64)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to parse room no: %v", row-5, err))
			processSuccess = false
		} else {
			roomModel := &models.RoomModel{}

			if err := s.roomRepository.GetRoomInfoForUpload(roomModel, buildingID, int(roomNo)); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: failed to get room info: %v", row-5, err))
				processSuccess = false
			} else if roomModel.ID == 0 {
				fileError = append(fileError, fmt.Errorf("row %d: failed to get info of room no: %d", row-5, roomNo))
				processSuccess = false
			} else if roomModel.Status != constants.Common.RoomStatus.AVAILABLE && roomModel.Status != constants.Common.RoomStatus.RENTED {
				fileError = append(fileError, fmt.Errorf("row %d: room is unavailable, status: %d", row-5, roomModel.Status))
				processSuccess = false
			} else {
				// Get customer info
				customerNo, _ := f.GetCellValue(listSheet, fmt.Sprintf("J%d", row))
				customerNo = strings.TrimSpace(customerNo)
				// customerName, _ := f.GetCellValue(listSheet, fmt.Sprintf("K%d", row))
				// customerName = strings.TrimSpace(customerName)
				customerModel := &models.UserModel{}

				if err := s.userRepository.GetCustomerByUserNo(nil, customerModel, customerNo); err != nil {
					fileError = append(fileError, fmt.Errorf("row %d: failed to get customer info: %v", row-5, err))
					processSuccess = false
				} else if customerModel.ID == 0 { //|| !utils.CompareStringRaw(strings.ReplaceAll(utils.GetUserFullName(customerModel), " ", ""), customerName) {
					fileError = append(fileError, fmt.Errorf("row %d: customer of number %s not found", row-5, customerNo))
					processSuccess = false
				}

				if processSuccess {
					// Get contract no
					contractNo, _ := f.GetCellValue(listSheet, fmt.Sprintf("A%d", row))
					contractNo = strings.TrimSpace(contractNo)

					// Get contract value
					valueStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("D%d", row))
					value, err := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(valueStr, ",", "")), 64)
					if err != nil {
						fileError = append(fileError, fmt.Errorf("row %d: failed to parse contract value: %v", row-5, err))
						processSuccess = false
					} else {
						// Get contract type
						contractTypeStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("E%d", row))
						contractType := 0
						if utils.CompareStringRaw("Thuê", strings.TrimSpace(contractTypeStr)) || utils.CompareStringRaw("Rental", strings.TrimSpace(contractTypeStr)) {
							contractType = constants.Common.ContractType.RENT
						} else if utils.CompareStringRaw("Mua", strings.TrimSpace(contractTypeStr)) || utils.CompareStringRaw("Purchase", strings.TrimSpace(contractTypeStr)) {
							contractType = constants.Common.ContractType.BUY
						}

						createdAtStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("F%d", row))
						startDateStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("G%d", row))
						endDateStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("H%d", row))
						signDateStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("I%d", row))

						newContract := &structs.NewUploadContract{
							ContractType:  contractType,
							ContractValue: value,
							HouseholderID: customerModel.ID,
							RoomID:        roomModel.ID,
							CreatorID:     upload.CreatorID,
							CreatedAt:     strings.TrimSpace(createdAtStr),
							StartDate:     strings.TrimSpace(startDateStr),
							EndDate:       strings.TrimSpace(endDateStr),
							SignDate:      strings.TrimSpace(signDateStr),
						}

						if err := constants.Validate.Struct(newContract); err != nil {
							fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
							processSuccess = false
						} else {
							startDate, _ := time.Parse("2006-01-02", newContract.StartDate)
							createAt, _ := time.Parse("2006-01-02", newContract.CreatedAt)
							endDate := utils.StringToNullTime(newContract.EndDate)
							signDate := utils.StringToNullTime(newContract.SignDate)

							status := 0
							currentDate := time.Now().Format("2006-01-02")

							if newContract.SignDate == "" {
								result, err := utils.CompareDates(newContract.StartDate, currentDate)
								if err != nil {
									fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
									processSuccess = false
								} else {
									switch result {
									case 1:
										status = constants.Common.ContractStatus.WAITING_FOR_SIGNATURE
									case 0, -1:
										status = constants.Common.ContractStatus.CANCELLED
									}
								}
							} else {
								result, err := utils.CompareDates(newContract.StartDate, currentDate)
								if err != nil {
									fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
									processSuccess = false
								} else {
									switch result {
									case 1:
										status = constants.Common.ContractStatus.NOT_IN_EFFECT
									case 0, -1:
										if newContract.EndDate == "" {
											status = constants.Common.ContractStatus.ACTIVE
										} else {
											result, err := utils.CompareDates(newContract.EndDate, currentDate)
											if err != nil {
												fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
												processSuccess = false
											} else {
												switch result {
												case -1:
													status = constants.Common.ContractStatus.EXPIRED
												case 0, 1:
													status = constants.Common.ContractStatus.ACTIVE
												}
											}
										}
									}
								}
							}

							if processSuccess {
								contracts := []models.ContractModel{}
								if err := s.contractRepository.GetOverlapContract(ctx, &contracts, roomModel.ID, newContract.StartDate); err != nil {
									fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
									processSuccess = false
								}

								if len(contracts) > 0 {
									fileError = append(fileError, fmt.Errorf("row %d: new contract overlaps with existing contracts in room %d", row-5, roomModel.No))
									processSuccess = false
								} else if newContract.EndDate != "" {
									if err := s.contractRepository.GetOverlapContract(ctx, &contracts, roomModel.ID, newContract.EndDate); err != nil {
										fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
										processSuccess = false
									}

									if len(contracts) > 0 {
										fileError = append(fileError, fmt.Errorf("row %d: new contract overlaps with existing contracts in room %d", row-5, roomModel.No))
										processSuccess = false
									}
								}

								if processSuccess {
									contractModel := &models.ContractModel{
										DefaultModel: models.DefaultModel{
											CreatedAt: createAt,
											CreatedBy: upload.CreatorID,
										},
										Value:         newContract.ContractValue,
										Type:          newContract.ContractType,
										StartDate:     startDate,
										EndDate:       endDate,
										SignDate:      signDate,
										CreatorID:     newContract.CreatorID,
										HouseholderID: newContract.HouseholderID,
										RoomID:        newContract.RoomID,
										Status:        status,
									}

									tx.SavePoint(fmt.Sprintf("sp_contract_%d", row-5))
									if err := s.contractRepository.CreateContract(ctx, tx, contractModel); err != nil {
										fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
										processSuccess = false
										tx.RollbackTo(fmt.Sprintf("sp_contract_%d", row-5))
									} else {
										residentSheetRows, err := f.GetRows(contractNo)
										residentRows := [][]string{}
										if err == nil {
											if len(residentSheetRows) >= 6 {
												residentRows = residentSheetRows[5:]
											}
										}

										if len(residentRows) > 0 {
											for index2 := range residentRows {
												residentRowNo := index2 + 6
												subProcessSuccess := true

												residentUserModel := &models.UserModel{}
												residentNo, _ := f.GetCellValue(contractNo, fmt.Sprintf("B%d", residentRowNo))
												residentNo = strings.TrimSpace(residentNo)
												residentLastName, _ := f.GetCellValue(contractNo, fmt.Sprintf("C%d", residentRowNo))
												residentLastName = strings.TrimSpace(residentLastName)
												residentMiddleName, _ := f.GetCellValue(contractNo, fmt.Sprintf("D%d", residentRowNo))
												residentMiddleName = strings.TrimSpace(residentMiddleName)
												residentFirstName, _ := f.GetCellValue(contractNo, fmt.Sprintf("E%d", residentRowNo))
												residentFirstName = strings.TrimSpace(residentFirstName)

												residentGenderStr, _ := f.GetCellValue(contractNo, fmt.Sprintf("F%d", residentRowNo))
												residentGender := 0
												if utils.CompareStringRaw("Nam", strings.TrimSpace(residentGenderStr)) || utils.CompareStringRaw("Male", strings.TrimSpace(residentGenderStr)) {
													residentGender = constants.Common.UserGender.MALE
												} else if utils.CompareStringRaw("Nữ", strings.TrimSpace(residentGenderStr)) || utils.CompareStringRaw("Female", strings.TrimSpace(residentGenderStr)) {
													residentGender = constants.Common.UserGender.FEMALE
												} else if utils.CompareStringRaw("Khác", strings.TrimSpace(residentGenderStr)) || utils.CompareStringRaw("Other", strings.TrimSpace(residentGenderStr)) {
													residentGender = constants.Common.UserGender.OTHER
												}

												residentSSN, _ := f.GetCellValue(contractNo, fmt.Sprintf("G%d", residentRowNo))
												residentSSN = strings.TrimSpace(residentSSN)
												residentOldSSN, _ := f.GetCellValue(contractNo, fmt.Sprintf("H%d", residentRowNo))
												residentOldSSN = strings.TrimSpace(residentOldSSN)
												residentDOB, _ := f.GetCellValue(contractNo, fmt.Sprintf("I%d", residentRowNo))
												residentDOB = strings.TrimSpace(residentDOB)
												residentPOB, _ := f.GetCellValue(contractNo, fmt.Sprintf("J%d", residentRowNo))
												residentPOB = strings.TrimSpace(residentPOB)
												residentEmail, _ := f.GetCellValue(contractNo, fmt.Sprintf("K%d", residentRowNo))
												residentEmail = strings.TrimSpace(residentEmail)
												residentPhone, _ := f.GetCellValue(contractNo, fmt.Sprintf("L%d", residentRowNo))
												residentPhone = strings.TrimSpace(residentPhone)

												residentRelationStr, _ := f.GetCellValue(contractNo, fmt.Sprintf("M%d", residentRowNo))
												residentRelation := 0
												if utils.CompareStringRaw("Con cái", strings.TrimSpace(residentRelationStr)) || utils.CompareStringRaw("Child", strings.TrimSpace(residentRelationStr)) {
													residentRelation = constants.Common.ResidentRelationship.CHILD
												} else if utils.CompareStringRaw("Vợ/Chồng", strings.TrimSpace(residentRelationStr)) || utils.CompareStringRaw("Spouse", strings.TrimSpace(residentRelationStr)) {
													residentRelation = constants.Common.ResidentRelationship.SPOUSE
												} else if utils.CompareStringRaw("Cha/Mẹ", strings.TrimSpace(residentRelationStr)) || utils.CompareStringRaw("Parent", strings.TrimSpace(residentRelationStr)) {
													residentRelation = constants.Common.ResidentRelationship.PARENT
												} else if utils.CompareStringRaw("Khác", strings.TrimSpace(residentRelationStr)) || utils.CompareStringRaw("Other", strings.TrimSpace(residentRelationStr)) {
													residentRelation = constants.Common.ResidentRelationship.OTHER
												}

												if residentNo != "" {
													if err := s.userRepository.GetCustomerByUserNo(nil, residentUserModel, residentNo); err != nil {
														fileError = append(fileError, fmt.Errorf("row %d, resident %d: failed to get resident's customer account info: %v", row-5, residentRowNo-5, err))
														processSuccess = false
														subProcessSuccess = false
													} else {
														if residentUserModel.ID == 0 {
															fileError = append(fileError, fmt.Errorf("row %d, resident %d: customer of number %s not found", row-5, residentRowNo-5, residentNo))
															processSuccess = false
															subProcessSuccess = false
														}
													}
												} else {
													residentData := &structs.ContractResidents{
														FirstName:               residentFirstName,
														LastName:                residentLastName,
														MiddleName:              residentMiddleName,
														SSN:                     residentSSN,
														OldSSN:                  residentOldSSN,
														DOB:                     residentDOB,
														POB:                     residentPOB,
														Phone:                   residentPhone,
														Email:                   residentEmail,
														Gender:                  residentGender,
														RelationWithHouseholder: residentRelation,
													}

													if err := constants.Validate.Struct(residentData); err != nil {
														fileError = append(fileError, fmt.Errorf("row %d, resident %d: %v", row-5, residentRowNo-5, err))
														processSuccess = false
														subProcessSuccess = false
													}
												}

												if subProcessSuccess {
													var residentData *models.RoomResidentModel

													if residentUserModel.ID != 0 {
														residentData = &models.RoomResidentModel{
															FirstName:               residentUserModel.FirstName,
															LastName:                residentUserModel.LastName,
															MiddleName:              residentUserModel.MiddleName,
															SSN:                     sql.NullString{String: residentUserModel.SSN, Valid: residentUserModel.SSN != ""},
															OldSSN:                  residentUserModel.OldSSN,
															DOB:                     residentUserModel.DOB,
															POB:                     residentUserModel.POB,
															Phone:                   sql.NullString{String: residentUserModel.Phone, Valid: residentUserModel.Phone != ""},
															Email:                   sql.NullString{String: residentUserModel.Email, Valid: residentUserModel.Email != ""},
															Gender:                  residentUserModel.Gender,
															RelationWithHouseholder: residentRelation,
															UserAccountID:           sql.NullInt64{Int64: residentUserModel.ID, Valid: true},
														}
													} else {
														dob, _ := time.Parse("2006-01-02", residentDOB)

														residentData = &models.RoomResidentModel{
															FirstName:               residentFirstName,
															LastName:                residentLastName,
															MiddleName:              sql.NullString{String: residentMiddleName, Valid: residentMiddleName != ""},
															SSN:                     sql.NullString{String: residentSSN, Valid: residentSSN != ""},
															OldSSN:                  sql.NullString{String: residentOldSSN, Valid: residentOldSSN != ""},
															DOB:                     dob,
															POB:                     residentPOB,
															Phone:                   sql.NullString{String: residentPhone, Valid: residentPhone != ""},
															Email:                   sql.NullString{String: residentEmail, Valid: residentEmail != ""},
															Gender:                  residentGender,
															RelationWithHouseholder: residentRelation,
															UserAccountID:           sql.NullInt64{Int64: 0, Valid: false},
														}
													}

													tx.SavePoint(fmt.Sprintf("sp_contract_%d_resident_%d", row-5, residentRowNo-5))
													if err := s.contractRepository.AddNewRoomResident(ctx, tx, residentData, contractModel.ID); err != nil {
														fileError = append(fileError, fmt.Errorf("row %d, resident %d: %v", row-5, residentRowNo-5, err))
														subProcessSuccess = false
														processSuccess = false
														tx.RollbackTo(fmt.Sprintf("sp_contract_%d_resident_%d", row-5, residentRowNo-5))
													}
												}

												if subProcessSuccess {
													f.SetCellValue(contractNo, fmt.Sprintf("N%d", residentRowNo), "✔")
													f.SetCellStyle(contractNo, fmt.Sprintf("N%d", residentRowNo), fmt.Sprintf("N%d", residentRowNo), successCell)
												} else {
													f.SetCellValue(contractNo, fmt.Sprintf("N%d", residentRowNo), "✘")
													f.SetCellStyle(contractNo, fmt.Sprintf("N%d", residentRowNo), fmt.Sprintf("N%d", residentRowNo), failCell)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if !processSuccess {
			f.SetCellValue(listSheet, fmt.Sprintf("L%d", row), "✘")
			f.SetCellStyle(listSheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), failCell)
		} else {
			f.SetCellValue(listSheet, fmt.Sprintf("L%d", row), "✔")
			f.SetCellStyle(listSheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), successCell)
		}
	}

	return fileError
}

func (s *UploadService) ProcessAddBill(f *excelize.File, tx *gorm.DB, upload *models.UploadFileModel) []error {
	failCell, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#FF0000", // Red
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1}, // Style 1 = thin
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	successCell, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#00FF00", // Green
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1}, // Style 1 = thin
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})

	listSheet := f.GetSheetName(0)
	if listSheet == "" {
		return []error{errors.New("first sheet not found")}
	}

	rows, err := f.GetRows(listSheet)
	if err != nil {
		return []error{err}
	}

	if len(rows) < 6 {
		return []error{errors.New("no data rows found")}
	}
	ctx := &gin.Context{}
	ctx.Set("userID", upload.CreatorID)

	buildingInfoRow := rows[0]
	// buildingName := buildingInfoRow[1]
	buildingIDStr := buildingInfoRow[3]
	buildingID, err := strconv.ParseInt(buildingIDStr, 10, 64)
	if err != nil {
		return []error{err}
	}

	if !s.CheckUserPermissionForBuilding(ctx, buildingID) {
		return []error{fmt.Errorf("user: %d does not have permission for building ID: %d", upload.CreatorID, buildingID)}
	}
	paymentPeriodRow := rows[1]
	paymentYearStr := paymentPeriodRow[1]
	paymentMonthStr := paymentPeriodRow[3]
	paymentPeriodStr := paymentYearStr + "-" + fmt.Sprintf("%02s", paymentMonthStr)
	paymentPeriod, err := time.Parse("2006-01-02", paymentPeriodStr+"-01")
	if err != nil {
		return []error{errors.New("invalid payment period")}
	}

	dataRows := rows[5:]
	fileError := []error{}

	for index := range dataRows {
		processSuccess := true
		row := index + 6

		// // Get room floor
		// floorStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("B%d", row))
		// floor, err := strconv.ParseInt(strings.TrimSpace(floorStr), 10, 64)
		// if err != nil {
		// 	fileError = append(fileError, fmt.Errorf("row %d: failed to parse room floor: %v", row-5, err))
		// 	processSuccess = false
		// }

		// Get room no
		roomNoStr, _ := f.GetCellValue(listSheet, fmt.Sprintf("C%d", row))
		roomNo, err := strconv.ParseInt(strings.TrimSpace(roomNoStr), 10, 64)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to parse room no: %v", row-5, err))
			processSuccess = false
		} else {
			roomModel := &models.RoomModel{}
			if err := s.roomRepository.GetRoomInfoForUpload(roomModel, buildingID, int(roomNo)); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: failed to get room info: %v", row-5, err))
				processSuccess = false
			} else if roomModel.ID == 0 {
				fileError = append(fileError, fmt.Errorf("row %d: failed to get info of room no: %d", row-5, roomNo))
				processSuccess = false
			} else if len(roomModel.Contracts) == 0 {
				fileError = append(fileError, fmt.Errorf("row %d: no active contract in this room", row-5))
				processSuccess = false
			} else {
				billNo, _ := f.GetCellValue(listSheet, fmt.Sprintf("A%d", row))
				billNo = strings.TrimSpace(billNo)

				paymentDetailSheet, err := f.GetRows(billNo)
				if err != nil {
					fileError = append(fileError, fmt.Errorf("row %d: failed to get bill's payment detail sheet: %v", row-5, err))
					processSuccess = false
				} else if len(paymentDetailSheet) < 6 {
					fileError = append(fileError, fmt.Errorf("row %d: bill's payment detail sheet does not contain any data", row-5))
					processSuccess = false
				}

				billTitle, _ := f.GetCellValue(listSheet, fmt.Sprintf("D%d", row))
				billTitle = strings.TrimSpace(billTitle)
				billNote, _ := f.GetCellValue(listSheet, fmt.Sprintf("H%d", row))
				billNote = strings.TrimSpace(billNote)
				paymentTime, _ := f.GetCellValue(listSheet, fmt.Sprintf("G%d", row))
				paymentTime = strings.TrimSpace(paymentTime)
				payerCustomerNo, _ := f.GetCellValue(listSheet, fmt.Sprintf("E%d", row))
				payerCustomerNo = strings.TrimSpace(payerCustomerNo)

				customerModel := &models.UserModel{}
				if payerCustomerNo != "" && paymentTime != "" {
					if err := s.userRepository.GetCustomerByUserNo(nil, customerModel, payerCustomerNo); err != nil {
						fileError = append(fileError, fmt.Errorf("row %d: failed to get customer info: %v", row-5, err))
						processSuccess = false
					} else if customerModel.ID == 0 { //|| !utils.CompareStringRaw(strings.ReplaceAll(utils.GetUserFullName(customerModel), " ", ""), customerName) {
						fileError = append(fileError, fmt.Errorf("row %d: customer of number %s not found", row-5, payerCustomerNo))
						processSuccess = false
					}

					residentList := &[]models.RoomResidentModel{}
					if err := s.contractRepository.GetContractResidents(nil, roomModel.Contracts[0].ID, residentList); err != nil {
						fileError = append(fileError, fmt.Errorf("row %d: failed to get contract residents: %v", row-5, err))
						processSuccess = false
					} else {
						isPayerBelongToContract := false

						if roomModel.Contracts[0].HouseholderID == customerModel.ID {
							isPayerBelongToContract = true
						}

						for _, resident := range *residentList {
							if resident.UserAccountID.Valid && resident.UserAccountID.Int64 == customerModel.ID {
								isPayerBelongToContract = true
								break
							}
						}

						if !isPayerBelongToContract {
							fileError = append(fileError, fmt.Errorf("row %d: payer does not belong to this contract", row-5))
							processSuccess = false
						}
					}
				} else if (payerCustomerNo == "" && paymentTime != "") || (payerCustomerNo != "" && paymentTime == "") {
					fileError = append(fileError, fmt.Errorf("row %d: payer customer no and payment time must be both filled or both empty", row-5))
					processSuccess = false
				}

				status := 0
				if payerCustomerNo != "" && paymentTime != "" {
					status = constants.Common.BillStatus.PAID
				} else {
					lastDayOfPaymentMonth := utils.GetLastDayOfMonth(paymentPeriodStr)

					result, _ := utils.CompareDates(lastDayOfPaymentMonth, time.Now().Format("2006-01-02"))

					switch result {
					case 0, 1:
						status = constants.Common.BillStatus.UN_PAID
					case -1:
						status = constants.Common.BillStatus.OVERDUE
					}
				}

				if processSuccess {
					billStruct := &structs.UploadBill{
						Title:       billTitle,
						Period:      paymentPeriodStr,
						Status:      status,
						Note:        billNote,
						PayerID:     customerModel.ID,
						PaymentTime: paymentTime,
						ContractID:  roomModel.Contracts[0].ID,
					}

					if err := constants.Validate.Struct(billStruct); err != nil {
						fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
						processSuccess = false
					}
				}

				if processSuccess {
					payTime, _ := utils.ParseTimeWithZone(paymentTime + " 00:00:00")

					billModel := &models.BillModel{
						Title:       billTitle,
						Period:      paymentPeriod,
						Status:      status,
						Note:        sql.NullString{String: billNote, Valid: billNote != ""},
						PayerID:     sql.NullInt64{Int64: customerModel.ID, Valid: customerModel.ID != 0},
						PaymentTime: sql.NullTime{Time: payTime, Valid: paymentTime != ""},
						ContractID:  roomModel.Contracts[0].ID,
						Amount:      0.0,
					}

					tx.SavePoint(fmt.Sprintf("sp_bill_%d", row-5))
					if err := s.billRepository.CreateBill(ctx, tx, billModel); err != nil {
						fileError = append(fileError, fmt.Errorf("row %d: %v", row-5, err))
						processSuccess = false
						tx.RollbackTo(fmt.Sprintf("sp_bill_%d", row-5))
					} else {
						paymentRows := paymentDetailSheet[5:]
						totalAmount := 0.0
						for index2 := range paymentRows {
							paymentRowNo := index2 + 6
							subProcessSuccess := true

							paymentName, _ := f.GetCellValue(billNo, fmt.Sprintf("B%d", paymentRowNo))
							paymentName = strings.TrimSpace(paymentName)

							paymentAmountStr, _ := f.GetCellValue(billNo, fmt.Sprintf("C%d", paymentRowNo))
							paymentAmount, err := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(paymentAmountStr, ",", "")), 64)
							if err != nil {
								fileError = append(fileError, fmt.Errorf("row %d, payment detail %d: failed to parse payment amount: %v", row-5, paymentRowNo-5, err))
								processSuccess = false
								subProcessSuccess = false
							}
							paymentNote, _ := f.GetCellValue(billNo, fmt.Sprintf("D%d", paymentRowNo))
							paymentNote = strings.TrimSpace(paymentNote)

							if subProcessSuccess {
								paymentStruct := &structs.NewPayment{
									Name:   paymentName,
									Amount: paymentAmount,
									Note:   paymentNote,
								}

								if err := constants.Validate.Struct(paymentStruct); err != nil {
									fileError = append(fileError, fmt.Errorf("row %d, payment detail %d: %v", row-5, paymentRowNo-5, err))
									processSuccess = false
									subProcessSuccess = false
								}
							}

							if subProcessSuccess {
								totalAmount += paymentAmount

								paymentModel := &models.BillPaymentModel{
									Name:   paymentName,
									Amount: paymentAmount,
									Note:   sql.NullString{String: paymentNote, Valid: paymentNote != ""},
									BillID: billModel.ID,
								}

								tx.SavePoint(fmt.Sprintf("sp_bill_%d_payment_%d", row-5, paymentRowNo-5))
								if err := s.billRepository.AddNewPayment2(ctx, tx, paymentModel); err != nil {
									fileError = append(fileError, fmt.Errorf("row %d, payment detail %d: %v", row-5, paymentRowNo-5, err))
									processSuccess = false
									subProcessSuccess = false
									tx.RollbackTo(fmt.Sprintf("sp_bill_%d_payment_%d", row-5, paymentRowNo-5))
								}
							}

							if subProcessSuccess {
								f.SetCellValue(billNo, fmt.Sprintf("E%d", paymentRowNo), "✔")
								f.SetCellStyle(billNo, fmt.Sprintf("E%d", paymentRowNo), fmt.Sprintf("E%d", paymentRowNo), successCell)
							} else {
								f.SetCellValue(billNo, fmt.Sprintf("E%d", paymentRowNo), "✘")
								f.SetCellStyle(billNo, fmt.Sprintf("E%d", paymentRowNo), fmt.Sprintf("E%d", paymentRowNo), failCell)
							}
						}

						if processSuccess {
							billModel.Amount = totalAmount
							if err := s.billRepository.UpdateBill(ctx, tx, billModel, billModel.ID); err != nil {
								fileError = append(fileError, fmt.Errorf("row %d: failed to update bill total amount: %v", row-5, err))
								processSuccess = false
							}
						}
					}
				}
			}
		}

		if processSuccess {
			f.SetCellValue(listSheet, fmt.Sprintf("I%d", row), "✔")
			f.SetCellStyle(listSheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), successCell)
		} else {
			f.SetCellValue(listSheet, fmt.Sprintf("I%d", row), "✘")
			f.SetCellStyle(listSheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), failCell)
		}
	}

	return fileError
}

func (s *UploadService) ProcessUploadFile(upload *models.UploadFileModel) error {
	filePath := strings.ReplaceAll(upload.URLPath, "/api/", "")
	fileDecompositions := strings.Split(upload.URLPath, "/")
	fileName := fileDecompositions[len(fileDecompositions)-1]

	// Read the file content
	bytes, err := utils.ReadFile(filePath)
	if err != nil {
		return err
	}

	// From the file content create a .xlsx/.xls file for go-exelize to read on
	currentDate := time.Now().Format("2006-01-02")
	currentYear := time.Now().Format("2006")
	currentMonth := time.Now().Format("2006-01")
	filePath = filepath.Join("assets", "cron", currentYear, currentMonth, currentDate, filePath)
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	logFile, err := os.Create(filepath.Join(filepath.Dir(filePath), "result.log"))
	if err != nil {
		return err
	}
	defer logFile.Close()

	if _, err := file.Write(bytes); err != nil {
		return err
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		fmt.Fprintf(logFile, "failed to start transaction for file %s: %v\n", fileName, tx.Error)
		return tx.Error
	}

	// Use go-excelize to read the file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		s.CronFileFail(tx, upload)
		tx.Commit()
		fmt.Fprintf(logFile, "failed to open file %s: %v\n", fileName, err)
		return err
	}
	defer f.Close()

	var fileError []error = []error{}

	tx.SavePoint("sp_start")
	switch upload.UploadType {
	case constants.Common.UploadType.ADD_CUSTOMERS:
		fileError = s.ProcessAddCustomer(f, tx, upload)
	case constants.Common.UploadType.ADD_CONTRACTS:
		fileError = s.ProcessAddContract(f, tx, upload)
	case constants.Common.UploadType.ADD_BILLS:
		fileError = s.ProcessAddBill(f, tx, upload)
	}
	f.Save()

	if len(fileError) > 0 {
		tx.RollbackTo("sp_start")
		for _, err := range fileError {
			// Write the result to logFile
			fmt.Fprintf(logFile, "%v\n", err)
		}
		fmt.Fprintf(logFile, "File %s failed to process with %d errors\n", fileName, len(fileError))

		if err := utils.OverWriteFile(strings.ReplaceAll(upload.URLPath, "/api/", ""), filePath); err != nil {
			fmt.Fprintf(logFile, "failed to overwrite file %s result: %v\n", fileName, err)
			tx.Rollback()
			return err
		}

		if err := s.CronFileFail(tx, upload); err != nil {
			fmt.Fprintf(logFile, "failed to update upload file %s: %v\n", fileName, err)
			tx.Rollback()
			return err
		}

		tx.Commit()
		return fmt.Errorf("file %s failed to process with %d errors", fileName, len(fileError))
	}

	if err := utils.OverWriteFile(strings.ReplaceAll(upload.URLPath, "/api/", ""), filePath); err != nil {
		fmt.Fprintf(logFile, "failed to overwrite file %s result: %v\n", fileName, err)
		tx.Rollback()
		return err
	}

	if err := s.CronFileSuccess(tx, upload); err != nil {
		fmt.Fprintf(logFile, "failed to update upload file %s: %v\n", fileName, err)
		tx.Rollback()
		return err
	}

	fmt.Fprintf(logFile, "File %s processed successfully\n", fileName)

	defer func() {
		if r := recover(); r != nil {
			fileDecompositions := strings.Split(upload.URLPath, "/")
			fileName := fileDecompositions[len(fileDecompositions)-1]
			fmt.Printf("Recovered from panic while processing file %s: %v\n", fileName, r)
			if err := s.CronFileFail(tx, upload); err != nil {
				fmt.Fprintf(logFile, "failed to update upload file %s: %v\n", fileName, err)
				tx.Rollback()
			}
		}
	}()

	tx.Commit()

	return nil
}

func (s *UploadService) RunUploadCron() {
	uploads := &[]models.UploadFileModel{}

	if err := s.repository.GetUploadFileForCron(uploads); err != nil {
		return
	}

	// // Use a WaitGroup to keep track of active goroutines
	// var wg sync.WaitGroup
	// wg.Add(len(*uploads))

	for _, upload := range *uploads {
		// go
		func(upload *models.UploadFileModel) {
			// defer wg.Done()
			// Process each upload file
			fileDecompositions := strings.Split(upload.URLPath, "/")
			fileName := fileDecompositions[len(fileDecompositions)-1]
			if err := s.ProcessUploadFile(upload); err == nil {
				fmt.Printf("File %s processed successfully\n", fileName)
			} else {
				fmt.Printf("Failed to process file %s\n", fileName)
			}
		}(&upload)
	}
	if len(*uploads) > 0 {
		fmt.Println("Check the log files for more details")
	}

	// wg.Wait()
}
