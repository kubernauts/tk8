// Copyright © 2018 The TK8 Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateSSHKey is used to create a new SSH key in AWS for the user.
// It is called when a default config needs to be generated.
func CreateSSHKey(pairName, region string) {
	// Start a new AWS auth session.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create an EC2 service client.
	svc := ec2.New(sess)

	result, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(pairName),
	})
	errorCheck(err)

	// Create the SSH Key on disk.
	sshKey, err := os.OpenFile(pairName, os.O_CREATE|os.O_WRONLY, 0600)
	errorCheck(err)
	fmt.Fprintf(sshKey, *result.KeyMaterial)
	fmt.Printf("\n" + "Successfully created config file and SSH key." + "\n" +
		"You can use the newly created SSH key by adding it to the SSH agent. More info: " +
		"https://www.ssh.com/ssh/add" + "\n\n")
}

func errorCheck(err error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
			ExitErrorf("Specified keypair already exists.")
		}
		ExitErrorf("Error while trying to create the specified key pair: %v.", err)
	}
}
