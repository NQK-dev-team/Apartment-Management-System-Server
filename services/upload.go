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
	// if tx == nil {
	// 	return config.DB.Transaction(func(tx *gorm.DB) error {
	// 		upload.ProcessResult = sql.NullInt64{
	// 			Int64: constants.Common.CronUploadProcessResult.FAILED,
	// 			Valid: true,
	// 		}

	// 		upload.ProcessDate = sql.NullTime{
	// 			Time:  time.Now(),
	// 			Valid: true,
	// 		}

	// 		return s.repository.Update(nil, tx, upload)
	// 	})
	// }

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
	// if tx == nil {
	// 	return config.DB.Transaction(func(tx *gorm.DB) error {
	// 		upload.ProcessResult = sql.NullInt64{
	// 			Int64: constants.Common.CronUploadProcessResult.SUCCESS,
	// 			Valid: true,
	// 		}

	// 		upload.ProcessDate = sql.NullTime{
	// 			Time:  time.Now(),
	// 			Valid: true,
	// 		}

	// 		return s.repository.Update(nil, tx, upload)
	// 	})
	// }

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
	listSheet := f.GetSheetName(0)
	if listSheet == "" {
		return []error{errors.New("first sheet not found"), s.CronFileFail(tx, upload)}
	}

	rows, err := f.GetRows(listSheet)
	if err != nil {
		return []error{err, s.CronFileFail(tx, upload)}
	}

	if len(rows) < 6 {
		return []error{errors.New("no data rows found"), s.CronFileFail(tx, upload)}
	}

	dataRows := rows[5:]

	fileError := []error{}

	for index, data := range dataRows {
		processSuccess := true
		row := index + 5

		// Process customer gender string
		gender := 0
		if utils.CompareStringRaw("Nam", strings.TrimSpace(data[4])) || utils.CompareStringRaw("Male", strings.TrimSpace(data[4])) {
			gender = constants.Common.UserGender.MALE
		} else if utils.CompareStringRaw("Nữ", strings.TrimSpace(data[4])) || utils.CompareStringRaw("Female", strings.TrimSpace(data[4])) {
			gender = constants.Common.UserGender.FEMALE
		} else if utils.CompareStringRaw("Khác", strings.TrimSpace(data[4])) || utils.CompareStringRaw("Other", strings.TrimSpace(data[4])) {
			gender = constants.Common.UserGender.OTHER
		}

		// Generate customer's password
		newPassword, err := utils.GeneratePassword(constants.Common.NewPasswordLength)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to generate password", row))
			newPassword = "123456"
		}
		hashedPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to hash password", row))
			hashedPassword = "$12$xG0qlWDqXflwTqTBgFRnjuA1J5zZRSd6dbzOT353TAQS7ScjVfqXW"
		}

		// Create new customer's account
		newUser := &structs.NewUploadCustomer{
			LastName:         strings.TrimSpace(data[1]),
			FirstName:        strings.TrimSpace(data[3]),
			MiddleName:       strings.TrimSpace(data[2]),
			Dob:              strings.TrimSpace(data[7]),
			Pob:              strings.TrimSpace(data[8]),
			Gender:           gender,
			SSN:              strings.TrimSpace(data[5]),
			OldSSN:           strings.TrimSpace(data[6]),
			Email:            strings.TrimSpace(data[9]),
			Phone:            strings.TrimSpace(data[10]),
			PermanentAddress: strings.TrimSpace(data[11]),
			TemporaryAddress: strings.TrimSpace(data[12]),
			ProfileImage:     "/image/placeholder_image.png",
			FrontSSNImage:    "/image/placeholder_image.png",
			BackSSNImage:     "/image/placeholder_image.png",
		}

		if err := constants.Validate.Struct(newUser); err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
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
			if err := s.userRepository.Create(ctx, tx, customerModel); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: failed to create user: %v", row, err))
				processSuccess = false
			}

			if err := s.emailService.SendAccountCreationEmail(config.GetEnv("APM_CLIENT_BASE_URL")+"/login", customerModel.Email, utils.GetUserFullName(customerModel), newPassword); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: failed to send account creation email: %v", row, err))
				processSuccess = false
			}
		}

		if !processSuccess {
			f.SetCellValue(listSheet, fmt.Sprintf("N%d", row), "FALSE")
		} else {
			f.SetCellValue(listSheet, fmt.Sprintf("N%d", row), "TRUE")
		}
	}

	if len(fileError) > 0 {
		return []error{s.CronFileFail(tx, upload)}
	}

	return fileError
}

