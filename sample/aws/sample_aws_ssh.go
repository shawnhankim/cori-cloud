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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	util "github.com/shawnhankim/cori-cloud/pkg/util"
	"golang.org/x/crypto/ssh"
)

func ExampleExecCmd(inst *CommonInstanceInfo) {
	res, err := ExecCmd(*inst.elasticIP, "ls -l /")
	//res, err := ExecCmd("54.191.245.224", "ls -l /") //whoami")
	if err != nil {
		util.CoriPrintln("Failed to connect ssh", res, err)
	} else {
		util.CoriPrintln("Succeed to connect ssh", res)
	}
}

func ExampleSSH() {
	res, err := ExecCmd("52.41.115.59", "ls -l /") //whoami")
	if err != nil {
		util.CoriPrintln("Failed to connect ssh", *res, err)
	} else {
		util.CoriPrintln("Succeed to connect ssh", *res)
	}
}

func ExecCmd(publicIP, cmd string) (*string, error) {

	// Open PEM file
	pemPath := os.Getenv("PEM_PATH")
	pemBytes, err := ioutil.ReadFile(pemPath)
	if err != nil {
		util.CoriPrintln("Unable to open PEM file.", pemPath, "error: ", err)
		return nil, err
	}

	// Obtain private key
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		util.CoriPrintln("Unable to open private key.", signer, "error: ", err)
		return nil, err
	}

	// Connect to the remote server and perform the SSH handshake
	config := &ssh.ClientConfig{
		User: "centos",
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
	config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	addr := fmt.Sprintf("%s:22", publicIP)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		util.CoriPrintln("Unable to dial to", addr, "error: ", err)
		return nil, err
	}
	util.CoriPrintln("Succeed to dial to", addr)
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		util.CoriPrintln("Unable to create session: ", err)
		return nil, err
	}
	util.CoriPrintln("Succeed to create session")

	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	if err != nil {
		util.CoriPrintln("Unable to get return value: ", err)
		return nil, err
	}
	return aws.String(stdoutBuf.String()), nil
}
