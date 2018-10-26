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
	"math/big"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	coriAWS "github.com/shawnhankim/cori-cloud/pkg/aws"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

func InitCommonInstanceInfo() *CommonInstanceInfo {
	return &CommonInstanceInfo{
		isInstanceCreated:   false,
		isNetworkCreated:    false,
		instanceName:        aws.String(""),
		instanceID:          aws.String(""),
		networkInterfaceID:  aws.String(""),
		elasticIP:           aws.String(""),
		elasticAllocationID: aws.String(""),
		publicIP:            aws.String(""),
	}
}

// CreateAWSEC2Instance creates an EC2 instance on AWS and wait until instances exists
func CreateAWSEC2InstanceWitWaitInstanceExists() (*CommonInstanceInfo, error) {

	ret := InitCommonInstanceInfo()

	// Display starting message
	start := time.Now()
	util.CoriPrintln("Start creating a sample EC2 instance on AWS.")
	util.CoriPrintln("- ", new(big.Int).Binomial(1000, 10))

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintln("Failed to create session", err)
		return nil, err
	}

	// Create EC2 instance session
	svc := ec2.New(sess)

	// Run EC2 instance
	input := GetSampleSecurityGroupInput()
	runResult, err := svc.RunInstances(input)
	if err != nil {
		util.CoriPrintln("Failed to create instance", err)
		return nil, err
	}

	// Get instance ID and network interface ID
	ret.instanceID = runResult.Instances[0].InstanceId
	ret.networkInterfaceID = runResult.Instances[0].NetworkInterfaces[0].NetworkInterfaceId
	//instanceID := *runResult.Instances[0].InstanceId
	//networkInterfaceID := *runResult.Instances[0].NetworkInterfaces[0].NetworkInterfaceId
	util.CoriPrintf("Created instance: ID(%s), network ID(%s) \n", *ret.instanceID, *ret.networkInterfaceID)

	// Modify network interface attribute : SourceDestCheck(False)
	err = ExampleEC2_ModifyNetworkInterfaceAttribute(svc, *ret.networkInterfaceID)
	if err != nil {
		util.CoriPrintf("Failed to modify network interface attribute", err)
		return ret, err
	}
	util.CoriPrintln("Updated network interface: ", runResult.Instances[0].NetworkInterfaces)

	// Check whether public IP is created in the instance on AWS
	statusInput := ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			ret.instanceID,
		},
		Filters: []*ec2.Filter{
			{
				Name: aws.String("network-interface.attachment.status"),
				Values: []*string{
					aws.String("attached"), // (attaching | attached | detaching | detached).
				},
			},
		},
	}

	// Wait whether all of network interfaces are attached to find public IP
	util.CoriPrintln("Waiting for all network interfaces to be attached to find public IP...")
	networkAttachedStatusErr := svc.WaitUntilInstanceExists(&statusInput)
	if networkAttachedStatusErr != nil {
		util.CoriPrintln("Failed to wait until instances exist: %v", networkAttachedStatusErr)
		return nil, networkAttachedStatusErr
	}
	util.CoriPrintln("Attached all network interfaces to find public IP...")

	// Get the latest EC2 instance information
	statusInput = ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			ret.instanceID,
		},
	}
	result, err := svc.DescribeInstances(&statusInput)
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
		return ret, err
	}

	// Display public IP
	ret.publicIP = result.Reservations[0].Instances[0].PublicIpAddress
	//publicIP := *result.Reservations[0].Instances[0].PublicIpAddress
	util.CoriPrintln("Found public IP: ", *ret.publicIP)

	// Create elastic IP
	ret.elasticIP, ret.elasticAllocationID, err = ExampleEC2_AllocateAddress(svc, *ret.instanceID, *ret.publicIP)
	if err != nil {
		util.CoriPrintln("Failed to create elastic IP", err)
		return ret, err
	}
	util.CoriPrintln("Successfully created elasticIP: ", *ret.elasticIP)

	// Associate elastic IP to instance
	_, err = ExampleEC2_AssociateAddress(svc, *ret.instanceID, *ret.elasticIP)
	if err != nil {
		util.CoriPrintf("Failed to assotiated elasticIP (%s) to instance (%s), %v \n", *ret.elasticIP, *ret.instanceID, err)
		return ret, err
	}
	//result.Reservations[0].
	util.CoriPrintf("Successfully assotiated elasticIP (%s) to instance (%s) \n", *ret.elasticIP, *ret.instanceID)

	// Display elapsed time
	elapsed := time.Since(start)
	util.CoriPrintf("Elapsed time : %s\n", elapsed)

	return ret, nil

}
func fmtAddress(addr *ec2.Address) string {
	out := fmt.Sprintf("IP: %s,  allocation id: %s",
		aws.StringValue(addr.PublicIp), aws.StringValue(addr.AllocationId))
	if addr.InstanceId != nil {
		out += fmt.Sprintf(", instance-id: %s", *addr.InstanceId)
	}
	return out
}
