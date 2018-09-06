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
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kubernauts/tk8/internal/addon"
	"github.com/spf13/cobra"
)

var monitor, rancher bool

// addonCmd represents the addon command
var addonCmd = &cobra.Command{
	Use:   "addon [command]",
	Short: "Manage addon packages",
	Long:  `Manage your addons on kubernetes with this cli.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

// addonCmd represents the addon command
var addonInstallCmd = &cobra.Command{
	Use:   "install [addon name or git url]",
	Short: "Install a addon on your kubernetes cluster",
	Long:  `Install additional packages on top of your kubernetes deployment.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		addon.InstallAddon(args[0])
	},
}

// addonCmd represents the addon command
var addonDestroyCmd = &cobra.Command{
	Use:   "destroy [addon name]",
	Short: "Destroy a running addon on your kubernetes cluster",
	Long:  `You don´t need a addon on your cluster anymore then easly destroy it.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		addon.DestroyAddon(args[0])
	},
}

// addonCmd represents the addon command
var addonCreateCmd = &cobra.Command{
	Use:   "create [addon name]",
	Short: "create a new kubernetes addon packages on your local machine for development",
	Long: `Create your own addons for your kubernetes cluster. 
This command will prepare a example package in a folder with the addon name`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		addon.PrepareExample(args[0])
	},
}

// addonCmd represents the addon command
var addonGetCmd = &cobra.Command{
	Use:   "get [git url]",
	Short: "Get a addon to your local machine",
	Long:  `If you need to prepare some configs for a addon on your local machine then you can get it with this command in a local folder.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		addon.GetAddon(args[0])
	},
}

func getKubeConfig() string {
	/* This function gets the path to the kubeconfig, cluster details and auth
	   for using with the kubectl.
	   Then use this to install the addon on this cluster
	*/
	fmt.Println("Please enter the path to your kubeconfig:")
	var kubeConfig string
	fmt.Scanln(&kubeConfig)
	fmt.Printf("path: %s\n", kubeConfig)
	if _, err := os.Stat(kubeConfig); err != nil {
		fmt.Println("Kubeconfig file not found, kindly check")
		os.Exit(1)
	}
	return kubeConfig
}

func checkKubectl(kubeConfig string) {
	/*This function is used to check the whether kubectl command is installed &
	  it works with the kubeConfig provided
	*/
	kerr, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatal("kubectl command not found, kindly check")
	}
	fmt.Printf("Found kubectl at %s\n", kerr)
	rr, err := exec.Command("kubectl", "--kubeconfig", kubeConfig, "version", "--short").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(rr))
}

func init() {
	addonCmd.AddCommand(addonCreateCmd)
	addonCmd.AddCommand(addonInstallCmd)
	addonCmd.AddCommand(addonGetCmd)
	addonCmd.AddCommand(addonDestroyCmd)
	rootCmd.AddCommand(addonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// addonCmd.Flags().BoolVarP(&monitor, "monitor", "m", false, "Deploy Monitoring and Alerting")
	addonCmd.PersistentFlags().StringVar(&addon.KubeConfig, "kubeconfig", "", "kubeconfig path")
}
