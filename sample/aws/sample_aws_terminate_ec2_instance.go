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
	"time"

	"github.com/aws/aws-sdk-go/aws"
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
	err = TerminateAWSEC2Instance(inst.ec2Service, *inst.instanceID)

	// Display elapsed time
	elapsed := time.Since(start)
	util.CoriPrintf("Elapsed time (Get Information) : %s\n", elapsed)
	return err
}

/*
// TerminateInstance is the sample code to terminate EC2 instance on AWS
func TerminateInstance(instanceID string) error {
	// Display starting message
	util.CoriPrintln("Start terminating a sample EC2 instance on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintf("Failed to create session : %v\n", err)
		return err
	}

	// Create EC2 instance session
	svc := ec2.New(sess)

	// Terminate EC2 instnce on AWS
	return TerminateAWSEC2Instance(svc, instanceID)
}
*/

// TerminateAWSEC2Instance terminates an EC2 instance on AWS
func TerminateAWSEC2Instance(svc *ec2.EC2, instanceID string) error {

	// Get input parameter to be terminated
	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	// Germinate EC2 instance on AWS
	resp, err := svc.TerminateInstances(params)
	if err != nil {
		util.CoriPrintf("Failed to terminate instance : %s, error : %v", instanceID, err)
	} else {
		util.CoriPrintf("Successfully terminated instance : %s, %v", instanceID, resp)
	}
	return err
}

/*
// TerminateAWSEC2Instance terminates an EC2 instance on AWS
func TerminateAWSEC2Instance(instanceId string) error {

	// Display starting message
	util.CoriPrintln("Start terminating a sample EC2 instance on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintf("Failed to create session : %v\n", err)
		return err
	}

	// Create EC2 instance
	svc := ec2.New(sess)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		DryRun: aws.Bool(true),
	}

	// Run EC2 instance
	runResult, err := svc.RunInstances(input)

	if err != nil {
		log.Printf("Could not terminate instance : %s, error : %v", instanceID, err)
		return err
	}

	log.Println("Created instance", *runResult.Instances[0].InstanceId)

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("Cori-instance"),
			},
			{
				Key:   aws.String("AutoPrune"),
				Value: aws.String("False"),
			},
			{
				Key:   aws.String("Owner"),
				Value: aws.String("Shawn"),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return err
	}

	log.Println("Successfully tagged instance")
	return nil

}

*/
