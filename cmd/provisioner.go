package cmd

import (
	"errors"
	"strings"

	"github.com/kubernauts/tk8/internal/cluster"
	aws "github.com/kubernauts/tk8/internal/provisioner/aws"
	azure "github.com/kubernauts/tk8/internal/provisioner/azure"
	baremetal "github.com/kubernauts/tk8/internal/provisioner/baremetal"
	eks "github.com/kubernauts/tk8/internal/provisioner/eks"
	nutanix "github.com/kubernauts/tk8/internal/provisioner/nutanix"
	openstack "github.com/kubernauts/tk8/internal/provisioner/openstack"
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
	Short:            "Manages the infrastructure on AWS",
	Long: `
Create, delete and show current status of the deployment that is running on AWS.
Kindly ensure that terraform is installed also.`,
	Args: ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Init()
			val.Setup()
		}
	},
}

var provisionerDestroyCmd = &cobra.Command{
	Use:              "destroy [" + GetAvaibleProvisioner() + "]",
	TraverseChildren: true,
	Short:            "Manages the infrastructure on AWS",
	Long: `
Create, delete and show current status of the deployment that is running on AWS.
Kindly ensure that terraform is installed also.`,
	Args: ArgsValidation,
	Run: func(cmd *cobra.Command, args []string) {
		if val, ok := provisioners[args[0]]; ok {
			val.Destroy()
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
			val.Upgrade()
		}
	},
}

func ArgsValidation(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires at least one arg")
	}
	if _, ok := provisioners[args[0]]; ok {
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
	clusterCmd.AddCommand(provisionerUpgradeCmd)
	clusterCmd.AddCommand(provisionerDestroyCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// awsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// awsCmd.Flags().BoolVarP(&install, "install", "i", false, "Install Kubernetes on the AWS infrastructure")
	// Flags to initiate the terraform installation
	// awsCmd.Flags().BoolVarP(&create, "create", "c", false, "Deploy the AWS infrastructure using terraform")
	// Flag to destroy the AWS infrastructure using terraform
	provisionerInstallCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
	provisionerUpgradeCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
	provisionerDestroyCmd.Flags().StringVar(&cluster.Name, "name", cluster.Name, "name of the cluster workspace")
}
