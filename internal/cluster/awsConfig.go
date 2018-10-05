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
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/kubernauts/tk8/internal/templates"
	"github.com/spf13/viper"
)

// Config holds the variables to be used in the default configuration.
type Config struct {
	AccessKey   string
	SecretKey   string
	ClusterName string
	SSHName     string
}

func namer(name string) Config {
	return Config{
		ClusterName: name,
		SSHName:     name,
	}
}

func generateName() string {
	var (
		words     = flag.Int("words", 2, "The number of words in generated name")
		separator = flag.String("separator", "", "The separator between words in the name")
	)
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	generatedName := petname.Generate(*words, *separator)
	return generatedName
}

func getCreds() (string, string) {
	var accessKey, secretKey string
	fmt.Print("Enter AWS Access Key: ")
	fmt.Scanln(&accessKey)
	fmt.Print("Enter AWS Secret Key: ")
	fmt.Scanln(&secretKey)

	err := os.Setenv("AWS_ACCESS_KEY_ID", accessKey)
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)
	ErrorCheck("Error setting the credentials environment variable: ", err)

	return accessKey, secretKey
}

func checkCredentials() error {
	sess, err := session.NewSession(&aws.Config{})
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})
	_, err = creds.Get()
	if err != nil {
		return err
	}
	return nil
}

// CreateConfig is responsible for creating a default config incase when none is provided.
func CreateConfig() {
	generatedName := generateName()
	fmt.Printf("\nNo default config was provided. Generating one for you...\n")
	err := checkCredentials()
	confStruct := Config{}
	if err != nil {
		accessKey, secretKey := getCreds()
		confStruct = Config{AccessKey: accessKey, SecretKey: secretKey, ClusterName: generatedName, SSHName: generatedName}
	}
	ParseTemplate(templates.Config, "./config.yaml", confStruct)
	ReadViperConfigFile("config")
	region := viper.GetString("aws.aws_default_region")
	CreateSSHKey(generatedName, region)
	fmt.Printf("\nCluster Name:\t%s\nSSH Key name:\t%s\nAWS Region:\t%s\n", generatedName, generatedName, region)
}
