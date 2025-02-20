package utils

import (
	"context"
	"mime/multipart"
	"strconv"
	"time"

	appConfig "api/config"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Error error
var s3Client *s3.Client

func InitS3Connection() {
	awsRegion, err := appConfig.GetEnv("AWS_REGION")
	if err != nil {
		s3Error = err
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))

	if err != nil {
		s3Error = err
		return
	}

	localS3Client := s3.NewFromConfig(cfg)
	_, err = localS3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{})

	if err != nil {
		s3Error = err
		return
	}

	s3Client = localS3Client
}

func generateUniqueFileName() string {
	// Create a unique file name
	now := time.Now()
	milliseconds := now.Nanosecond() / 1000000
	formattedTime := now.Format("2006-01-02_15-04-05") + "-" + strconv.Itoa(milliseconds)
	fileName := formattedTime

	return fileName
}

func checkS3Connection() bool {
	return s3Error == nil
}

func saveFileToS3(file *multipart.FileHeader, folder string) error {
	fileName := generateUniqueFileName()
	fileName = folder + "/" + fileName
	return nil
}

func saveFileToLocal(file *multipart.FileHeader, folder string) error {
	fileName := generateUniqueFileName()
	fileName = folder + "/" + fileName
	return nil
}

func removeFileFromS3(fileName string, folder string) error {
	return nil
}

func removeFileFromLocal(fileName string, folder string) error {
	return nil
}

func StoreFile(file *multipart.FileHeader, folder string) error {
	if checkS3Connection() {
		return saveFileToS3(file, folder)
	}
	return saveFileToLocal(file, folder)
}

func RemoveFile(fileName string, folder string) error {
	if checkS3Connection() {
		return removeFileFromS3(fileName, folder)
	}
	return removeFileFromLocal(fileName, folder)
}
