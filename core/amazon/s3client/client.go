package s3client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Init set credentials to use s3 service
func Init(accessKeyID, secretAccessKey, region string) (svc *s3.S3, err error) {

	var creds *credentials.Credentials

	// credentials of AWS service
	if accessKeyID != "" && secretAccessKey != "" {
		creds = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	}

	// initialize a session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		return
	}

	// create s3 service client
	svc = s3.New(sess)

	return
}
