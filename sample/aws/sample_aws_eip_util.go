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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

func GetElasticIPList(svc *ec2.EC2) {
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("domain"),
				Values: aws.StringSlice([]string{"vpc"}),
			},
		},
	})
	if err != nil {
		util.CoriPrintln("Unable to elastic IP address: ", err)
	}

	// Printout the IP addresses if there are any.
	if len(result.Addresses) == 0 {
		util.CoriPrintf("No elastic IPs for %s region\n", *svc.Config.Region)
	} else {
		util.CoriPrintln("Elastic IPs")
		for _, addr := range result.Addresses {
			util.CoriPrintln("*", fmtAddress(addr))
		}
	}
}

// To allocate an Elastic IP address for EC2-VPC
//
// This example allocates an Elastic IP address to use with an instance in a VPC.
func ExampleEC2_AllocateAddress(svc *ec2.EC2, instanceID, publicIP string) (*string, *string, error) {
	input := &ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	}

	result, err := svc.AllocateAddress(input)
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
		return nil, nil, err
	}

	elasticIP := result.PublicIp
	elasticAllocationID := result.AllocationId
	util.CoriPrintln(result)
	return elasticIP, elasticAllocationID, nil
}

// To associate an Elastic IP address in EC2-Classic
//
// This example associates an Elastic IP address to publicIP with an instance in EC2-Classic.
func ExampleEC2_AssociateAddress(svc *ec2.EC2, instanceID, publicIP string) (*ec2.AssociateAddressOutput, error) {

	input := &ec2.AssociateAddressInput{
		InstanceId: aws.String(instanceID),
		PublicIp:   aws.String(publicIP),
	}

	result, err := svc.AssociateAddress(input)
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
		return nil, err
	}

	util.CoriPrintln(result)
	return result, nil
}

func ReleaseElasticIP(svc *ec2.EC2, allocationID *string) error {
	// Attempt to release the Elastic IP address.
	_, err := svc.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: allocationID,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidAllocationID.NotFound" {
			util.CoriPrintln("Allocation ID %s does not exist", allocationID)
		}
		util.CoriPrintf("Unable to release IP address for allocation %s, %v",
			*allocationID, err)
		return err
	}

	util.CoriPrintf("Successfully released allocation ID %s\n", allocationID)
	return nil
}
