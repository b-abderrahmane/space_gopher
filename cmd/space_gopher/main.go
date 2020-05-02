package main

import (
	"fmt"
	"os"

	"../../pkg/awsstorage"
	"github.com/akamensky/argparse"
)

type S3BucketNamespace struct {
	BucketParser  *argparse.Command
	DeleteCommand *argparse.Command
	CreateCommand *argparse.Command
	ListCommand   *argparse.Command
}

type S3FileNamespace struct {
	fileCommand     *argparse.Command
	uploadCommand   *argparse.Command
	downloadCommand *argparse.Command
	deleteCommand   *argparse.Command
	listCommand     *argparse.Command
}

type S3Namespace struct {
	s3Command       *argparse.Command
	bucketNamespace *S3BucketNamespace
	fileNamespace   *S3FileNamespace
}

type CLIParser struct {
	parser      *argparse.Parser
	s3Namespace *S3Namespace
}

type BucketManifest struct {
	bucket string
	file   map[string]string
}

func main() {

	cliParser := createCLIParser()

	bucketName := cliParser.s3Namespace.s3Command.String("b", "bucket_name", &argparse.Options{Help: "Name of the S3 bucket to be created", Required: false})
	deleteBucketPurge := cliParser.s3Namespace.bucketNamespace.DeleteCommand.Flag("p", "purge", &argparse.Options{Help: "If the bucket is not empty, delete all it's content", Default: false})
	uploadName := cliParser.s3Namespace.fileNamespace.uploadCommand.String("f", "filename", &argparse.Options{Help: "Name/path of the file to be uploaded", Required: true})
	uploadOverwrite := cliParser.s3Namespace.fileNamespace.uploadCommand.Flag("", "overwrite", &argparse.Options{Help: "Name/path of the target bucket", Default: false})
	downloadName := cliParser.s3Namespace.fileNamespace.downloadCommand.String("f", "filename", &argparse.Options{Help: "Name/path of the file to be uploaded", Required: true})

	err := cliParser.parser.Parse(os.Args)
	if err != nil {
		fmt.Println(cliParser.parser.Usage(err))
		return
	}

	if cliParser.s3Namespace.bucketNamespace.CreateCommand.Happened() {
		awsstorage.CreateBucket(*bucketName)
	} else if cliParser.s3Namespace.bucketNamespace.DeleteCommand.Happened() {
		awsstorage.DeleteBucket(*bucketName, *deleteBucketPurge)
	} else if cliParser.s3Namespace.bucketNamespace.ListCommand.Happened() {
		awsstorage.ListBucket()
	} else if cliParser.s3Namespace.fileNamespace.uploadCommand.Happened() {
		awsstorage.UploadFile(*bucketName, *uploadName, *uploadOverwrite)
	} else if cliParser.s3Namespace.fileNamespace.downloadCommand.Happened() {
		awsstorage.DownloadFile(awsstorage.GetAwsSession(), *bucketName, *downloadName)
	} else if cliParser.s3Namespace.fileNamespace.listCommand.Happened() {
		awsstorage.ListFiles(*bucketName)
	}
}

func createCLIParser() *CLIParser {
	var cliParser *CLIParser
	cliParser = new(CLIParser)
	cliParser.parser = argparse.NewParser(os.Args[0], "Prints provided string to stdout")
	cliParser.s3Namespace = new(S3Namespace)
	cliParser.s3Namespace.s3Command = cliParser.parser.NewCommand("s3", "Manage AWS S3 resources")

	cliParser.s3Namespace.bucketNamespace = new(S3BucketNamespace)
	cliParser.s3Namespace.bucketNamespace.BucketParser = cliParser.s3Namespace.s3Command.NewCommand("bucket", "Manage S3 buckets")
	cliParser.s3Namespace.bucketNamespace.DeleteCommand = cliParser.s3Namespace.bucketNamespace.BucketParser.NewCommand("delete", "Delete an S3 bucket")
	cliParser.s3Namespace.bucketNamespace.CreateCommand = cliParser.s3Namespace.bucketNamespace.BucketParser.NewCommand("create", "Create an S3 bucket")
	cliParser.s3Namespace.bucketNamespace.ListCommand = cliParser.s3Namespace.bucketNamespace.BucketParser.NewCommand("list", "List S3 buckets")

	cliParser.s3Namespace.fileNamespace = new(S3FileNamespace)
	cliParser.s3Namespace.fileNamespace.fileCommand = cliParser.s3Namespace.s3Command.NewCommand("file", "Manage S3 buckets content")
	cliParser.s3Namespace.fileNamespace.uploadCommand = cliParser.s3Namespace.fileNamespace.fileCommand.NewCommand("upload", "upload file to an S3 bucket")
	cliParser.s3Namespace.fileNamespace.downloadCommand = cliParser.s3Namespace.fileNamespace.fileCommand.NewCommand("download", "download file to an S3 bucket")
	cliParser.s3Namespace.fileNamespace.listCommand = cliParser.s3Namespace.fileNamespace.fileCommand.NewCommand("list", "list files within a given bucket")
	return cliParser
}
