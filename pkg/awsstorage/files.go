package awsstorage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// BucketManifestFilename is the file storing the information
const BucketManifestFilename = "bucket-manifest.json"

type FileEntry struct {
	Filename     string
	Size         int64
	LastModified string
	URL          string
}

func fileExists(sess *session.Session, bucketName string, filename string) bool {
	for _, item := range listFiles(GetS3Client(sess), bucketName) {
		if *item.Key == filename {
			return true
		}
	}
	return false
}

func uploadFile(sess *session.Session, bucketName string, filePath string, uploadOverwrite bool, fileContent io.Reader) {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
		Body:   fileContent,
	})
	if err != nil {
		ExitErrorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
}

func UploadFile(sess *session.Session, bucketName string, filePath string, uploadOverwrite bool) {

	fileContent, err := os.Open(filePath)
	if err != nil {
		ExitErrorf("failed to open file %q, %v", filePath, err)
	}

	if !fileExists(sess, bucketName, BucketManifestFilename) {
		fmt.Println("This bucket does not have a manifest file.")
	}

	if !uploadOverwrite && fileExists(sess, bucketName, filePath) {
		ExitErrorf("Upload canceled, a file with the same name (%q) already exists", filePath)
	}
	uploadFile(sess, bucketName, filePath, uploadOverwrite, fileContent)
}

func DownloadFile(sess *session.Session, bucketName string, filePath string) {

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filePath)
	if err != nil {
		ExitErrorf("failed to create file %q, %v", filePath, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
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
	generateFullManifest(bucketName, files)
	for _, item := range files {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}

func generateManifestEntry(bucket string, file *s3.Object) FileEntry {
	url := fmt.Sprintf("https://%s.s3-us-west-1.amazonaws.com/%s", bucket, *file.Key)
	fmt.Println(*file.Key, *file.Size, file.LastModified.String(), url)
	manifestEntry := FileEntry{*file.Key, *(file.Size), file.LastModified.String(), url}
	return manifestEntry
}

func generateFullManifest(bucket string, files []*s3.Object) string {
	var fileEntries []FileEntry
	for _, file := range files {
		manifestEntry := generateManifestEntry(bucket, file)
		fileEntries = append(fileEntries, manifestEntry)
	}
	b, _ := json.Marshal(fileEntries)
	return string(b)
}
