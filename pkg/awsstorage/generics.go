package awsstorage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GlobalRegion is the Default region for now
const GlobalRegion = "us-west-1"

func GetAwsSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(GlobalRegion)},
	)
	if err != nil {
		ExitErrorf("Unable to create an AWS session, %v", err)
	}
	return sess
}

func GetS3Client(sess *session.Session) *s3.S3 {
	svc := s3.New(sess)
	return svc
}
