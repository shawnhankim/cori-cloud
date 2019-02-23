package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

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

	// Get new session
	sess, err := ses.NewSession(&aws.Config{
		Region: aws.String(sampleRegion), // us-west-2
	})
	if err != nil {
		util.CoriPrintln("Failed to get a new session", err)
		return
	}
	util.CoriPrintln("Got a new session")

	// Get EC2 service
	ec2svc := ec2.New(sess)

	instanceID := "i-0c14e9064ee1c53e2"
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			&instanceID,
		},
	}
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		return
	}
	util.CoriPrintln("Got the instance", resp)

	for idx, res := range resp.Reservations {
		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID: ", *inst.InstanceId)
		}
	}
}

/* Sample Instance Data Structure Information


 Reservations: [{
      Instances: [{
          AmiLaunchIndex: 0,
          Architecture: "x86_64",
          BlockDeviceMappings: [{
              DeviceName: "/dev/xvda",
              Ebs: {
                AttachTime: 2019-02-19 20:34:49 +0000 UTC,
                DeleteOnTermination: true,
                Status: "attached",
                VolumeId: "vol-0095858007cdc6215"
              }
            }],
          ClientToken: "",
          CpuOptions: {
            CoreCount: 1,
            ThreadsPerCore: 1
          },
          EbsOptimized: false,
          EnaSupport: true,
          Hypervisor: "xen",
          ImageId: "ami-0cd3dfa4e37921605",
          InstanceId: "i-0c14e9064ee1c53e2",
          InstanceType: "t2.micro",
          KeyName: "shawnkim-ssh",
          LaunchTime: 2019-02-19 20:34:48 +0000 UTC,
          Monitoring: {
            State: "disabled"
          },
          NetworkInterfaces: [{
              Attachment: {
                AttachTime: 2019-02-19 20:34:48 +0000 UTC,
                AttachmentId: "eni-attach-0090da5e9af3b25e8",
                DeleteOnTermination: true,
                DeviceIndex: 0,
                Status: "attached"
              },
              Description: "",
              Groups: [{
                  GroupId: "sg-058b38faa5645b89e",
                  GroupName: "default"
                }],
              MacAddress: "02:5b:44:2b:a5:32",
              NetworkInterfaceId: "eni-0f2032dfdf2c02fe2",
              OwnerId: "641222230151",
              PrivateDnsName: "ip-10-0-0-149.us-east-2.compute.internal",
              PrivateIpAddress: "10.0.0.149",
              PrivateIpAddresses: [{
                  Primary: true,
                  PrivateDnsName: "ip-10-0-0-149.us-east-2.compute.internal",
                  PrivateIpAddress: "10.0.0.149"
                }],
              SourceDestCheck: false,
              Status: "in-use",
              SubnetId: "subnet-05ae0dc4bbfaad921",
              VpcId: "vpc-0e9e0dda2f47cab21"
            }],
          Placement: {
            AvailabilityZone: "us-east-2a",
            GroupName: "",
            Tenancy: "default"
          },
          PrivateDnsName: "ip-10-0-0-149.us-east-2.compute.internal",
          PrivateIpAddress: "10.0.0.149",
          PublicDnsName: "",
          RootDeviceName: "/dev/xvda",
          RootDeviceType: "ebs",
          SecurityGroups: [{
              GroupId: "sg-058b38faa5645b89e",
              GroupName: "default"
            }],
          SourceDestCheck: false,
          State: {
            Code: 16,
            Name: "running"
          },
          StateTransitionReason: "",
          SubnetId: "subnet-05ae0dc4bbfaad921",
          Tags: [
            {
              Key: "f5:blueops:maid_status",
              Value: "Instance marked for stop: stop@2019/02/22"
            },
            {
              Key: "f5:blueops:owner",
              Value: "Team_Atlas"
            },
            {
              Key: "KubernetesCluster",
              Value: "Shawn-sample"
            },
            {
              Key: "Owner",
              Value: "Shawn"
            },
            {
              Key: "Name",
              Value: "shawn-linux-master"
            },
            {
              Key: "AutoPrune",
              Value: "False"
            }
          ],
          VirtualizationType: "hvm",
          VpcId: "vpc-0e9e0dda2f47cab21"
        }],
      OwnerId: "641222230151",
      ReservationId: "r-0cae53ab4c920c8ac"
    }]

*/
