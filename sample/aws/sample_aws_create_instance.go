package sample

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	util "github.com/shawnhankim/cori-cloud/pkg/util"
)

// CreateAWSInstance creates an instance on AWS
func CreateAWSInstance() error {

	util.CoriPrintln("Create a sample instance on AWS")
	sess, err := CreateNewSession()
	if err != nil {
		return err
	}
	err = CreateAWSKeyPair(sess)
	return err
}

// CreateAWSKeyPair creates Key Pair on AWS
func CreateAWSKeyPair(sess *session.Session) error {

	svc := ec2.New(session.New())
	input := &ec2.CreateKeyPairInput{
		KeyName: aws.String("cori-key-pair"),
	}

	result, err := svc.CreateKeyPair(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
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
