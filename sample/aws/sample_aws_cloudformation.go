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
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func ExampleCloudFormation() {
	// Set stack name, template url
	stackName := "my-groovy-stack"
	templateUrl := "https://my-groovy-bucket.s3.us-west-2.amazonaws.com/my-groovy-template"

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create CloudFormation client in region
	svc := cloudformation.New(sess)

	input := &cloudformation.CreateStackInput{TemplateURL: aws.String(templateUrl), StackName: aws.String(stackName)}

	_, err := svc.CreateStack(input)
	if err != nil {
		fmt.Println("Got error creating stack:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Waiting for stack to be created")

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}
	err = svc.WaitUntilStackCreateComplete(desInput)
	if err != nil {
		fmt.Println("Got error waiting for stack to be created")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Created stack " + stackName)
}
