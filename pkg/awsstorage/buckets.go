package awsstorage

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CreateBucket(svc *s3.S3, bucketName string) {
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		ExitErrorf("Unable to create bucket %q, %v", bucketName, err)
	} else {
		fmt.Printf("Bucket %s created successfully\n", bucketName)
	}
}

func DeleteBucket(svc *s3.S3, bucketName string, purgeBucket bool) {
	if purgeBucket {
		fmt.Printf("Bucket %s contains some elements, those files will be deleted.\n", bucketName)
		iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})

		if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
			ExitErrorf("Unable to delete objects from bucket %q, %v", bucketName, err)
		}
		fmt.Printf("Bucket %s content purged successfully\n", bucketName)
	}

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		ExitErrorf("Unable to delete bucket %q, %v", bucketName, err)
	} else {
		fmt.Printf("Bucket %s deleted successfully\n", bucketName)
	}
}

func ListBucket(svc *s3.S3) {
	result, err := svc.ListBuckets(nil)

	if err != nil {
		ExitErrorf("Unable to list buckets", err)
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
}

func ExitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
