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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Usage:
// go run iam_createaccesskey.go > newuserkeys.txt
func CreateAWSRole() error {

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String("IAM_USER_NAME"),
	})

	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	fmt.Println("Success", *result.AccessKey)

	return nil
}

/*
// CreateAWSRole creates IAM role on AWS
func CreateAWSRole() error {

	// Display starting message
	util.CoriPrintln("Start creating a IAM role on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"

	// Create Session
	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		util.CoriPrintf("Failed to create session : %v\n", err)
		return err
	}

	svc := iam.New(sess)
	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String("<URL-encoded-JSON>"),
		Path:     aws.String("/"),
		RoleName: aws.String("Test-Role"),
	}

	result, err := svc.CreateRole(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(iam.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}
	fmt.Println(result)

	return nil
}
*/
