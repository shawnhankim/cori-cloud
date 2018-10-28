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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

// TerminateInstance is the sample code to terminate EC2 instance on AWS
func TerminateInstance(inst *CommonInstanceInfo) error {

	// Display starting message
	start := time.Now()
	util.CoriPrintln("Start terminating a sample EC2 instance on AWS.")

	// Deassociate elastic IP
	err := DisassociateAddress(inst.ec2Service, inst.elasticIP)
	if err != nil {
		util.CoriPrintln("Unable to deassociate instance since elastic IP has been already deassociated", err)
	} else {
		util.CoriPrintln("Complated to deassociate instance", *inst.elasticIP, err)
	}

	// Release elastic IP
	err = ReleaseElasticIP(inst.ec2Service, inst.elasticAllocationID)
	if err != nil {
		util.CoriPrintln("Unable to deassociate instance since elastic IP has been already released", err)
	} else {
		util.CoriPrintln("Completed to release elastic IP", *inst.elasticIP)
	}

	// Terminate EC2 instnce on AWS
	err = TerminateAWSEC2Instance(inst.ec2Service, inst.instanceID)

	// Display elapsed time
	elapsed := time.Since(start)
	util.CoriPrintf("Elapsed time (Get Information) : %s\n", elapsed)
	return err
}

// TerminateAWSEC2Instance terminates an EC2 instance on AWS
func TerminateAWSEC2Instance(svc *ec2.EC2, instanceID *string) error {

	// Check parameter
	if instanceID == nil {
		msg := "There is no instance to be terminated"
		util.CoriPrintln(msg)
		return errors.New(msg)
	}

	// Get input parameter to be terminated
	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			instanceID,
		},
	}

	// Germinate EC2 instance on AWS
	resp, err := svc.TerminateInstances(params)
	if err != nil {
		util.CoriPrintf("Failed to terminate instance : %s, error : %v", *instanceID, err)
	} else {
		util.CoriPrintf("Successfully terminated instance : %s, %v", *instanceID, resp)
	}

	// Check whether instance state is terminated on AWS
	//    * instance-state-code - The state of the instance, as a 16-bit unsigned
	//    integer. The high byte is used for internal purposes and should be ignored.
	//    The low byte is set based on the state represented. The valid values are:
	//    0 (pending), 16 (running), 32 (shutting-down), 48 (terminated), 64 (stopping),
	//    and 80 (stopped).
	//
	//    * instance-state-name - The state of the instance (pending | running |
	//    shutting-down | terminated | stopping | stopped).
	statusInput := ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			instanceID,
		},
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-code"),
				Values: []*string{
					aws.String("48"), // 0 (pending), 16 (running), 32 (shutting-down), 48 (terminated), 64 (stopping), 80 (stopped).
				},
			},
		},
	}

	util.CoriPrintln("Waiting for the instance to be terminated...")
	instanceStateRet := svc.WaitUntilInstanceTerminated(&statusInput)
	if instanceStateRet != nil {
		util.CoriPrintln("Failed to wait until instance is terminated: %v", instanceStateRet)
		return instanceStateRet
	}

	// Get the latest EC2 instance information
	statusInput = ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			instanceID,
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
		return err
	}
	util.CoriPrintf("The terminated instance information : %+v \n", result)
	util.CoriPrintln("The instance is terminated")

	return nil
}
