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
	"github.com/shawnhankim/cori-cloud/pkg/util"
	sampleAWS "github.com/shawnhankim/cori-cloud/sample/aws"
)

func main() {

	// Display title of cori-cloud
	util.CoriPrintln("* This is the cori-cloud.")

	// Sample Code : create an instance on AWS
	// sampleAWS.CreateAWSInstance()

	// Sample Code : create an instance on AWS
	// sampleAWS.CreateAWSSecurityGroup()

	// Sample Code : create an IAM role on AWS
	// sampleAWS.CreateAWSRole()

	// Sample Code : create an EC2 instance on AWS
	ret, err := sampleAWS.CreateAWSEC2InstanceWitWaitInstanceExists()
	if err == nil {
		sampleAWS.DisplayCommonInstanceInfo(ret)
	}

	// Sample Code : connect ssh to EC2 instance on AWS
	// SampleAWS.ExampleSSH()

	// Sample Code : Get common instance information on AWS
	//inst, err := sampleAWS.GetCommonInstance()

	// Sample Code : terminate instance on AWS
	//if err == nil {
	//		sampleAWS.TerminateInstance(inst)/
	//	}
}
