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
	coriAWS "github.com/shawnhankim/cori-cloud/pkg/aws"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

// CreateAWSInstance creates an instance on AWS
func CreateAWSInstance() error {

	// Display starting message
	util.CoriPrintln("Start creating a sample instance on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"
	sampleKeyName := "cori-key"

	// Create an instance on AWS
	err := coriAWS.CreateInstance(sampleRegion, sampleProfile, sampleKeyName)
	if err != nil {
		util.CoriPrintln("Unable to create a sample instance on AWS. ")
		return err
	}

	// Display completion message
	util.CoriPrintln("Successfully created a sample instance on AWS.")
	return err
}
