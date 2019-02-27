package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	ses "github.com/aws/aws-sdk-go/aws/session"
)

var (
	sampleRegion = "us-west-1"
)

func main() {
	TestGetProductCode()
}

// TestGetProductCode is the example to get product code through AWS Marketplace
func TestGetProductCode() {

	// Get new session
	sess, err := ses.NewSession(&aws.Config{
		Region: aws.String(sampleRegion),
	})
	if err != nil {
		fmt.Println("Failed to get a new session: ", err)
		return
	}
	fmt.Println("Got a new session")

	// Get EC2 service
	ec2svc := ec2.New(sess)

	//instanceID := "i-001bebdda2a630711"
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			// {
			// 	Name: aws.String("image-id"),
			// 	Values: []*string{
			// 		aws.String("ami-05e518340f5c4837a"),
			// 	},
			// },
			{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String("i-0381255a7afd3ba0f"),
				},
			},
		},
		//InstanceIds: []*string{
		//	&instanceID,
		//},
	}
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		return
	}
	fmt.Println("Got the instance", resp)
	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("Test product code")
	fmt.Println("\n-----------------------------------------------------")

	for idx, res := range resp.Reservations {
		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID : ", *inst.InstanceId)
			fmt.Printf("    - Product Code: %s\n", inst.ProductCodes)
			for _, prod := range inst.ProductCodes {
				fmt.Println("    - Product Code: ", *prod.ProductCodeId)
				fmt.Println("    - Product Type: ", *prod.ProductCodeType)
			}
		}
	}
	/*
		- Product Code:  e23svmiz0xn3z81073a0s9noc
		- Product Type:  marketplace
	*/
}

// To confirm the product instance
// ExampleEC2_ConfirmProductInstance_shared00
// This example determines whether the specified product code is associated with the
// specified instance.
func ExampleEC2_ConfirmProductInstance_shared00() {
	sess, err := ses.NewSession(&aws.Config{
		Region: aws.String(sampleRegion), // us-west-2
	})
	svc := ec2.New(sess)

	input := &ec2.ConfirmProductInstanceInput{
		InstanceId:  aws.String("i-0f63d17d222e9a375"),
		ProductCode: aws.String("774F4FF8"),
	}
	result, err := svc.ConfirmProductInstance(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}
