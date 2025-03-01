package utils

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	appConfig "api/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Error error
var s3Client *s3.Client

func InitS3Connection() {
	awsRegion := appConfig.GetEnv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "ap-southeast-2"
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

func saveFileToS3(file *multipart.FileHeader, filePath string) error {
	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(filePath),
		Body:   fileContent,
		// ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		return err
	}

	return nil
}

func saveFileToLocal(file *multipart.FileHeader, filePath string) error {
	filePath = filepath.Join("assets", "files", filePath)
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, fileContent)
	if err != nil {
		return err
	}

	return nil
}

func removeFileFromS3(filePath string) error {
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return err
	}
	return nil
}

func removeFileFromLocal(filePath string) error {
	filePath = filepath.Join("assets", "files", filePath)
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func StoreFile(file *multipart.FileHeader, folder string) (string, error) {
	fileName := generateUniqueFileName() + filepath.Ext(file.Filename)
	filePath := folder + "/" + fileName

	if checkS3Connection() {
		if err := saveFileToS3(file, filePath); err != nil {
			if err := saveFileToLocal(file, filePath); err != nil {
				return "", err
			}
			return "/api/" + filePath, nil
		}
		return filePath, nil
	}
	if err := saveFileToLocal(file, filePath); err != nil {
		return "", err
	}
	return "/api/" + filePath, nil
}

func RemoveFile(filePath string) {
	filePath = strings.Replace(filePath, "/api/", "", -1)

	if checkS3Connection() {
		if err := removeFileFromS3(filePath); err != nil {
			removeFileFromLocal(filePath)
		}
		return
	}
	removeFileFromLocal(filePath)
}
