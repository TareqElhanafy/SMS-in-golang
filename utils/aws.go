package utils

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//ConnectAws function to create session to provide the configuration for SDK's service clients
func ConnectAws() *session.Session {
	godotenv.Load()
	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"", // a token will be created when the session it's used.
			),
			Endpoint: aws.String(os.Getenv("AWS_ENDPOINT")),
		})
	if err != nil {
		panic(err)
	}
	return session
}

//UploadToS3 function
func UploadToS3(ctx *gin.Context, fileName string, file multipart.File) (*s3manager.UploadOutput, error) {
	session := ctx.Value("sess").(*session.Session) //getting the session value
	uploader := s3manager.NewUploader(session)      //using S3 service client
	godotenv.Load()
	//upload to the s3 bucket
	fmt.Print(fileName)
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		ACL:    aws.String("public-read"),
		Key:    aws.String(fileName),
		Body:   file,
	})
	fmt.Print(err)
	if err != nil {
		return nil, err
	}
	return up, nil
}
