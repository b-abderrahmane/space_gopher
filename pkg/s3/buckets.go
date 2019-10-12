package s3

import (
	"fmt"
    "github.com/space_gopher/pkg/generics"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func createBucket(svc *s3.S3, bucketName string) {
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucketName, err)
	} else {
		fmt.Printf("Bucket %s created successfully\n", bucketName)
	}
}

func deleteBucket(svc *s3.S3, bucketName string, purgeBucket bool) {
	if purgeBucket {
		fmt.Printf("Bucket %s contains some elements, those files will be deleted.\n", bucketName)
		iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})

		if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
			exitErrorf("Unable to delete objects from bucket %q, %v", bucketName, err)
		}
		fmt.Printf("Bucket %s content purged successfully\n", bucketName)
	}

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		exitErrorf("Unable to delete bucket %q, %v", bucketName, err)
	} else {
		fmt.Printf("Bucket %s deleted successfully\n", bucketName)
	}
}

func listBucket(svc *s3.S3) {
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
}
