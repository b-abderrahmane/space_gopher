package awsstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

func generateManifestEntry(bucketName string, fileName string, fileSize int64, fileLastModified string) FileEntry {
	url := fmt.Sprintf("https://%s.s3-us-west-1.amazonaws.com/%s", bucketName, fileName)
	manifestEntry := FileEntry{fileName, fileSize, fileLastModified, url}
	return manifestEntry
}

func generateFullManifest(bucket string, files []*s3.Object) string {
	var fileEntries []FileEntry
	for _, file := range files {
		manifestEntry := generateManifestEntry(bucket, *file.Key, *(file.Size), file.LastModified.String())
		fileEntries = append(fileEntries, manifestEntry)
	}
	jsonFileEntries, _ := json.Marshal(fileEntries)
	return string(jsonFileEntries)
}

func generateUpdatedManifestDefinition(bucket string, fullManifest string, newFile FileEntry) string {
	var fileEntries []FileEntry

	manifestEntry := generateManifestEntry(bucket, newFile.Filename, newFile.Size, newFile.LastModified)
	err := json.Unmarshal([]byte(fullManifest), &fileEntries)
	fileEntries = append(fileEntries, manifestEntry)

	if err != nil {
		ExitErrorf("Unable to parse existing manifest, therefor unable to update it", err)
	}
	jsonFileEntries, _ := json.Marshal(fileEntries)
	return string(jsonFileEntries)
}

func updateManifest(bucketName string, file *os.File, fileName string, oldManifest string) {
	if oldManifest == "" {
		oldManifest = generateFullManifest(bucketName, listFiles(GetS3Client(), bucketName))
	}

	fi, err := file.Stat()
	if err != nil {
		ExitErrorf("Unable to get file information", err)
	}

	manifestEntry := generateManifestEntry(bucketName, fileName, fi.Size(), time.Now().String())
	oldManifest = generateUpdatedManifestDefinition(bucketName, oldManifest, manifestEntry)
	uploadFile(bucketName, BucketManifestFilename, false, strings.NewReader(oldManifest))
}
