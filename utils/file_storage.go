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
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	s3Client    *s3.Client    = nil
	minioClient *minio.Client = nil
)

func initS3Connection() {
	awsRegion := appConfig.GetEnv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "ap-southeast-2"
	}

	if appConfig.GetEnv("AWS_BUCKET") == "" {
		return
	}

	if appConfig.GetEnv("AWS_ACCESS_KEY_ID") == "" {
		return
	}

	if appConfig.GetEnv("AWS_SECRET_ACCESS_KEY") == "" {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))

	if err != nil {
		return
	}

	localS3Client := s3.NewFromConfig(cfg)
	s3Client = localS3Client
}

func initMinioConnection() {
	endpoint := appConfig.GetEnv("MINIO_ENDPOINT")
	if endpoint == "" {
		return
	}

	accessKey := appConfig.GetEnv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		return
	}

	secretKey := appConfig.GetEnv("MINIO_SECRET_KEY")
	if secretKey == "" {
		return
	}

	useSSL := appConfig.GetEnv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return
	}

	minioClient = client
}

func InitStorageServices() {
	initS3Connection()
	initMinioConnection()
}

func generateUniqueFileName() string {
	// Create a unique file name
	now := time.Now()
	milliseconds := now.Nanosecond() / 1000000
	formattedTime := now.Format("2006-01-02_15-04-05") + "-" + strconv.Itoa(milliseconds)
	fileName := formattedTime

	return fileName
}

func fileExistsInS3(path string) bool {
	_, err := s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(path),
	})

	return err == nil
}

func fileExistsInMinio(path string) bool {
	_, err := minioClient.StatObject(context.TODO(), appConfig.GetEnv("MINIO_BUCKET"), path, minio.StatObjectOptions{})
	return err == nil
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

func saveFileToMinio(file *multipart.FileHeader, filePath string) error {
	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()

	_, err = minioClient.PutObject(context.TODO(), appConfig.GetEnv("MINIO_BUCKET"), filePath, fileContent, file.Size, minio.PutObjectOptions{})
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
	if folder == "" {
		return "", errors.New("folder cannot be empty")
	}

	fileName := generateUniqueFileName() + filepath.Ext(file.Filename)
	if folder[len(folder)-1:] != "/" {
		folder = folder + "/"
	}

	filePath := folder + fileName

	var err1, err2, err3 error

	if s3Client != nil {
		err1 = saveFileToS3(file, filePath)
	}

	if minioClient != nil {
		err2 = saveFileToMinio(file, filePath)
	}

	err3 = saveFileToLocal(file, filePath)
	if err1 == nil || err2 == nil || err3 == nil {
		return "/api/" + filePath, nil
	}

	return "", errors.New("failed to save file to all storage services")
}

func StoreFileAllMedia(file *multipart.FileHeader, folder string) (string, error) {
	if folder == "" {
		return "", errors.New("folder cannot be empty")
	}

	fileName := generateUniqueFileName() + filepath.Ext(file.Filename)
	if folder[len(folder)-1:] != "/" {
		folder = folder + "/"
	}

	filePath := folder + fileName

	var err1, err2, err3 error

	if s3Client != nil {
		err1 = saveFileToS3(file, filePath)
	} else {
		err1 = errors.New("can not save file to S3")
	}

	if minioClient != nil {
		err2 = saveFileToMinio(file, filePath)
	} else {
		err2 = errors.New("can not save file to MinIO")
	}

	err3 = saveFileToLocal(file, filePath)

	if err1 == nil || err2 == nil || err3 == nil {
		return "/api/" + filePath, nil
	}

	return "", errors.New("failed to save file to any storage service")
}

func removeFileFromS3(filePath string) error {
	if filePath == "" {
		return nil
	}

	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return err
	}
	return nil
}