func (s *UploadService) ProcessAddContract(f *excelize.File, tx *gorm.DB, upload *models.UploadFileModel) []error {
	listSheet := f.GetSheetName(0)
	if listSheet == "" {
		return []error{errors.New("first sheet not found"), s.CronFileFail(tx, upload)}
	}

	rows, err := f.GetRows(listSheet)
	if err != nil {
		return []error{err, s.CronFileFail(tx, upload)}
	}

	if len(rows) < 6 {
		return []error{errors.New("no data rows found"), s.CronFileFail(tx, upload)}
	}

	ctx := &gin.Context{}
	ctx.Set("userID", upload.CreatorID)

	buildingInfoRow := rows[0]
	// buildingName := strings.TrimSpace(buildingInfoRow[1])
	buildingIDStr := strings.TrimSpace(buildingInfoRow[3])
	buildingID, err := strconv.ParseInt(buildingIDStr, 10, 64)
	if err != nil {
		return []error{err, s.CronFileFail(tx, upload)}
	}

	buildingModel := &models.BuildingModel{}
	if err := s.buildingRepository.GetById(nil, buildingModel, buildingID); err != nil {
		return []error{err, s.CronFileFail(tx, upload)}
	} else if buildingModel.ID == 0 { // || !utils.CompareStringRaw(strings.ReplaceAll(buildingModel.Name, " ", ""), strings.ReplaceAll(buildingName, " ", ""))
		return []error{fmt.Errorf("building with ID %d not found", buildingID), s.CronFileFail(tx, upload)}
	}

	dataRows := rows[5:]

	fileError := []error{}

	for index, data := range dataRows {
		processSuccess := true
		row := index + 6

		// // Get room floor
		// floor, err := strconv.ParseInt(strings.TrimSpace(data[1]), 10, 64)
		// if err != nil {
		// 	fileError = append(fileError, fmt.Errorf("row %d: failed to parse room floor: %v", row, err))
		// 	processSuccess = false
		// }

		// Get room no
		roomNo, err := strconv.ParseInt(strings.TrimSpace(data[2]), 10, 64)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to parse room no: %v", row, err))
			processSuccess = false
		}

		roomModel := &models.RoomModel{}

		if err := s.roomRepository.GetRoomInfoForUpload(roomModel, buildingID, int(roomNo)); err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to get room info: %v", row, err))
			processSuccess = false
		} else if roomModel.ID == 0 {
			fileError = append(fileError, fmt.Errorf("row %d: failed to get info of room no: %d", row, roomNo))
			processSuccess = false
		} else if len(roomModel.Contracts) > 0 {
			fileError = append(fileError, fmt.Errorf("row %d: room is already contracted", row))
			processSuccess = false
		}

		// Get customer info
		customerNo := strings.TrimSpace(data[9])
		// customerName := strings.ReplaceAll(strings.TrimSpace(data[10]), " ", "")
		customerModel := &models.UserModel{}

		if err := s.userRepository.GetCustomerByUserNo(nil, customerModel, customerNo); err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to get customer info: %v", row, err))
			processSuccess = false
		} else if customerModel.ID == 0 { //|| !utils.CompareStringRaw(strings.ReplaceAll(utils.GetUserFullName(customerModel), " ", ""), customerName) {
			fileError = append(fileError, fmt.Errorf("row %d: customer of number %s not found", row, customerNo))
			processSuccess = false
		}

		if processSuccess {
			// Get contract no
			contractNo := strings.TrimSpace(data[0])

			// Get contract value
			value, err := strconv.ParseFloat(strings.TrimSpace(data[3]), 64)
			if err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: failed to parse contract value: %v", row, err))
				processSuccess = false
			}

			// Get contract type
			contractType := 0
			if utils.CompareStringRaw("Thuê", strings.TrimSpace(data[4])) || utils.CompareStringRaw("Rental", strings.TrimSpace(data[4])) {
				contractType = constants.Common.ContractType.RENT
			} else if utils.CompareStringRaw("Mua", strings.TrimSpace(data[4])) || utils.CompareStringRaw("Purchase", strings.TrimSpace(data[4])) {
				contractType = constants.Common.ContractType.BUY
			}

			newContract := &structs.NewUploadContract{
				ContractType:  contractType,
				ContractValue: value,
				HouseholderID: customerModel.ID,
				RoomID:        roomModel.ID,
				CreatorID:     upload.CreatorID,
				CreatedAt:     strings.TrimSpace(data[5]),
				StartDate:     strings.TrimSpace(data[6]),
				EndDate:       strings.TrimSpace(data[7]),
				SignDate:      strings.TrimSpace(data[8]),
			}

			if err := constants.Validate.Struct(newContract); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
				processSuccess = false
			}

			startDate, _ := time.Parse("2006-01-02", newContract.StartDate)
			createAt, _ := time.Parse("2006-01-02", newContract.CreatedAt)
			endDate := utils.StringToNullTime(newContract.EndDate)
			signDate := utils.StringToNullTime(newContract.SignDate)

			status := 0
			currentDate := time.Now().Format("2006-01-02")

			if newContract.SignDate == "" {
				result, err := utils.CompareDates(newContract.StartDate, currentDate)
				if err != nil {
					fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
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
					fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
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
								fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
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

				if err := s.contractRepository.CreateContract(ctx, tx, contractModel); err != nil {
					fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
					processSuccess = false
				} else {
					residentSheetRows, err := f.GetRows(contractNo)
					residentRows := [][]string{}
					if err == nil {
						if len(residentSheetRows) >= 6 {
							residentRows = residentSheetRows[5:]
						}
					}

					if len(residentRows) > 0 {
						for index2, residentRow := range residentRows {
							residentRowNo := index2 + 6
							subProcessSuccess := true

							residentUserModel := &models.UserModel{}
							residentNo := strings.TrimSpace(residentRow[1])
							residentLastName := strings.TrimSpace(residentRow[2])
							residentMiddleName := strings.TrimSpace(residentRow[3])
							residentFirstName := strings.TrimSpace(residentRow[4])
							residentGender := 0
							if utils.CompareStringRaw("Nam", strings.TrimSpace(data[5])) || utils.CompareStringRaw("Male", strings.TrimSpace(data[5])) {
								residentGender = constants.Common.UserGender.MALE
							} else if utils.CompareStringRaw("Nữ", strings.TrimSpace(data[5])) || utils.CompareStringRaw("Female", strings.TrimSpace(data[5])) {
								residentGender = constants.Common.UserGender.FEMALE
							} else if utils.CompareStringRaw("Khác", strings.TrimSpace(data[5])) || utils.CompareStringRaw("Other", strings.TrimSpace(data[5])) {
								residentGender = constants.Common.UserGender.OTHER
							}
							residentSSN := strings.TrimSpace(residentRow[6])
							residentOldSSN := strings.TrimSpace(residentRow[7])
							residentDOB := strings.TrimSpace(residentRow[8])
							residentPOB := strings.TrimSpace(residentRow[9])
							residentEmail := strings.TrimSpace(residentRow[10])
							residentPhone := strings.TrimSpace(residentRow[11])
							residentRelation := 0
							if utils.CompareStringRaw("Con cái", strings.TrimSpace(data[12])) || utils.CompareStringRaw("Child", strings.TrimSpace(data[12])) {
								residentRelation = constants.Common.ResidentRelationship.CHILD
							} else if utils.CompareStringRaw("Vợ/Chồng", strings.TrimSpace(data[12])) || utils.CompareStringRaw("Spouse", strings.TrimSpace(data[12])) {
								residentRelation = constants.Common.ResidentRelationship.SPOUSE
							} else if utils.CompareStringRaw("Cha/Mẹ", strings.TrimSpace(data[12])) || utils.CompareStringRaw("Parent", strings.TrimSpace(data[12])) {
								residentRelation = constants.Common.ResidentRelationship.PARENT
							} else if utils.CompareStringRaw("Khác", strings.TrimSpace(data[12])) || utils.CompareStringRaw("Other", strings.TrimSpace(data[12])) {
								residentRelation = constants.Common.ResidentRelationship.OTHER
							}

							if residentNo != "" {
								if err := s.userRepository.GetCustomerByUserNo(nil, residentUserModel, residentNo); err != nil {
									fileError = append(fileError, fmt.Errorf("row %d, resident %d: failed to get resident's customer account info: %v", row, residentRowNo, err))
									processSuccess = false
									subProcessSuccess = false
								} else {
									if residentUserModel.ID == 0 {
										fileError = append(fileError, fmt.Errorf("row %d, resident %d: customer of number %s not found", row, residentRowNo, residentNo))
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
									fileError = append(fileError, fmt.Errorf("row %d, resident %d: %v", row, residentRowNo, err))
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

								if err := s.contractRepository.AddNewRoomResident(ctx, tx, residentData, contractModel.ID); err != nil {
									fileError = append(fileError, fmt.Errorf("row %d, resident %d: %v", row, residentRowNo, err))
									subProcessSuccess = false
									processSuccess = false
								}
							}

							if subProcessSuccess {
								f.SetCellValue(contractNo, fmt.Sprintf("N%d", residentRowNo), "TRUE")
							} else {
								f.SetCellValue(contractNo, fmt.Sprintf("N%d", residentRowNo), "FALSE")
							}
						}
					}
				}
			}
		}

		if !processSuccess {
			f.SetCellValue(listSheet, fmt.Sprintf("L%d", row), "FALSE")
		} else {
			f.SetCellValue(listSheet, fmt.Sprintf("L%d", row), "TRUE")
		}
	}

	if len(fileError) > 0 {
		return []error{s.CronFileFail(tx, upload)}
	}

	return fileError
}

