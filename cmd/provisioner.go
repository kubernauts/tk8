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

package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/kubernauts/tk8/internal"

	aws "github.com/kubernauts/tk8-provisioner-aws"
	azure "github.com/kubernauts/tk8-provisioner-azure"
	baremetal "github.com/kubernauts/tk8-provisioner-baremetal"
	eks "github.com/kubernauts/tk8-provisioner-eks"
	nutanix "github.com/kubernauts/tk8-provisioner-nutanix"
	openstack "github.com/kubernauts/tk8-provisioner-openstack"
	"github.com/kubernauts/tk8/internal/cluster"

	"github.com/spf13/cobra"
)

var name string
var provisioners = map[string]cluster.Provisioner{
	"aws":       aws.NewAWS(),
	"azure":     azure.NewAzure(),
	"baremetal": baremetal.NewBaremetal(),
	"eks":       eks.NewEKS(),
	"nutanix":   nutanix.NewNutanix(),
	"openstack": openstack.NewOpenstack(),
}

var provisionerInstallCmd = &cobra.Command{
	Use:              "install [" + GetAvaibleProvisioner() + "]",
	TraverseChildren: true,

	Short: "install the infrastructure",
	Long:  `This command will provide the infrastructure which is needed to install and run kubernetes on your selected platform.`,
	Args:  ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Init(args[1:])
			val.Setup(args[1:])
		}
	},
}

var provisionerScaleCmd = &cobra.Command{
	Use:              "scale [" + GetAvaibleProvisioner() + "]",
	TraverseChildren: true,

	Short: "scale up the infrastructure",
	Long:  `This command will scale up the existing infrastructure and the kubernetes cluster to the desired strength as per the config file.`,
	Args:  ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Scale(args[1:])
		}
	},
}

var provisionerResetCmd = &cobra.Command{
	Use:              "reset [" + GetAvaibleProvisioner() + "]",
	TraverseChildren: true,

	Short: "reset the infrastructure",
	Long:  `This command will reset the existing infrastructure and will remove kubernetes from it.`,
	Args:  ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Reset(args[1:])
		}
	},
}

var provisionerRemoveCmd = &cobra.Command{
	Use:              "remove [" + GetAvaibleProvisioner() + "]",
	TraverseChildren: true,

	Short: "scale down the infrastructure",
	Long:  `This command will scale down the existing infrastructure and the kubernetes cluster to the desired strength as per the config file.`,
	Args:  ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Remove(args[1:])
		}
	},
}

var provisionerDestroyCmd = &cobra.Command{
	Use:              "destroy [" + GetAvaibleProvisioner() + "]",
	TraverseChildren: true,
	Short:            "destroy the infrastructure",
	Long:             `This command will destroy the infrastructure which was created with the install command.`,
	Args:             ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Destroy(args[1:])
		}
	},
}

var provisionerUpgradeCmd = &cobra.Command{
	Use:              "upgrade [" + GetAvaibleProvisioner() + "]",
	TraverseChildren: true,
	Short:            "Manages the infrastructure on AWS",
	Long: `
Create, delete and show current status of the deployment that is running on AWS.
Kindly ensure that terraform is installed also.`,
	Args: ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Upgrade(args[1:])
		}
	},
}

func ArgsValidation(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires at least one arg")
	}
	if _, ok := provisioners[args[0]]; ok {
		if _, err := os.Stat("./provisioner/" + args[0]); err == nil {
			return nil
		}
		os.Mkdir("./provisioner", 0755)
		common.CloneGit("./provisioner", "https://github.com/kubernauts/tk8-provisioner-"+args[0], args[0])
		common.ReplaceGit("./provisioner/" + args[0])
		return nil

	}
	return errors.New("provisioner not supported")

}

func GetAvaibleProvisioner() string {
	keys := make([]string, 0, len(provisioners))
	for k := range provisioners {
		keys = append(keys, k)
	}
	return strings.Join(keys, "|")
}
func init() {
	clusterCmd.AddCommand(provisionerInstallCmd)
	clusterCmd.AddCommand(provisionerScaleCmd)
	clusterCmd.AddCommand(provisionerResetCmd)
	clusterCmd.AddCommand(provisionerRemoveCmd)
	clusterCmd.AddCommand(provisionerUpgradeCmd)
	clusterCmd.AddCommand(provisionerDestroyCmd)

	provisionerInstallCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
	provisionerScaleCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
	provisionerResetCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
	provisionerRemoveCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
	provisionerUpgradeCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
	provisionerDestroyCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
}
