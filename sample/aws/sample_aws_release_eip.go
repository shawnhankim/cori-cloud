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

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

func ReleaseElasticIP(svc *ec2.EC2, allocationID *string) error {

	if allocationID == nil {
		msg := "The instance's elastic IP is already released"
		util.CoriPrintln(msg)
		return errors.New(msg)
	}

	// Attempt to release the Elastic IP address.
	_, err := svc.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: allocationID,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidAllocationID.NotFound" {
			util.CoriPrintln("Allocation ID %s does not exist", allocationID)
		}
		util.CoriPrintf("Unable to release IP address for allocation %s, %v",
			*allocationID, err)
		return err
	}
	return nil
}

// To disassociate an Elastic IP address in EC2-VPC
//
// This example disassociates an Elastic IP address from an instance in a VPC.
func DisassociateAddress(svc *ec2.EC2, elasticIP *string) error {
	if elasticIP == nil {
		msg := "The instance's elastic IP is already released"
		util.CoriPrintln(msg)
		return errors.New(msg)
	}
	input := &ec2.DisassociateAddressInput{
		PublicIp: elasticIP,
	}

	_, err := svc.DisassociateAddress(input)
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
	return nil
}
