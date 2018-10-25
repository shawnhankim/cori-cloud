// Copyright 2018 The Cori Cloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package sample provides code examples how to provision/deprovision resources
// on the multiple cloud providers to help software engineers.
package sample

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	coriAWS "github.com/shawnhankim/cori-cloud/pkg/aws"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

/*
// GetSampleEC2InstanceInput returns EC2 instance input parameter
func GetSampleEC2InstanceInput() *ec2.RunInstancesInput {
	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(100),
				},
			},
			{
				DeviceName:  aws.String("/dev/sda1"),
				VirtualName: aws.String("ephemeral1"),
			},
		},
		ImageId:      aws.String("ami-0690d7168760bcb2d"),
		InstanceType: aws.String("t2.large"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		NetworkInterfaces: []*ec2.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest{
			{
				AssociatePublicIpAddress: aws.Bool(true),
				DeviceIndex:              aws.Int64(0),
				Ipv6AddressCount:         aws.Int64(1),
				SubnetId:                 aws.String("subnet-7b16de0c"),
			},
		},
	}
	return input
}
*/

// ExampleEC2CreateLaunchTemplate creates EC2 instance and returns instance ID
func ExampleEC2CreateLaunchTemplate() (string, error) {

	// Display starting message
	util.CoriPrintln("Start creating a sample EC2 instance on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintf("Failed to create session : %v\n", err)
		return "", err
	}

	svc := ec2.New(sess)
	input := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: &ec2.RequestLaunchTemplateData{
			ImageId:      aws.String("ami-0690d7168760bcb2d"),
			InstanceType: aws.String("t2.large"),
			KeyName:      aws.String("shawn-ssh"),
			NetworkInterfaces: []*ec2.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest{
				{
					AssociatePublicIpAddress: aws.Bool(true),
					DeviceIndex:              aws.Int64(0),
					Ipv6AddressCount:         aws.Int64(1),
					SubnetId:                 aws.String("subnet-059d49181a476ccdb"),
				},
			},
			SecurityGroupIds: []*string{
				aws.String("sg-002b4e2ccb97b66c7"),
			},
			TagSpecifications: []*ec2.LaunchTemplateTagSpecificationRequest{
				{
					ResourceType: aws.String("instance"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("Shawn-sample")},
						{
							Key:   aws.String("AutoPrune"),
							Value: aws.String("False")},
						{
							Key:   aws.String("Owner"),
							Value: aws.String("Shawn")},
					},
				},
			},
		},
		LaunchTemplateName: aws.String("shawn-template"),
		VersionDescription: aws.String("shawn-sample-version1"),
	}

	result, err := svc.CreateLaunchTemplate(input)
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
		return "", err
	}
	fmt.Println(result)
	return "", nil
}

/*
// To associate an IAM instance profile with an instance
//
// This example associates an IAM instance profile named admin-role with the specified
// instance.
func ExampleEC2AssociateIamInstanceProfile() {
	svc := ec2.New(session.New())
	input := &ec2.AssociateIamInstanceProfileInput{
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String("shawn-sample-iam-role"),
		},
		InstanceId: aws.String("i-123456789abcde123"),
	}

	result, err := svc.AssociateIamInstanceProfile(input)
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
*/

func GetSampleCommonInput() *ec2.RunInstancesInput {
	return &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0690d7168760bcb2d"),
		InstanceType: aws.String("t2.large"),
		KeyName:      aws.String("shawnkim-ssh"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("Shawn-sample")},
					{
						Key:   aws.String("AutoPrune"),
						Value: aws.String("False")},
					{
						Key:   aws.String("Owner"),
						Value: aws.String("Shawn")},
				},
			},
		},
	}
}

func GetSampleNetworkInput() *ec2.RunInstancesInput {
	input := GetSampleCommonInput()
	input.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
		{
			AssociatePublicIpAddress: aws.Bool(true),
			DeviceIndex:              aws.Int64(0),
			SubnetId:                 aws.String("subnet-059d49181a476ccdb"),
		},
	}
	return input
}

func GetSampleSecurityGroupInput() *ec2.RunInstancesInput {
	input := GetSampleCommonInput()
	input.SecurityGroupIds = []*string{
		aws.String("sg-002b4e2ccb97b66c7"), //sg-002b4e2ccb97b66c7"),
	}
	return input
}

func GetSampleInstanceInput() *ec2.RunInstancesInput {
	input := GetSampleNetworkInput()
	input.SecurityGroups = []*string{
		aws.String("shawnkim-ssh"),
	}
	return input
}

// CreateAWSEC2Instance creates an EC2 instance on AWS
func CreateAWSEC2Instance() (string, error) {

	// Display starting message
	util.CoriPrintln("Start creating a sample EC2 instance on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintf("Failed to create session : %v\n", err)
		return "", err
	}

	// Create EC2 instance session
	svc := ec2.New(sess)

	// Run EC2 instance
	runResult, err := svc.RunInstances(GetSampleSecurityGroupInput())
	if err != nil {
		log.Println("Could not create instance", err)
		return "", err
	}
	instanceID := *runResult.Instances[0].InstanceId
	log.Println("Created instance", instanceID)

	// Attach security group to EC2 instance
	/*
		input := &ec2.DescribeSecurityGroupsInput{
			GroupIds: []*string{
				aws.String("sg-002b4e2ccb97b66c7"),
			},
		}

			result, err := svc.DescribeSecurityGroups(input)
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
				return instanceID, err
			}
			if err != nil {
				log.Println("Could not attach security group to EC2 instance", err)
				return "", err
			}
			util.CoriPrintf("Attached security group : %v\n", result)
	*/
	return instanceID, nil

}
