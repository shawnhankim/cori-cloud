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

// Package main provides code examples how to provision/deprovision resources
// on the multiple cloud providers to help software engineers.
package main

import (
	"fmt"

	SampleAWS "github.com/shawnhankim/cori-cloud/sample/aws"
)

func main() {

	// Display title of cori-cloud
	fmt.Printf("* This is the cori-cloud.\n")

	// Sample Code : create an instance on AWS
	// sampleAWS.CreateAWSInstance()

	// Sample Code : create an instance on AWS
	// sampleAWS.CreateAWSSecurityGroup()

	// Sample Code : create an IAM role on AWS
	// sampleAWS.CreateAWSRole()

	// Sample Code : create an EC2 instance on AWS
	SampleAWS.ExampleSSH()

	//ret, err := sampleAWS.CreateAWSEC2InstanceWitWaitInstanceExists()
	//if err == nil {
	//sampleAWS.ExampleExecCmd(ret)
	//util.CoriPrintln("Waiting 60 seconds to terminate instance")
	//time.Sleep(60 * time.Second)
	//sampleAWS.TerminateInstance(ret)
	//util.CoriPrintln("Terminated instance")
	//}
	//sampleAWS.CreateAWSEC2Instance()
	//sampleAWS.ExampleEC2CreateLaunchTemplate()

	// Sample Code : terminate an EC2 instance on AWS
	//sampleAWS.TerminateInstance("i-0ff22d13ab2c30744")

}