func (s *UploadService) ProcessAddBill(f *excelize.File, tx *gorm.DB, upload *models.UploadFileModel) []error {
	listSheet := f.GetSheetName(0)
	if listSheet == "" {
		return []error{errors.New("first sheet not found"), s.CronFileFail(tx, upload)}
	}

	rows, err := f.GetRows(listSheet)
	if err != nil {
		return []error{err, s.CronFileFail(tx, upload)}
	}

	if len(rows) < 6 {
		return []error{errors.New("no data rows found"), s.CronFileFail(tx, upload)}
	}

	ctx := &gin.Context{}
	ctx.Set("userID", upload.CreatorID)

	buildingInfoRow := rows[0]
	// buildingName := buildingInfoRow[1]
	buildingIDStr := buildingInfoRow[3]
	buildingID, err := strconv.ParseInt(buildingIDStr, 10, 64)
	if err != nil {
		return []error{err, s.CronFileFail(tx, upload)}
	}
	paymentPeriodRow := rows[1]
	paymentYearStr := paymentPeriodRow[1]
	paymentMonthStr := paymentPeriodRow[3]
	paymentPeriodStr := paymentYearStr + "-" + fmt.Sprintf("%02s", paymentMonthStr)
	paymentPeriod, err := time.Parse("2006-01-02", paymentPeriodStr+"-01")
	if err != nil {
		return []error{errors.New("invalid payment period"), s.CronFileFail(tx, upload)}
	}

	dataRows := rows[5:]
	fileError := []error{}

	for index, data := range dataRows {
		processSuccess := true
		row := index + 6

		billNo := strings.TrimSpace(data[0])
		paymentDetailSheet, err := f.GetRows(billNo)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to get bill's payment detail sheet: %v", row, err))
			processSuccess = false
		} else if len(paymentDetailSheet) < 6 {
			fileError = append(fileError, fmt.Errorf("row %d: bill's payment detail sheet does not contain any data", row))
			processSuccess = false
		}

		// // Get room floor
		// floor, err := strconv.ParseInt(strings.TrimSpace(data[1]), 10, 64)
		// if err != nil {
		// 	fileError = append(fileError, fmt.Errorf("row %d: failed to parse room floor: %v", row, err))
		// 	processSuccess = false
		// }

		// Get room no
		roomNo, err := strconv.ParseInt(strings.TrimSpace(data[2]), 10, 64)
		if err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to parse room no: %v", row, err))
			processSuccess = false
		}

		roomModel := &models.RoomModel{}
		if err := s.roomRepository.GetRoomInfoForUpload(roomModel, buildingID, int(roomNo)); err != nil {
			fileError = append(fileError, fmt.Errorf("row %d: failed to get room info: %v", row, err))
			processSuccess = false
		} else if roomModel.ID == 0 {
			fileError = append(fileError, fmt.Errorf("row %d: failed to get info of room no: %d", row, roomNo))
			processSuccess = false
		} else if len(roomModel.Contracts) == 0 {
			fileError = append(fileError, fmt.Errorf("row %d: no active contract in this room", row))
			processSuccess = false
		}

		billTitle := strings.TrimSpace(data[3])
		billNote := strings.TrimSpace(data[7])
		paymentTime := strings.TrimSpace(data[6])
		payerCustomerNo := strings.TrimSpace(data[4])

		customerModel := &models.UserModel{}
		if payerCustomerNo != "" && paymentTime != "" {
			if err := s.userRepository.GetCustomerByUserNo(nil, customerModel, payerCustomerNo); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: failed to get customer info: %v", row, err))
				processSuccess = false
			} else if customerModel.ID == 0 { //|| !utils.CompareStringRaw(strings.ReplaceAll(utils.GetUserFullName(customerModel), " ", ""), customerName) {
				fileError = append(fileError, fmt.Errorf("row %d: customer of number %s not found", row, payerCustomerNo))
				processSuccess = false
			}
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
				fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
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

			if err := s.billRepository.CreateBill(ctx, tx, billModel); err != nil {
				fileError = append(fileError, fmt.Errorf("row %d: %v", row, err))
				processSuccess = false
			} else {
				paymentRows := paymentDetailSheet[5:]
				totalAmount := 0.0
				for index2, paymentRow := range paymentRows {
					paymentRowNo := index2 + 6
					subProcessSuccess := true

					paymentName := strings.TrimSpace(paymentRow[1])
					paymentAmount, err := strconv.ParseFloat(strings.TrimSpace(paymentRow[2]), 64)
					if err != nil {
						fileError = append(fileError, fmt.Errorf("row %d, payment detail %d: failed to parse payment amount: %v", row, paymentRowNo, err))
						processSuccess = false
						subProcessSuccess = false
					}
					paymentNote := strings.TrimSpace(paymentRow[3])

					if subProcessSuccess {
						paymentStruct := &structs.NewPayment{
							Name:   paymentName,
							Amount: paymentAmount,
							Note:   paymentNote,
						}

						if err := constants.Validate.Struct(paymentStruct); err != nil {
							fileError = append(fileError, fmt.Errorf("row %d, payment detail %d: %v", row, paymentRowNo, err))
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

						if err := s.billRepository.AddNewPayment2(ctx, tx, paymentModel); err != nil {
							fileError = append(fileError, fmt.Errorf("row %d, payment detail %d: %v", row, paymentRowNo, err))
							processSuccess = false
							subProcessSuccess = false
						}
					}

					if subProcessSuccess {
						f.SetCellValue(listSheet, fmt.Sprintf("E%d", paymentRowNo), "TRUE")
					} else {
						f.SetCellValue(listSheet, fmt.Sprintf("E%d", paymentRowNo), "FALSE")
					}
				}

				if processSuccess {
					billModel.Amount = totalAmount
					if err := s.billRepository.UpdateBill(ctx, tx, billModel, billModel.ID); err != nil {
						fileError = append(fileError, fmt.Errorf("row %d: failed to update bill total amount: %v", row, err))
						processSuccess = false
					}
				}
			}
		}

		if processSuccess {
			f.SetCellValue(listSheet, fmt.Sprintf("I%d", row), "TRUE")
		} else {
			f.SetCellValue(listSheet, fmt.Sprintf("I%d", row), "FALSE")
		}
	}

	if len(fileError) > 0 {
		return []error{s.CronFileFail(tx, upload)}
	}

	return fileError
}

