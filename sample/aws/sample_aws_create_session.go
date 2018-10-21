package sample

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// CreateNewSession creates new session on AWS
func CreateNewSession() (*session.Session, error) {
	//sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "test-account"),
	})
	return sess, err
}
