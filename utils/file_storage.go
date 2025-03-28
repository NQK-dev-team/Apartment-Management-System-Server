package utils

import (
	"context"
	"errors"
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
	"github.com/gin-gonic/gin"
)

var s3Error error
var s3Client *s3.Client

func InitS3Connection() {
	awsRegion := appConfig.GetEnv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "ap-southeast-2"
	}

	if appConfig.GetEnv("AWS_BUCKET") == "" {
		s3Error = errors.New("AWS_BUCKET is not set")
		return
	}

	if appConfig.GetEnv("AWS_ACCESS_KEY_ID") == "" {
		s3Error = errors.New("AWS_ACCESS_KEY_ID is not set")
		return
	}

	if appConfig.GetEnv("AWS_SECRET_ACCESS_KEY") == "" {
		s3Error = errors.New("AWS_SECRET_ACCESS_KEY is not set")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))

	if err != nil {
		s3Error = err
		return
	}

	localS3Client := s3.NewFromConfig(cfg)
	// _, err = localS3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{})

	// if err != nil {
	// 	s3Error = err
	// 	return
	// }

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

func fileExistsInS3(path string) bool {
	_, err := s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(path),
	})

	if err != nil {
		return false
	}
	return true
}

func saveFileToS3(file *multipart.FileHeader, filePath string) error {
	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()

	// uploader := manager.NewUploader(s3Client)
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

func StoreFile(file *multipart.FileHeader, folder string) (string, error) {
	fileName := generateUniqueFileName() + filepath.Ext(file.Filename)
	if folder[len(folder)-1:] != "/" {
		folder = folder + "/"
	}

	filePath := folder + fileName

	if checkS3Connection() {
		if err := saveFileToS3(file, filePath); err != nil {
			if err := saveFileToLocal(file, filePath); err != nil {
				return "", err
			}
		}
		return "/api/" + filePath, nil
	}
	if err := saveFileToLocal(file, filePath); err != nil {
		return "", err
	}
	return "/api/" + filePath, nil
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

func RemoveFile(filePath string) {
	filePath = strings.Replace(filePath, "/api/", "", -1)

	if checkS3Connection() && fileExistsInS3(filePath) {
		if err := removeFileFromS3(filePath); err != nil {
			removeFileFromLocal(filePath)
		}
		return
	}
	removeFileFromLocal(filePath)
}

// func getFileFromS3(file *structs.CustomFileStruct, filePath string) error {
// 	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
// 		Key:    aws.String(filePath),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	defer result.Body.Close()

// 	// Read file content into buffer
// 	buf := new(bytes.Buffer)
// 	if _, err := io.Copy(buf, result.Body); err != nil {
// 		return err
// 	}

// 	// Extract file metadata
// 	contentLength := aws.ToInt64(result.ContentLength)

// 	if contentLength == 0 {
// 		return errors.New("file is empty")
// 	}

// 	// Populate the file struct
// 	file.Filename = filepath.Base(filePath)
// 	file.Content = buf.Bytes()
// 	file.Size = contentLength
// 	file.Header = make(map[string][]string)
// 	file.Header.Set("Content-Type", "application/octet-stream")

// 	return nil
// }

func getFileFromS3(ctx *gin.Context, filePath string) error {
	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()

	// // Read file content into buffer
	// buf := new(bytes.Buffer)
	// 	if _, err := io.Copy(buf, result.Body); err != nil {
	// 		return err
	// 	}

	// Extract file metadata
	contentLength := aws.ToInt64(result.ContentLength)

	if contentLength == 0 {
		return errors.New("file is empty")
	}

	// // Populate the file struct
	// file.Filename = filepath.Base(filePath)
	// file.Content = buf.Bytes()
	// file.Size = contentLength
	// file.Header = make(map[string][]string)
	// file.Header.Set("Content-Type", "application/octet-stream")

	ctx.Header("Content-Type", "image/"+strings.TrimPrefix(filepath.Ext(filepath.Base(filePath)), "."))
	ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	ctx.Header("Content-Length", strconv.FormatInt(contentLength, 10))

	if _, err := io.Copy(ctx.Writer, result.Body); err != nil {
		return err
	}

	return nil
}

// func getFileFromLocal(file *structs.CustomFileStruct, filePath string) error {
// 	filePath = filepath.Join("assets", "files", filePath)
// 	fileContent, err := os.Open(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer fileContent.Close()

// 	fileInfo, err := fileContent.Stat()
// 	if err != nil {
// 		return err
// 	}

// 	fileBytes, err := io.ReadAll(fileContent)

// 	if err != nil {
// 		return err
// 	}

// 	// *file = *fileHeader
// 	file.Filename = filepath.Base(filePath)
// 	file.Size = fileInfo.Size()
// 	file.Header = make(map[string][]string)
// 	file.Content = fileBytes
// 	file.Header.Set("Content-Type", "application/octet-stream")

// 	return nil
// }

func getFileFromLocal(ctx *gin.Context, filePath string) error {
	filePath = filepath.Join("assets", "files", filePath)
	fileContent, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	fileInfo, err := fileContent.Stat()
	if err != nil {
		return err
	}

	// fileBytes, err := io.ReadAll(fileContent)

	// if err != nil {
	// 	return err
	// }

	// // *file = *fileHeader
	// file.Filename = filepath.Base(filePath)
	// file.Size = fileInfo.Size()
	// file.Header = make(map[string][]string)
	// file.Content = fileBytes
	// file.Header.Set("Content-Type", "application/octet-stream")

	ctx.Header("Content-Type", "image/"+strings.TrimPrefix(filepath.Ext(filepath.Base(filePath)), "."))
	ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	ctx.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	if _, err := io.Copy(ctx.Writer, fileContent); err != nil {
		return err
	}

	return nil
}

// func GetFile(file *structs.CustomFileStruct, filePath string) error {
// 	if checkS3Connection() {
// 		if err := getFileFromS3(file, filePath); err != nil {
// 			return getFileFromLocal(file, filePath)
// 		}
// 		return nil
// 	}

// 	return getFileFromLocal(file, filePath)
// }

func GetFile(ctx *gin.Context, filePath string) error {
	if checkS3Connection() && fileExistsInS3(filePath) {
		if err := getFileFromS3(ctx, filePath); err != nil {
			return getFileFromLocal(ctx, filePath)
		}
		return nil
	}

	return getFileFromLocal(ctx, filePath)
}
