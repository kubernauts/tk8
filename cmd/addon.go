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
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var monitor, heapster, rancher bool

// addonCmd represents the addon command
var addonCmd = &cobra.Command{
	Use:   "addon",
	Short: "Install kubernetes addon packages",
	Long: `
Install additional packages on top of your kubernetes deployment. Examples: Prometheus,
Zipkin, Kibana, Load Testing As A Service`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if monitor {
			// Get kubeconfig file location
			fmt.Println("Please enter the path to your kubeconfig")
			var kubeConfig string
			fmt.Scanln(&kubeConfig)

			if _, err := os.Stat(kubeConfig); err != nil {
				fmt.Println("Kubeconfig not found, kindly check")
				os.Exit(1)
			}

			// check if kubectl is installed
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

			if _, err := os.Stat("./prometheus-grafana-alerting"); err == nil {
				fmt.Println("Addon Files already exist on this system ... skip")
			} else {

				err := exec.Command("git", "clone", "https://github.com/arashkaffamanesh/prometheus-grafana-alerting").Run()
				if err != nil {
					log.Fatalf("Seems there is a problem cloning the Addon Files repo, %v", err)
					os.Exit(1)
				}
			}

			cmdSet := exec.Command("./build.sh")
			cmdSet.Dir = "./prometheus-grafana-alerting"
			cmdOut, er := cmdSet.CombinedOutput()
			if er != nil {
				log.Fatal(er)
			}
			cmdSet.Wait()
			fmt.Println(cmdOut)
			cmdSet = exec.Command("kubectl", "--kubeconfig", kubeConfig, "apply", "-f", "./prometheus-grafana-alerting/manifests-all.yaml")
			cmdOut, er = cmdSet.CombinedOutput()
			if er != nil {
				log.Fatal(er)
			}
			cmdSet.Wait()
			fmt.Println(cmdOut)
			os.Exit(0)
		}

		if rancher {
			/* This is to install the Rancher addon where all the k8s objects
			   for this are provided with main.yml
			   This is applied with the kubectl create -f command
			*/
			//get the kubeconfig file full path
			var kubeConfig = getKubeConfig()
			//Check if kubectl cmd is installed
			checkKubectl(kubeConfig)
			//Check if yaml for Rancher is present relative to current directory
			pwd, _ := os.Getwd()
			_, err1 := os.Stat(pwd + "/addons/rancher/master.yaml")
			if err1 != nil {
				print(err1.Error())
			} else {
				fmt.Println("Deploying Rancher")
				RancherDeploy := exec.Command("kubectl", "--kubeconfig", kubeConfig, "create", "-f", pwd+"/addons/rancher/master.yaml")
				stdout, _ := RancherDeploy.StdoutPipe()
				RancherDeploy.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
				}
				RancherDeploy.Wait()
			}
		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

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
	rootCmd.AddCommand(addonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addonCmd.Flags().BoolVarP(&monitor, "monitor", "l", false, "Deploy Monitoring and Alerting")
	addonCmd.Flags().BoolVarP(&heapster, "heapster", "m", false, "Deploy Heapster")
	addonCmd.Flags().BoolVarP(&rancher, "rancher", "r", false, "Deploy Rancher")
}
