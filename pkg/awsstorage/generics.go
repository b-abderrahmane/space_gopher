package awsstorage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GlobalRegion is the Default region for now
const GlobalRegion = "us-west-1"

var AWSSession *session.Session = nil

var S3Client *s3.S3 = nil

func GetAwsSession() *session.Session {
	if AWSSession == nil {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(GlobalRegion)},
		)
		if err != nil {
			ExitErrorf("Unable to create an AWS session, %v", err)
		}
		AWSSession = sess
	}
	return AWSSession
}

func GetS3Client() *s3.S3 {
	if S3Client == nil {
		S3Client = s3.New(GetAwsSession())
	}
	return S3Client
}
