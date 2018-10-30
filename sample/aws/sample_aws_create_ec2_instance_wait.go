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

	output := InitCommonInstanceInfo()

	// Display starting message
	start := time.Now()
	util.CoriPrintln("Start creating a sample EC2 instance on AWS.")
	util.CoriPrintln("- ", new(big.Int).Binomial(1000, 10))

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintln("Failed to create session", err)
		return nil, err
	}

	// Create EC2 instance session
	output.ec2Service = ec2.New(sess)

	// Run EC2 instance
	input := GetSampleSecurityGroupInput()
	runResult, err := output.ec2Service.RunInstances(input)
	if err != nil {
		util.CoriPrintln("Failed to create instance", err)
		return nil, err
	}

	// Get instance ID and network interface ID
	output.instanceID = runResult.Instances[0].InstanceId
	output.networkInterfaceID = runResult.Instances[0].NetworkInterfaces[0].NetworkInterfaceId
	util.CoriPrintf("Created instance: ID(%s), network ID(%s) \n", *output.instanceID, *output.networkInterfaceID)

	// Modify network interface attribute : SourceDestCheck(False)
	err = ExampleEC2_ModifyNetworkInterfaceAttribute(output.ec2Service, *output.networkInterfaceID)
	if err != nil {
		util.CoriPrintf("Failed to modify network interface attribute", err)
		return output, err
	}
	util.CoriPrintln("Updated network interface: ", runResult.Instances[0].NetworkInterfaces)

	// Check whether public IP is created in the instance on AWS
	statusInput := ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			output.instanceID,
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
	networkAttachedStatusErr := output.ec2Service.WaitUntilInstanceExists(&statusInput)
	if networkAttachedStatusErr != nil {
		util.CoriPrintln("Failed to wait until instances exist: %v", networkAttachedStatusErr)
		return nil, networkAttachedStatusErr
	}
	util.CoriPrintln("Attached all network interfaces to find public IP...")

	// Get the latest EC2 instance information
	statusInput = ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			output.instanceID,
		},
	}
	result, err := output.ec2Service.DescribeInstances(&statusInput)
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
		return output, err
	}

	// Set instance name
	for _, inst := range result.Reservations[0].Instances {
		if *inst.State.Code == 16 { // "running"
			for _, tag := range inst.Tags {
				if *tag.Key == "Name" {
					output.instanceName = tag.Value
					break
				}
			}
			output.isInstanceCreated = true
			output.isNetworkCreated = true
			break
		}
	}

	// Set and display public IP
	output.publicIP = result.Reservations[0].Instances[0].PublicIpAddress
	util.CoriPrintln("Found public IP: ", *output.publicIP)

	/*
		// Create elastic IP
		output.elasticIP, output.elasticAllocationID, err = ExampleEC2_AllocateAddress(output.ec2Service, *output.instanceID, *output.publicIP)
		if err != nil {
			util.CoriPrintln("Failed to create elastic IP", err)
			return output, err
		}
		util.CoriPrintln("Successfully created elasticIP: ", *output.elasticIP)

		// Associate elastic IP to instance
		_, err = ExampleEC2_AssociateAddress(output.ec2Service, *output.instanceID, *output.elasticIP)
		if err != nil {
			util.CoriPrintf("Failed to assotiated elasticIP (%s) to instance (%s), %v \n", *output.elasticIP, *output.instanceID, err)
			return output, err
		}
		//result.Reservations[0].
		util.CoriPrintf("Successfully assotiated elasticIP (%s) to instance (%s) \n", *output.elasticIP, *output.instanceID)
	*/
	// Display elapsed time
	elapsed := time.Since(start)
	util.CoriPrintf("Elapsed time : %s\n", elapsed)

	return output, nil

}
func fmtAddress(addr *ec2.Address) string {
	out := fmt.Sprintf("IP: %s,  allocation id: %s",
		aws.StringValue(addr.PublicIp), aws.StringValue(addr.AllocationId))
	if addr.InstanceId != nil {
		out += fmt.Sprintf(", instance-id: %s", *addr.InstanceId)
	}
	return out
}
