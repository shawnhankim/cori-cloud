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
	"time"

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

	return &ec2.RunInstancesInput{
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String("Shawn-1025-BlueInstanceIAMProfile-696Z2E99TZGV"),
		},
		ImageId:      aws.String("ami-0cf88cd96d9b08d38"),
		InstanceType: aws.String("c4.large"),
		KeyName:      aws.String("shawnkim-ssh"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String("sg-0551fcf43e7f4039f"),
		},
		SubnetId: aws.String("subnet-059d49181a476ccdb"),
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
					{
						Key:   aws.String("KubernetesCluster"),
						Value: aws.String("Shawn-sample")},
				},
			},
		},
	}
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
		util.CoriPrintln("Failed to create session", err)
		return "", err
	}

	// Create EC2 instance session
	svc := ec2.New(sess)

	// Run EC2 instance
	runResult, err := svc.RunInstances(GetSampleSecurityGroupInput())
	if err != nil {
		util.CoriPrintln("Failed to create instance", err)
		return "", err
	}

	// Modify network interface attribute
	instanceID := *runResult.Instances[0].InstanceId
	networkInterfaceId := *runResult.Instances[0].NetworkInterfaces[0].NetworkInterfaceId
	for i := 0; i < 10; i++ {
		attachment := *runResult.Instances[0].NetworkInterfaces[0].Attachment
		if attachment.Status == aws.String(ec2.AttachmentStatusAttached) {
			break
		}
		time.Sleep(time.Second)
	}

	publicIP := ""
	if nil != runResult.Instances[0].NetworkInterfaces[0].Association {
		publicIP = *runResult.Instances[0].NetworkInterfaces[0].Association.PublicIp
	}
	util.CoriPrintf("Created instance : %s, network interface ID: %s, public IP : %s \n",
		instanceID, networkInterfaceId, publicIP)

	err = ExampleEC2_ModifyNetworkInterfaceAttribute(svc, networkInterfaceId)
	if err != nil {
		util.CoriPrintf("Failed to modify network interface attribute", err)
		return "", err
	}
	util.CoriPrintln("Successfully updated network interface: ", networkInterfaceId)
	util.CoriPrintln("network interface: ", runResult.Instances[0].NetworkInterfaces)

	// Create elastic IP
	elasticIP, err := ExampleEC2_CreateElasticIP(svc, "", instanceID)
	//ExampleEC2AssociateIamInstanceProfile(svc, instanceID)
	if err != nil {
		util.CoriPrintln("Failed to create elastic IP", err)
		return "", err
	}
	util.CoriPrintln("Successfully created elasticIP: ", elasticIP)

	return instanceID, nil

}

// To associate an IAM instance profile with an instance
//
// This example associates an IAM instance profile named admin-role with the specified
// instance.
func ExampleEC2AssociateIamInstanceProfile(ec2Svc *ec2.EC2, instanceID string) error { // instance id

	input := &ec2.AssociateIamInstanceProfileInput{
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String("Shawn-1024-09pm-BlueInstanceIAMRole-4S3S9V0G1O60"),
		},
		InstanceId: aws.String(instanceID),
	}

	result, err := ec2Svc.AssociateIamInstanceProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				util.CoriPrintln(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			util.CoriPrintln(err.Error())
		}
		return err
	}
	util.CoriPrintln(result)
	return nil

}

// To modify the sourceDestCheck attribute of a network interface
//
// This example command modifies the sourceDestCheck attribute of the specified network
// interface.
func ExampleEC2_ModifyNetworkInterfaceAttribute(ec2Svc *ec2.EC2, networkInterfaceId string) error {
	input := &ec2.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(networkInterfaceId),
		SourceDestCheck: &ec2.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	}

	result, err := ec2Svc.ModifyNetworkInterfaceAttribute(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				util.CoriPrintln(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			util.CoriPrintln(err.Error())
		}
		return err
	}
	util.CoriPrintln(result)
	return nil
}

func ExampleEC2_CreateElasticIP(ec2Svc *ec2.EC2, publicIP, instanceID string) (string, error) {
	// Attempt to allocate the Elastic IP address.
	allocRes, err := ec2Svc.AllocateAddress(&ec2.AllocateAddressInput{
		PublicIpv4Pool: aws.String(publicIP),
		Domain:         aws.String("vpc"),
	})
	if err != nil {
		util.CoriPrintln("Unable to allocate IP address", err)
		return "", err
	}
	util.CoriPrintln("Allocated IP address", allocRes)

	// Associate the new Elastic IP address with an existing EC2 instance.
	assocRes, err := ec2Svc.AssociateAddress(&ec2.AssociateAddressInput{
		AllocationId: allocRes.AllocationId,
		InstanceId:   aws.String(instanceID),
	})
	if err != nil {
		util.CoriPrintf("Unable to associate IP address with %s, %v\n",
			instanceID, err)
		return "", err
	}

	util.CoriPrintf("Successfully allocated %s with instance %s.\n\tallocation id: %s, association id: %s\n",
		*allocRes.PublicIp, instanceID, *allocRes.AllocationId, *assocRes.AssociationId)
	elasticIP := *allocRes.PublicIp
	return elasticIP, nil
}
