package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	if len(os.Args) != 3 {
		exitErrorf("Usage: %s action bucket_name",
			os.Args[0])
	}
	action := os.Args[1]
	bucketName := os.Args[2]

	if action == "create" {
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("us-west-1")},
		)

		svc := s3.New(sess)

		_, err := svc.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			exitErrorf("Unable to create bucket %q, %v", bucketName, err)
		} else {
			fmt.Printf("Bucket %s created successfully\n", bucketName)
		}

	} else if action == "delete" {
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("us-west-1")},
		)

		svc := s3.New(sess)

		_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			exitErrorf("Unable to delete bucket %q, %v", bucketName, err)
		} else {
			fmt.Printf("Bucket %s deleted successfully\n", bucketName)
		}

	} else {
		exitErrorf("Unvalid action %q\n", bucketName)
	}

	//_, err := sess.Config.Credentials.Get()

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
