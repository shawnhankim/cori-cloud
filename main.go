package main

import (
	"fmt"

	sampleAWS "github.com/shawnhankim/cori-cloud/sample/aws"
)

func main() {

	// Display title of cori-cloud
	fmt.Printf("This is the cori-cloud.\n")

	// Sample Code : create an instance on AWS
	sampleAWS.CreateAWSInstance()
}
