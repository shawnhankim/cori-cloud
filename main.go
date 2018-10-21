package main

import (
	"fmt"

	util "github.com/shawnhankim/cori-cloud/pkg/util"
	sampleAWS "github.com/shawnhankim/cori-cloud/sample/aws"
)

func main() {
	fmt.Printf("This is the cori-cloud.\n")

	err := sampleAWS.CreateAWSInstance()
	if err != nil {
		util.CoriPrintf("Unable to create AWS instance! Error : %v", err)
	} else {
		util.CoriPrintf("Successfully created AWS instance! \n")
	}
}
