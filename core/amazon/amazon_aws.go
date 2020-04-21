package amazon

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gofrs/uuid"
	"github.com/tuckyapps/lit-go-tools/core/amazon/s3client"
	"github.com/tuckyapps/lit-go-tools/core/amazon/sesclient"
)

var (
	// S3 Holds the s3 service.
	S3 s3Service

	// SES Holds the ses service.
	SES sesService
)

type s3Service struct {
}

type sesService struct {
}

// Holds the s3 service client
var s3ServiceClient *s3.S3

// Holds the ses service client
var sesServiceClient *ses.SES

// Init configures AWS service to use s3 and ses with same credentials.
func Init(accessKey string, secretKey string, awsRegion string) (err error) {
	err = initS3(accessKey, secretKey, awsRegion)

	if err == nil {
		err = initSES(accessKey, secretKey, awsRegion)
	}

	return
}

// ------ S3 OPERATIONS ------

// initS3 configures AWS service to use s3
func initS3(accessKey string, secretKey string, awsRegion string) (err error) {
	s3ServiceClient, err = s3client.Init(
		accessKey,
		secretKey,
		awsRegion,
	)

	return
}

// GetS3Object retrieves a file from the indicated bucket
func (s3s s3Service) GetS3Object(bucket, filePath string) (resp *s3.GetObjectOutput, err error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	}

	resp, err = s3ServiceClient.GetObject(params)
	return
}

// DeleteS3Object deletes an existing object from the inidcated bucket
func (s3s s3Service) DeleteS3Object(bucket, filePath string) (err error) {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	}

	_, err = s3ServiceClient.DeleteObject(params)
	return
}

// PutS3Object stores the file into the indicated bucket. Value from 'path' allows to store
// the file into a specific folder
func (s3s s3Service) PutS3Object(bucket, path string, fileName string, file multipart.File, fileHeader *multipart.FileHeader) (storedFileName string, err error) {

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

	if _, err = s3ServiceClient.PutObject(s3Object); err == nil {
		storedFileName = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, filePath[1:])
	}

	return
}

// ------ SES OPERATIONS ------

// Email options.
type Email struct {
	// From is the source email.
	From string

	// To is a set of destination emails.
	To []string

	// ReplyTo is a set of reply to emails.
	ReplyTo []string

	// Subject is the email subject text.
	Subject string

	// Text is the plain text representation of the body.
	Text string

	// HTML is the HTML representation of the body.
	HTML string
}

// EmailTemplate contains template the necessary data to send a template based email.
type EmailTemplate struct {
	TemplateName string
	From         string
	To           []string
	TemplateData string
	ReplyTo      []string
}

// initSES configures AWS service to use SES
func initSES(accessKey string, secretKey string, awsRegion string) (err error) {
	sesServiceClient, err = sesclient.Init(
		accessKey,
		secretKey,
		awsRegion)

	return
}

// SendEmail an email.
func (sess sesService) SendEmail(e *Email) error {
	if e.HTML == "" {
		e.HTML = e.Text
	}

	msg := &ses.Message{
		Subject: &ses.Content{
			Charset: aws.String("utf-8"),
			Data:    &e.Subject,
		},
		Body: &ses.Body{
			Html: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &e.HTML,
			},
			Text: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &e.Text,
			},
		},
	}

	dest := &ses.Destination{
		ToAddresses: aws.StringSlice(e.To),
	}

	_, err := sesServiceClient.SendEmail(&ses.SendEmailInput{
		Source:           &e.From,
		Destination:      dest,
		Message:          msg,
		ReplyToAddresses: aws.StringSlice(e.ReplyTo),
	})

	return err
}

// SendTemplateEmail is used to send an email based on a template.
func (sess sesService) SendTemplateEmail(et *EmailTemplate) error {

	dest := &ses.Destination{
		ToAddresses: aws.StringSlice(et.To),
	}

	_, err := sesServiceClient.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Source:           &et.From,
		Destination:      dest,
		ReplyToAddresses: aws.StringSlice(et.ReplyTo),
		Template:         &et.TemplateName,
		TemplateData:     &et.TemplateData,
	})

	return err
}
