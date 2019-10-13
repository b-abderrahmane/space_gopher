package awsstorage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func fileExists(sess *session.Session, bucketName string, filename string) bool {
	for _, item := range listFiles(GetS3Client(sess), bucketName) {
		if *item.Key == filename {
			return true
		}
	}
	return false
}

func UploadFile(sess *session.Session, bucketName string, file_path string, uploadOverwrite bool) {

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(file_path)
	if err != nil {
		ExitErrorf("failed to open file %q, %v", file_path, err)
	}

	if !uploadOverwrite && fileExists(sess, bucketName, file_path) {
		ExitErrorf("Upload canceled, a file with the same name (%q) already exists", file_path)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file_path),
		Body:   f,
	})
	if err != nil {
		ExitErrorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
}

func DownloadFile(sess *session.Session, bucketName string, file_path string) {

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(file_path)
	if err != nil {
		ExitErrorf("failed to create file %q, %v", file_path, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file_path),
	})
	if err != nil {
		ExitErrorf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)
}

func listFiles(svc *s3.S3, bucketName string) []*s3.Object {
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		ExitErrorf("Unable to list items in bucket %q, %v", bucketName, err)
	}
	return resp.Contents
}

func ListFiles(svc *s3.S3, bucketName string) {
	files := listFiles(svc, bucketName)
	for _, item := range files {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}
