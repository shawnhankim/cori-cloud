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

package main

import (
	//"k8s.io/apimachinery/pkg/api/meta"

	"math/big"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	ses "github.com/aws/aws-sdk-go/aws/session"

	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

var (
	sampleKeyID   = "shawnkim-ssh"
	sampleRegion  = "us-east-2"
	sampleProfile = "my-account"
	sampleName    = "Shawn-sample"
)

func main() {
	main1()
}

func main1() {

	// Create a EC2Metadata client from just a session.
	//svc := ec2metadata.New(mySession)

	// Create a EC2Metadata client with additional configuration
	//svc := ec2metadata.New(mySession, aws.NewConfig().WithLogLevel(aws.LogDebugHTTPBody))
	/*
		metaSession, _ := session.NewSession()
		metaClient := ec2metadata.New(metaSession)
		region, _ := metaClient.Region()

		conf := aws.NewConfig().WithRegion(region)
		realSession, _ := session.NewSession(conf)
		dynDb := dynamodb.New(realSession) // Or whatever AWS service we want to use
	*/

	// Display starting message
	start := time.Now()
	util.CoriPrintln("Start testing metadata service on AWS.")
	util.CoriPrintln("- ", new(big.Int).Binomial(1000, 10))

	// Get new session
	metaSession, err := ses.NewSession(&aws.Config{
		Region: aws.String(sampleRegion), // us-west-2
	})
	if err != nil {
		util.CoriPrintln("Failed to get a new session", err)
		return
	}
	util.CoriPrintln("Got a new session")

	// Get metadata service client
	metaClient := ec2metadata.New(metaSession)
	region, err := metaClient.Region()
	if err != nil {
		util.CoriPrintln("Failed to create meta data service client", err)
		return
	}
	util.CoriPrintln("Created AWS metadata service client : ", metaSession)

	// Get metadata service session
	conf := aws.NewConfig().WithRegion(region)
	realSession, _ := ses.NewSession(conf)
	if err != nil {
		util.CoriPrintln("Failed to create meta data service session", err)
		return
	}
	util.CoriPrintln("Created AWS metadata service session : ", realSession)

	//dynDb := dynamodb.New(realSession)

	instanceName := ""

	// Display elapsed time
	elapsed := time.Since(start)
	util.CoriPrintf("Elapsed time : %s [Instance : %s]\n", elapsed, instanceName)
}
