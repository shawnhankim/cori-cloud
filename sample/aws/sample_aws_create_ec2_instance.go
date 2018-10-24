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
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	coriAWS "github.com/shawnhankim/cori-cloud/pkg/aws"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

// CreateAWSEC2Instance creates an EC2 instance on AWS
func CreateAWSEC2Instance() error {

	// Display starting message
	util.CoriPrintln("Start creating a sample EC2 instance on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"
	sampleAMIid := "ami-0dfeec5739bc7c5d7" // ami-e7527ed7
	sampleInstType := "t2.large"

	// Create session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintf("Failed to create session : %v\n", err)
		return err
	}

	// Create EC2 instance
	svc := ec2.New(sess)

	// Run EC2 instance
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(sampleAMIid),
		InstanceType: aws.String(sampleInstType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		log.Println("Could not create instance", err)
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
