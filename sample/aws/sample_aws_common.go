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
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CommonInstanceInfo is the information which is created
type CommonInstanceInfo struct {
	ec2Service          *ec2.EC2
	isInstanceCreated   bool
	isNetworkCreated    bool
	instanceName        *string
	instanceID          *string
	networkInterfaceID  *string
	elasticIP           *string
	elasticAllocationID *string
	publicIP            *string
}

var (
	sampleKeyID   = "shawnkim-ssh"
	sampleRegion  = "us-west-2"
	sampleProfile = "my-account"
	sampleName    = "Shawn-sample"
)
