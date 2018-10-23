//snippet-sourcedescription:[<<FILENAME>> demonstrates how to ...]
//snippet-keyword:[Go]
//snippet-keyword:[Code Sample]
//snippet-service:[<<ADD SERVICE>>]
//snippet-sourcetype:[<<snippet or full-example>>]
//snippet-sourcedate:[]
//snippet-sourceauthor:[AWS]

/*
   Copyright 2010-2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/

package sample

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	coriAWS "github.com/shawnhankim/cori-cloud/pkg/aws"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

// Creates a new security group with the given name and description for
// open port 80 and 22 access. Associating the security group with the
// first VPC in the account if a VPC ID is not provided.
//
// Usage:
//    go run ec2_describe_security_groups.go -n name -d description -vpc vpcID
func CreateAWSSecurityGroup() error {
	var name, desc, vpcID string
	flag.StringVar(&name, "n", "", "Group Name")
	flag.StringVar(&desc, "d", "", "Group Description")
	flag.StringVar(&vpcID, "vpc", "", "(Optional) VPC ID to associate security group with")
	flag.Parse()

	if len(name) == 0 || len(desc) == 0 {
		flag.PrintDefaults()
		exitErrorf("Group name and description require")
	}

	// Display starting message
	util.CoriPrintln("Start creating a sample security group on AWS.")

	// Declare sample variables
	sampleRegion := "us-west-2"
	sampleProfile := "my-account"

	sess, err := coriAWS.CreateSession(sampleRegion, sampleProfile)
	if err != nil {
		return err
	}

	// Create an EC2 service client.
	svc := ec2.New(sess)

	// If the VPC ID wasn't provided in the CLI retrieve the first in the account.
	if len(vpcID) == 0 {
		// Get a list of VPCs so we can associate the group with the first VPC.
		result, err := svc.DescribeVpcs(nil)
		if err != nil {
			exitErrorf("Unable to describe VPCs, %v", err)
		}
		if len(result.Vpcs) == 0 {
			exitErrorf("No VPCs found to associate security group with.")
		}
		vpcID = aws.StringValue(result.Vpcs[0].VpcId)
	}

	// Create the security group with the VPC, name and description.
	createRes, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(name),
		Description: aws.String(desc),
		VpcId:       aws.String(vpcID),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidVpcID.NotFound":
				exitErrorf("Unable to find VPC with ID %q.", vpcID)
			case "InvalidGroup.Duplicate":
				exitErrorf("Security group %q already exists.", name)
			}
		}
		exitErrorf("Unable to create security group %q, %v", name, err)
	}
	fmt.Printf("Created security group %s with VPC %s.\n",
		aws.StringValue(createRes.GroupId), vpcID)

	// Add permissions to the security group
	_, err = svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupName: aws.String(name),
		IpPermissions: []*ec2.IpPermission{
			// Can use setters to simplify seting multiple values without the
			// needing to use aws.String or associated helper utilities.
			(&ec2.IpPermission{}).
				SetIpProtocol("tcp").
				SetFromPort(80).
				SetToPort(80).
				SetIpRanges([]*ec2.IpRange{
					{CidrIp: aws.String("0.0.0.0/0")},
				}),
			(&ec2.IpPermission{}).
				SetIpProtocol("tcp").
				SetFromPort(22).
				SetToPort(22).
				SetIpRanges([]*ec2.IpRange{
					(&ec2.IpRange{}).
						SetCidrIp("0.0.0.0/0"),
				}),
		},
	})
	if err != nil {
		exitErrorf("Unable to set security group %q ingress, %v", name, err)
		return err
	}

	fmt.Println("Successfully set security group ingress")
	return nil
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
