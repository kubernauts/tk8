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
	"github.com/spf13/viper"
)

var load bool

// openstackCmd represents the openstack command
var bareCmd = &cobra.Command{
	Use:   "baremetal",
	Short: "Manages the infrastructure on Baremetal",
	Long: `
Create and delete kubernetes deployment that is running on Baremetal.`,
	Run: func(cmd *cobra.Command, args []string) {

		if install {

			// check if ansible is installed
			terr, err := exec.LookPath("ansible")
			if err != nil {
				log.Fatal("Ansible command not found, kindly check")
			}
			fmt.Printf("Found Ansible at %s\n", terr)
			rr, err := exec.Command("ansible", "--version").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))

			//Start Kubernetes Installation

			//check if kubespray cluster folder for the (baremetal)exists

			// Copy sample-inventory as indicated in the kubespray docs

			if _, err := os.Stat("./kubespray/inventory/barecluster"); err == nil {
				fmt.Println("Inventory folder already exists")
			} else {
				exec.Command("cp", "-LRp", "./kubespray/inventory/sample", "./kubespray/inventory/barecluster").Run()

				//Make a copy of kubeconfig on Ansible host after installation
				f, err := os.OpenFile("./kubespray/inventory/barecluster/group_vars/k8s-cluster.yml", os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					panic(err)
				}

				defer f.Close()
				fmt.Fprintf(f, "kubeconfig_localhost: true\n")

			}

			//Check for the host.ini file if it exists

			if _, err := os.Stat("./baremetal/hosts.ini"); err != nil {
				fmt.Println("Kubespray hosts.ini does not exist in tk8/baremetal folder, please check")
				os.Exit(1)
			}

			// Copy the hosts.ini file from the baremetal folder to the inventory folder
			exec.Command("cp", "-rf", "./baremetal/hosts.ini", "./kubespray/inventory/barecluster/hosts.ini").Run()

			//Get the value of Operating system username and whether to become root user
			viper.SetConfigName("variables")

			viper.AddConfigPath(".")
			viper.AddConfigPath("./baremetal/")
			venv := viper.ReadInConfig() // Find and read the config file
			if venv != nil {             // Handle errors reading the config file
				panic(fmt.Errorf("Fatal error config file: %s \n", venv))
			}

			OsUser := viper.GetString("os.username")
			AnsibleUser := "-e ansible_user=" + OsUser
			AnsibleBecome := viper.GetString("os.become")

			fmt.Print(OsUser, AnsibleBecome, AnsibleUser)

			if AnsibleBecome == "yes" {
				kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./cluster.yml", AnsibleUser, "-b", "--become-user=root", "--flush-cache")
				kubeSet.Dir = "./kubespray/"
				stdout, _ := kubeSet.StdoutPipe()
				kubeSet.Stderr = kubeSet.Stdout
				kubeSet.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				kubeSet.Wait()

				os.Exit(0)
			}

			if AnsibleBecome == "no" {
				kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./cluster.yml", AnsibleUser, "--flush-cache")
				kubeSet.Dir = "./kubespray/"
				stdout, _ := kubeSet.StdoutPipe()
				kubeSet.Stderr = kubeSet.Stdout
				kubeSet.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				kubeSet.Wait()

				os.Exit(0)
			}
		}

		if destroy {

			// check if ansible is installed
			terr, err := exec.LookPath("ansible")
			if err != nil {
				log.Fatal("Ansible command not found, kindly check")
			}
			fmt.Printf("Found Ansible at %s\n", terr)
			rr, err := exec.Command("ansible", "--version").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))

			//Check for the host.ini file if it exists

			if _, err := os.Stat("./baremetal/hosts.ini"); err != nil {
				fmt.Println("Kubespray hosts.ini does not exist in tk8/baremetal folder, please check")
				os.Exit(1)
			}

			//Get the value of Operating system username and whether to become root user
			viper.SetConfigName("variables")

			viper.AddConfigPath(".")
			viper.AddConfigPath("./baremetal/")
			venv := viper.ReadInConfig() // Find and read the config file
			if venv != nil {             // Handle errors reading the config file
				panic(fmt.Errorf("Fatal error config file: %s \n", venv))
			}

			OsUser := viper.GetString("os.username")
			AnsibleUser := "-e ansible_user=" + OsUser
			AnsibleBecome := viper.GetString("os.become")

			if AnsibleBecome == "yes" {
				kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./reset.yml", AnsibleUser, "-b", "--become-user=root", "--flush-cache")
				kubeSet.Dir = "./kubespray/"
				stdout, _ := kubeSet.StdoutPipe()
				kubeSet.Stderr = kubeSet.Stdout
				kubeSet.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				kubeSet.Wait()

				// Remove the cluster inventory folder
				err = os.RemoveAll("./kubespray/inventory/barecluster")
				if err != nil {
					fmt.Println(err)
				}

				os.Exit(0)

			}

			if AnsibleBecome == "no" {
				kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/barecluster/hosts.ini", "./reset.yml", AnsibleUser, "--flush-cache")
				kubeSet.Dir = "./kubespray/"
				stdout, _ := kubeSet.StdoutPipe()
				kubeSet.Stderr = kubeSet.Stdout
				kubeSet.Start()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
					//log.Printf(m)
				}

				kubeSet.Wait()

				// Remove the cluster inventory folder
				err = os.RemoveAll("./kubespray/inventory/barecluster")
				if err != nil {
					fmt.Println(err)
				}

				os.Exit(0)
			}

		}

		if load {

			// Check for the load balancer configmap file in /tk8/baremetal

			if _, err := os.Stat("./baremetal/lb-config.yml"); err != nil {
				fmt.Println("The baremetal configmap does not exist, please check")
				os.Exit(1)
			}

			// Check for the kubeconfig file in /tk8/baremetal

			if _, err := os.Stat("./baremetal/kubeconfig"); err != nil {
				fmt.Println("The baremetal configmap does not exist, please check")
				os.Exit(1)
			}

			// Get kubeconfig file location
			fmt.Println("Please enter the path to your kubeconfig")
			var kubeConfig string
			fmt.Scanln(&kubeConfig)

			if _, err := os.Stat(kubeConfig); err != nil {
				fmt.Println("Kubeconfig not found, kindly check")
				os.Exit(1)
			}

			// check if kubectl is installed
			terr, err := exec.LookPath("kubectl")
			if err != nil {
				log.Fatal("Kubectl command not found, kindly check")
			}
			fmt.Printf("Found kubectl at %s\n", terr)
			rr, err := exec.Command("kubectl", "version", "--client=true").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))

			// Deploy MetallB Artifacts i.e deployment, configmap, etc

			// Deploy the Controller, speaker, RBAC, service account
			kubeSet := exec.Command("kubectl", "--kubeconfig", kubeConfig, "create", "-f", "https://raw.githubusercontent.com/google/metallb/v0.5.0/manifests/metallb.yaml")
			kubeSet.Dir = "./baremetal/"
			stdout, _ := kubeSet.StdoutPipe()
			kubeSet.Stderr = kubeSet.Stdout
			kubeSet.Start()
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			kubeSet.Wait()

			// Deploy the Configmap
			ConfigSet := exec.Command("kubectl", "--kubeconfig", kubeConfig, "create", "-f", "lb-config.yml")
			ConfigSet.Dir = "./baremetal/"
			Configout, _ := ConfigSet.StdoutPipe()
			ConfigSet.Stderr = ConfigSet.Stdout
			ConfigSet.Start()
			Configscanner := bufio.NewScanner(Configout)
			for Configscanner.Scan() {
				m := Configscanner.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			ConfigSet.Wait()

			// Print out the objects under metallb namespace
			PrintSet := exec.Command("kubectl", "--kubeconfig", kubeConfig, "-n", "metallb-system", "get", "all")
			PrintSet.Dir = "./baremetal/"
			Printout, _ := PrintSet.StdoutPipe()
			PrintSet.Stderr = ConfigSet.Stdout
			PrintSet.Start()
			Objectscanner := bufio.NewScanner(Printout)
			for Objectscanner.Scan() {
				m := Objectscanner.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			PrintSet.Wait()

			os.Exit(0)

		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

	},
}

func init() {
	clusterCmd.AddCommand(bareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	bareCmd.Flags().BoolVarP(&load, "loadbalancer", "l", false, "Deploy the load balancer")

	bareCmd.Flags().BoolVarP(&destroy, "destroy", "d", false, "Destroy the baremetal kubernetes")

	bareCmd.Flags().BoolVarP(&install, "install", "i", false, "Install kuberntes on the baremetal infrastructure")
}
