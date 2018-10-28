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
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	coriAWS "github.com/shawnhankim/cori-cloud/pkg/aws"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

// To describe the instances with a specific instance type
//
// This example describes the instances with the t2.micro instance type.
func DescribeInstances(svc *ec2.EC2, instanceName *string) (*ec2.DescribeInstancesOutput, error) {

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					instanceName,
				},
			},
		},
	}

	output, err := svc.DescribeInstances(input)
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

	//util.CoriPrintln("Describe instance", result)
	return output, nil
}

func GetCommonInstance() {

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintln("Failed to create session", err)
		return
	}

	// Create EC2 instance session
	svc := ec2.New(sess)
	GetCommonInstanceInfo(svc, aws.String(sampleName))
}

func GetCommonInstanceInfo(svc *ec2.EC2, instanceName *string) (*CommonInstanceInfo, error) {

	// Describe instances
	ret, err := DescribeInstances(svc, instanceName)
	if err != nil {
		util.CoriPrintf("Unable to get common instance information: %v\n", err)
		return nil, err
	}

	// Get common instance information
	output := new(CommonInstanceInfo)
	for _, inst := range ret.Reservations[0].Instances {
		if *inst.State.Code == 16 { // "running"
			for _, tag := range inst.Tags {
				if *tag.Key == "Name" {
					output.instanceName = tag.Value
					break
				}
			}
			output.instanceID = inst.InstanceId
			output.elasticIP = inst.PublicIpAddress
			output.networkInterfaceID = inst.NetworkInterfaces[0].NetworkInterfaceId
			output.publicIP = inst.NetworkInterfaces[0].Association.PublicIp
			output.isInstanceCreated = true
			output.isNetworkCreated = true
			break
		}
	}

	// Describe elastic IP
	output.elasticAllocationID, err = GetElasticAssociationID(svc, output.elasticIP)

	// Display common instance information
	DisplayCommonInstanceInfo(output)

	return output, err
}

func GetElasticAssociationID(svc *ec2.EC2, elasticIP *string) (*string, error) {
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("public-ip"),
				Values: aws.StringSlice([]string{*elasticIP}),
			},
		},
	})
	if err != nil {
		util.CoriPrintln("Unable to elastic association ID: ", err)
		return nil, err
	}

	// Printout the IP addresses if there are any.
	if len(result.Addresses) == 0 {
		err := fmt.Sprintf("No elastic IPs for %s region\n", *svc.Config.Region)
		util.CoriPrintf(err)
		return nil, errors.New(err)
	}
	allocationID := result.Addresses[0].AllocationId
	util.CoriPrintln("Elastic association ID", allocationID)
	return allocationID, nil
}

func DisplayCommonInstanceInfo(inst *CommonInstanceInfo) {

	util.CoriPrintf("+-----------------------------------------------------------+\n")
	util.CoriPrintf("| Common Instance Information                               |\n")
	util.CoriPrintf("+-----------------------------------------------------------+\n")
	util.CoriPrintf("  - isInstanceCreated   : %v \n", inst.isInstanceCreated)
	util.CoriPrintf("  - isNetworkCreated    : %v \n", inst.isNetworkCreated)
	util.CoriPrintf("  - instanceName        : %s \n", *inst.instanceName)
	util.CoriPrintf("  - instanceID          : %s \n", *inst.instanceID)
	util.CoriPrintf("  - networkInterfaceID  : %s \n", *inst.networkInterfaceID)
	util.CoriPrintf("  - elasticIP           : %s \n", *inst.elasticIP)
	util.CoriPrintf("  - elasticAllocationID : %s \n", *inst.elasticAllocationID)
	util.CoriPrintf("  - publicIP            : %s \n\n", *inst.publicIP)
}