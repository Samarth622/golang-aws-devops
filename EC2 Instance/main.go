package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	var (
		instanceId string // instance variable
		err error // error variable
	)

	ctx := context.Background()  // create default context
	region := "us-east-1"  // region of the EC2 instance

	if instanceId, err = createEC2(ctx, region); err != nil {
		fmt.Printf("EC2 Instance error : %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Instance Id : %s\n", instanceId)
}

func createEC2(ctx context.Context, region string) (string, error) {   // function to create ec2 instance which take context, and region
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))  // configure with context and region
    if err != nil {
        return "", fmt.Errorf("unable to load SDK config, %v", err)
    }

	ec2Client := ec2.NewFromConfig(cfg)  // Create a new EC2 client using the loaded configuration.

	// Describe the key pairs in the EC2 instance with the specified key name ("go-aws-demo").
	keypairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		KeyNames: []string{"go-aws-demo"},
	})

	if err != nil && !strings.Contains(err.Error(), "InvalidKeyPair.NotFound"){
        return "", fmt.Errorf("describeKeyPairs error : %v", err)
    }

	if keypairs == nil || len(keypairs.KeyPairs) == 0{
		_, err = ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{  // if keypair already not exist then create keypair
			KeyName: aws.String("go-aws-demo"),
	
		})
	
		if err != nil {
			return "", fmt.Errorf("create key pair error")
		}
	}

	imageOut, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{  // here we describe image to create 
		Filters : []types.Filter{
			{
				Name: aws.String("name"),
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"},  // image name
			},
			{
				Name: aws.String("virtualization-type"), // image type
				Values: []string{"hvm"},
			},
		},
		Owners: []string{"099720109477"},  // owner id of the image
	})

	if err != nil {
        return "", fmt.Errorf("DescribeImages error : %v", err)
    }

	if len(imageOut.Images) == 0 {
		return "", fmt.Errorf("imageOut.Image is empty")
	}

	instance, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{  // here we run instance from 
		ImageId: imageOut.Images[0].ImageId, 		// image id
		KeyName: aws.String("go-aws-demo"),			// keypair
		InstanceType: types.InstanceTypeT3Micro,	// instance type
		MinCount: aws.Int32(1),
		MaxCount: aws.Int32(1),
	})

	if err != nil {
		return "", fmt.Errorf("RunInstances error: %v", err)
	}

	if len(instance.Instances) == 0{
		return "", fmt.Errorf("instance.Instances is empty")
	}

	return *instance.Instances[0].InstanceId, nil  // return instance id if all things are good
}