// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"os"

	"github.com/kubernauts/tk8/internal/cluster"
	"github.com/spf13/cobra"
)

var install, create, destroy bool

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Manages the infrastructure on AWS",
	Long: `
Create, delete and show current status of the deployment that is running on AWS.
Kindly ensure that terraform is installed also.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if install {
			cluster.AWSInstall()
		}

		if create {
			cluster.AWSCreate()
		}

		if destroy {
			cluster.AWSDestroy()
		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

func init() {
	clusterCmd.AddCommand(awsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// awsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	awsCmd.Flags().BoolVarP(&install, "install", "i", false, "Install Kubernetes on the AWS infrastructure")
	// Flags to initiate the terraform installation
	awsCmd.Flags().BoolVarP(&create, "create", "c", false, "Deploy the AWS infrastructure using terraform")
	// Flag to destroy the AWS infrastructure using terraform
	awsCmd.Flags().BoolVarP(&destroy, "destroy", "d", false, "Destroy the AWS infrastructure")
}
