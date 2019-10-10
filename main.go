package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GlobalRegion is the Default region for now
const GlobalRegion = "us-west-1"

func main() {

	if len(os.Args) < 2 {
		exitErrorf("Usage: %s action bucket_name",
			os.Args[0])
	}
	action := os.Args[1]

	if action == "create" {
		bucketName := os.Args[2]
		svc := getS3Client(getAwsSession())

		_, err := svc.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			exitErrorf("Unable to create bucket %q, %v", bucketName, err)
		} else {
			fmt.Printf("Bucket %s created successfully\n", bucketName)
		}

	} else if action == "delete" {
		bucketName := os.Args[2]
		svc := getS3Client(getAwsSession())

		_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			exitErrorf("Unable to delete bucket %q, %v", bucketName, err)
		} else {
			fmt.Printf("Bucket %s deleted successfully\n", bucketName)
		}

	} else if action == "list" {

		svc := getS3Client(getAwsSession())

		result, err := svc.ListBuckets(nil)

		if err != nil {
			exitErrorf("Unable to list buckets")
		}
		if result.Buckets != nil {
			fmt.Println("Buckets:")

			for _, b := range result.Buckets {
				fmt.Printf("* %s created on %s\n",
					aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
			}
		} else {
			fmt.Println("No buckets found.")
		}

	} else {
		exitErrorf("Unvalid action %q\n", action)
	}

	//_, err := sess.Config.Credentials.Get()

}

func getAwsSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(GlobalRegion)},
	)
	if err != nil {
		exitErrorf("Unable to create an AWS session, %v", err)
	}
	return sess
}

func getS3Client(sess *session.Session) *s3.S3 {
	svc := s3.New(sess)

	return svc
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
