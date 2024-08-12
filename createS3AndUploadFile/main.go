package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucketName = "test-bucket-7db"
const region = "us-east-1"

func main() {
	var (
		s3Client *s3.Client
		err error
	)

	ctx := context.Background()
	s3Client, err = initS3Client(ctx)
	if err != nil{
		log.Printf("initS3Client error : %v", err)
		os.Exit(1)
	}

	fmt.Println("S3 Client Created Successfully")

	err = createS3Bucket(ctx, s3Client)
	if err != nil {
		log.Printf("creates3Bucket error : %v", err)
		os.Exit(1)
	}

	fmt.Println("S3 Bucket Created Successfully")

	err = uploadToS3Bucket(ctx, s3Client)
	if err != nil {
		log.Printf("uploadToS3Bucket error : %v", err)
		os.Exit(1);
	}

	fmt.Println("Uploaded Successfully")
}

func initS3Client(ctx context.Context) (*s3.Client, error){
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load sdk config error : %v", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func createS3Bucket(ctx context.Context, s3Client *s3.Client) error {

	allBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("listBucket error : %v", err)
	}

	found := false
	for _, bucket := range allBuckets.Buckets{
		if *bucket.Name == bucketName{
			found = true
			break
		}
	}

	if !found {
		_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			// CreateBucketConfiguration: &types.CreateBucketConfiguration{
			// 	LocationConstraint: region,
			// },
		})
	
		if err != nil {
			return fmt.Errorf("createBucket error : %v", err)
		}
	}

	return nil
}

func uploadToS3Bucket(ctx context.Context, s3Client *s3.Client) error{
	uploadFile, err := ioutil.ReadFile("test.txt")

	if err != nil {
		return fmt.Errorf("readFile error: %v", err)
	}

	uploader := manager.NewUploader(s3Client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String("test.txt"),
		// Body: strings.NewReader("Hello Brother !!!!"), // upload this text
		Body: bytes.NewReader(uploadFile), // upload file from your local machine
	})

	if err != nil {
		return fmt.Errorf("upload error : %v", err)
	}

	return nil
}
