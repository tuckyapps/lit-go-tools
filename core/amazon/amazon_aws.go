package amazon

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofrs/uuid"
	"github.com/tuckyapps/lit-go-tools/core/amazon/s3client"
)

// Holds the s3 service client
var svc *s3.S3

// InitS3 configures configures AWS service to use s3
func InitS3(accessKey string, secretKey string, awsRegion string) (err error) {
	svc, err = s3client.Init(
		accessKey,
		secretKey,
		awsRegion,
	)

	return
}

// GetS3Object retrieves a file from the indicated bucket
func GetS3Object(bucket, filePath string) (resp *s3.GetObjectOutput, err error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	}

	resp, err = svc.GetObject(params)
	return
}

// DeleteS3Object deletes an existing object from the inidcated bucket
func DeleteS3Object(bucket, filePath string) (err error) {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	}

	_, err = svc.DeleteObject(params)
	return
}

// PutS3Object stores the file into the indicated bucket. Value from 'path' allows to store
// the file into a specific folder
func PutS3Object(bucket, path string, fileName string, file multipart.File, fileHeader *multipart.FileHeader) (storedFileName string, err error) {

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	if fileName == "" {
		// create unique name
		uniqueID, _ := uuid.NewV4()
		fileName = uniqueID.String() + filepath.Ext(fileHeader.Filename)
	}

	if path == "" {
		path = "/"
	} else {
		path += "/"
	}

	filePath := fmt.Sprintf("%s%s", path, fileName)

	s3Object := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(filePath),
		ACL:           aws.String(s3.BucketCannedACLPublicRead),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(int64(size)),
		ContentType:   aws.String(http.DetectContentType(buffer)),
		StorageClass:  aws.String(s3.ObjectStorageClassStandard),
		// ContentDisposition:   aws.String("attachment"),
		// ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
	}

	if _, err = svc.PutObject(s3Object); err == nil {
		storedFileName = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, filePath[1:])
	}

	return
}
