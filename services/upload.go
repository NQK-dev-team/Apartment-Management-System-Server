package services

import (
	"api/config"
	"api/constants"
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadService struct {
	repository *repositories.UploadRepository
}

func NewUploadService() *UploadService {
	return &UploadService{
		repository: repositories.NewUploadRepository(),
	}
}

func (s *UploadService) UploadFile(ctx *gin.Context, upload *structs.UploadStruct) error {
	jwt := ctx.GetString("jwt")

	token, err := utils.ValidateJWTToken(jwt)

	if err != nil {
		return err
	}

	claim := &structs.JTWClaim{}

	utils.ExtractJWTClaim(token, claim)

	return config.DB.Transaction(func(tx *gorm.DB) error {
		uploadModel := &models.UploadFileModel{
			CreatorID:   claim.UserID,
			FileName:    upload.File.Filename,
			URLPath:     "",
			StoragePath: "",
			Size:        upload.File.Size,
			UploadType:  upload.UploadType,
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
		uploadModel.StoragePath = strings.ReplaceAll(filePath, "/api/", "")

		if err := s.repository.Update(ctx, tx, uploadModel); err != nil {
			return err
		}

		return nil
	})
}

func (s *UploadService) GetUploads(ctx *gin.Context, uploads *[]models.UploadFileModel, uploadType int, isProcessed bool) error {
	return s.repository.Get(ctx, uploads, uploadType, isProcessed)
}

func (s *UploadService) GetUploadByID(ctx *gin.Context, upload *models.UploadFileModel, uploadID int64) error {
	return s.repository.GetByID(ctx, upload, uploadID)
}
