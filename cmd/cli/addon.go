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
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kubernauts/tk8/internal/addon"
	"github.com/spf13/cobra"
)

var Addon addon.Addon

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
		Addon.Install(args[0])
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
		Addon.Destroy(args[0])
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
		Addon.Create(args[0])
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
		Addon.Get(args[0])
	},
}

/*
 * This function gets the path to the kubeconfig, cluster details and auth
 * for using with the kubectl.
 * Then use this to install the addon on this cluster
 */
func getKubeConfig() string {
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

/*
 * This function is used to check the whether kubectl command is installed &
 * it works with the kubeConfig provided
 */
func checkKubectl(kubeConfig string) {
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
	rootCmd.AddCommand(addonCmd)

	addonCmd.AddCommand(addonInstallCmd)
	addonCmd.AddCommand(addonCreateCmd)
	addonCmd.AddCommand(addonGetCmd)
	addonCmd.AddCommand(addonDestroyCmd)

	addonCmd.PersistentFlags().StringVar(&addon.KubeConfig, "kubeconfig", "", "kubeconfig path")
}