func (s *UploadService) ProcessUploadFile(upload *models.UploadFileModel) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
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

		// Use go-excelize to read the file
		f, err := excelize.OpenFile(filePath)
		if err != nil {
			return s.CronFileFail(tx, upload)
		}
		defer f.Close()

		var fileError []error = []error{}

		switch upload.UploadType {
		case constants.Common.UploadType.ADD_CUSTOMERS:
			fileError = s.ProcessAddCustomer(f, tx, upload)
		case constants.Common.UploadType.ADD_CONTRACTS:
			fileError = s.ProcessAddContract(f, tx, upload)
		case constants.Common.UploadType.ADD_BILLS:
			fileError = s.ProcessAddBill(f, tx, upload)
		}

		if len(fileError) > 0 {
			for _, err := range fileError {
				// fmt.Println(err)
				// Write the result to logFile
				fmt.Fprintf(logFile, "%v\n", err)
			}
			fmt.Fprintf(logFile, "File %s failed to process with %d errors\n", fileName, len(fileError))

			if err := utils.OverWriteFile(strings.ReplaceAll(upload.URLPath, "/api/", ""), filePath); err != nil {
				fmt.Fprintf(logFile, "failed to overwrite file %s result: %v\n", fileName, err)
				return fmt.Errorf("failed to overwrite file %s result: %v", fileName, err)
			}

			upload.ProcessResult = sql.NullInt64{Int64: constants.Common.CronUploadProcessResult.FAILED, Valid: true}
			upload.ProcessDate = sql.NullTime{Time: time.Now(), Valid: true}

			if err := s.repository.Update(nil, tx, upload); err != nil {
				fmt.Fprintf(logFile, "failed to update upload file %s: %v\n", fileName, err)
				return fmt.Errorf("failed to update upload file %s: %v", fileName, err)
			}

			return fmt.Errorf("file %s failed to process with %d errors", fileName, len(fileError))
		}

		if err := utils.OverWriteFile(strings.ReplaceAll(upload.URLPath, "/api/", ""), filePath); err != nil {
			fmt.Fprintf(logFile, "failed to overwrite file %s result: %v\n", fileName, err)
			return fmt.Errorf("failed to overwrite file %s result: %v", fileName, err)
		}

		upload.ProcessResult = sql.NullInt64{Int64: constants.Common.CronUploadProcessResult.SUCCESS, Valid: true}
		upload.ProcessDate = sql.NullTime{Time: time.Now(), Valid: true}

		if err := s.repository.Update(nil, tx, upload); err != nil {
			fmt.Fprintf(logFile, "failed to update upload file %s: %v\n", fileName, err)
			return fmt.Errorf("failed to update upload file %s: %v", fileName, err)
		}

		fmt.Fprintf(logFile, "File %s processed successfully\n", fileName)

		return nil
	})
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
				fmt.Printf("Failed to process file %s\nCheck the log file for more details\n", fileName)
			}
		}(&upload)
	}

	// wg.Wait()
}
