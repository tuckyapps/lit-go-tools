package sesclient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Init set credentials to use ses service
func Init(accessKeyID, secretAccessKey, region string) (svc *ses.SES, err error) {

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

	// create SES service client
	svc = ses.New(sess)

	return
}
