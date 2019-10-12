package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3BucketNamespace struct {
	bucketParser  *argparse.Command
	deleteCommand *argparse.Command
	createCommand *argparse.Command
	listCommand   *argparse.Command
}

type S3FileNamespace struct {
	fileCommand     *argparse.Command
	uploadCommand   *argparse.Command
	downloadCommand *argparse.Command
	deleteCommand   *argparse.Command
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

// GlobalRegion is the Default region for now
const GlobalRegion = "us-west-1"

func main() {

	cliParser := createCLIParser()

	createBucketName := cliParser.s3Namespace.bucketNamespace.createCommand.String("n", "name", &argparse.Options{Help: "Name of the S3 bucket to be created", Required: true})
	deleteBucketName := cliParser.s3Namespace.bucketNamespace.deleteCommand.String("n", "name", &argparse.Options{Help: "Name of the S3 bucket to be deleted", Required: true})
	deleteBucketPurge := cliParser.s3Namespace.bucketNamespace.deleteCommand.Flag("p", "purge", &argparse.Options{Help: "If the bucket is not empty, delete all it's content", Default: false})

	err := cliParser.parser.Parse(os.Args)
	if err != nil {
		fmt.Println(cliParser.parser.Usage(err))
		return
	}

	if cliParser.s3Namespace.bucketNamespace.createCommand.Happened() {
		svc := getS3Client(getAwsSession())

		_, err := svc.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(*createBucketName),
		})

		if err != nil {
			exitErrorf("Unable to create bucket %q, %v", *createBucketName, err)
		} else {
			fmt.Printf("Bucket %s created successfully\n", *createBucketName)
		}

	} else if cliParser.s3Namespace.bucketNamespace.deleteCommand.Happened() {
		svc := getS3Client(getAwsSession())

		if *deleteBucketPurge {
			fmt.Printf("Bucket %s contains some elements, those files will be deleted.\n", *deleteBucketName)
			iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
				Bucket: aws.String(*deleteBucketName),
			})

			if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
				exitErrorf("Unable to delete objects from bucket %q, %v", *deleteBucketName, err)
			}
			fmt.Printf("Bucket %s content purged successfully\n", *deleteBucketName)
		}

		_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(*deleteBucketName),
		})

		if err != nil {
			exitErrorf("Unable to delete bucket %q, %v", *deleteBucketName, err)
		} else {
			fmt.Printf("Bucket %s deleted successfully\n", *deleteBucketName)
		}

	} else if cliParser.s3Namespace.bucketNamespace.listCommand.Happened() {

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

	}

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

func createCLIParser() *CLIParser {
	var cliParser *CLIParser
	cliParser = new(CLIParser)
	cliParser.parser = argparse.NewParser(os.Args[0], "Prints provided string to stdout")
	cliParser.s3Namespace = new(S3Namespace)
	cliParser.s3Namespace.s3Command = cliParser.parser.NewCommand("s3", "Manage AWS S3 resources")
	cliParser.s3Namespace.bucketNamespace = new(S3BucketNamespace)
	cliParser.s3Namespace.bucketNamespace.bucketParser = cliParser.s3Namespace.s3Command.NewCommand("bucket", "Manage S3 buckets")
	cliParser.s3Namespace.bucketNamespace.deleteCommand = cliParser.s3Namespace.bucketNamespace.bucketParser.NewCommand("delete", "Delete an S3 bucket")
	cliParser.s3Namespace.bucketNamespace.createCommand = cliParser.s3Namespace.bucketNamespace.bucketParser.NewCommand("create", "Create an S3 bucket")
	cliParser.s3Namespace.bucketNamespace.listCommand = cliParser.s3Namespace.bucketNamespace.bucketParser.NewCommand("list", "List S3 buckets")
	cliParser.s3Namespace.fileNamespace = new(S3FileNamespace)
	cliParser.s3Namespace.fileNamespace.fileCommand = cliParser.s3Namespace.s3Command.NewCommand("file", "Manage S3 buckets content")
	return cliParser
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
