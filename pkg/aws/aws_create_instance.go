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

// Package aws helps to easily provision and deprovision resources for AWS.
package aws

// CreateInstance creates an instance on AWS
func CreateInstance(region string, profile string, keyName string) error {

	// Create session
	sess, err := CreateSession(region, profile)
	if err != nil {
		return err
	}

	// Create key pair
	err = CreateKeyPair(sess, keyName)
	return err
}