func removeFileFromMinio(filePath string) error {
	if filePath == "" {
		return nil
	}

	err := minioClient.RemoveObject(context.TODO(), appConfig.GetEnv("MINIO_BUCKET"), filePath, minio.RemoveObjectOptions{})
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

	if s3Client != nil && fileExistsInS3(filePath) {
		removeFileFromS3(filePath)
	}

	if minioClient != nil && fileExistsInMinio(filePath) {
		removeFileFromMinio(filePath)
	}

	removeFileFromLocal(filePath)
}

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

	// ctx.Header("Content-Type", "image/"+strings.TrimPrefix(filepath.Ext(filepath.Base(filePath)), "."))
	ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	ctx.Header("Content-Length", strconv.FormatInt(contentLength, 10))

	if _, err := io.Copy(ctx.Writer, result.Body); err != nil {
		return err
	}

	return nil
}

func getFileFromMinio(ctx *gin.Context, filePath string) error {
	object, err := minioClient.GetObject(context.TODO(), appConfig.GetEnv("MINIO_BUCKET"), filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	fileInfo, err := object.Stat()
	if err != nil {
		return err
	}

	contentLength := fileInfo.Size

	if contentLength == 0 {
		return errors.New("file is empty")
	}

	// ctx.Header("Content-Type", "image/"+strings.TrimPrefix(filepath.Ext(filepath.Base(filePath)), "."))
	ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	ctx.Header("Content-Length", strconv.FormatInt(contentLength, 10))

	if _, err := io.Copy(ctx.Writer, object); err != nil {
		return err
	}

	return nil
}

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

	// ctx.Header("Content-Type", "image/"+strings.TrimPrefix(filepath.Ext(filepath.Base(filePath)), "."))
	ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	ctx.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	if _, err := io.Copy(ctx.Writer, fileContent); err != nil {
		return err
	}

	return nil
}

func GetFile(ctx *gin.Context, filePath string) error {
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	var err error

	if s3Client != nil && fileExistsInS3(filePath) {
		err = getFileFromS3(ctx, filePath)
		if err == nil {
			return nil
		}
	}

	if minioClient != nil && fileExistsInMinio(filePath) {
		err = getFileFromMinio(ctx, filePath)
		if err == nil {
			return nil
		}
	}

	return getFileFromLocal(ctx, filePath)
}

func readFileFromS3(filePath string) ([]byte, error) {
	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func readFileFromMinio(filePath string) ([]byte, error) {
	object, err := minioClient.GetObject(context.TODO(), appConfig.GetEnv("MINIO_BUCKET"), filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return io.ReadAll(object)
}

func readFileFromLocal(filePath string) ([]byte, error) {
	filePath = filepath.Join("assets", "files", filePath)
	fileContent, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	return io.ReadAll(fileContent)
}

func ReadFile(filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, errors.New("file path cannot be empty")
	}

	if s3Client != nil && fileExistsInS3(filePath) {
		return readFileFromS3(filePath)
	}

	if minioClient != nil && fileExistsInMinio(filePath) {
		return readFileFromMinio(filePath)
	}

	return readFileFromLocal(filePath)
}

func ovewriteFileToS3(targetPath string, sourcePath string) error {
	fileContent, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(appConfig.GetEnv("AWS_BUCKET")),
		Key:    aws.String(targetPath),
		Body:   fileContent,
	})
	if err != nil {
		return err
	}

	return nil
}

func ovewriteFileToMinio(targetPath string, sourcePath string) error {
	fileContent, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	fileStat, err := fileContent.Stat()
	if err != nil {
		return err
	}

	_, err = minioClient.PutObject(context.TODO(), appConfig.GetEnv("MINIO_BUCKET"), targetPath, fileContent, fileStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ovewriteFileToLocal(targetPath string, sourcePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func OverWriteFile(targetPath string, sourcePath string) error {
	if targetPath == "" {
		return errors.New("target path cannot be empty")
	}

	if sourcePath == "" {
		return errors.New("source path cannot be empty")
	}

	var (
		s3Error    error = nil
		minioError error = nil
		localError error = nil
	)

	if s3Client != nil && fileExistsInS3(targetPath) {
		s3Error = ovewriteFileToS3(targetPath, sourcePath)
	}

	if minioClient != nil && fileExistsInMinio(targetPath) {
		minioError = ovewriteFileToMinio(targetPath, sourcePath)
	}

	localError = ovewriteFileToLocal(filepath.Join("assets", "files", targetPath), sourcePath)

	if s3Error != nil && minioError != nil && localError != nil {
		return errors.New("failed to overwrite file in all storage services")
	}

	return nil
}
