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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// openstackCmd represents the openstack command
var openstackCmd = &cobra.Command{
	Use:   "openstack",
	Short: "Manages the infrastructure on Openstack",
	Long: `
	Create, delete and show current status of the deployment that is running on Openstack.
	Kindly ensure that terraform is installed also.`,
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

			//check if ansible Cluster folder exists

			// Copy sample-inventory as indicated in the kubespray docs

			if _, err := os.Stat("./kubespray/inventory/stackcluster"); err == nil {
				fmt.Println("Inventory folder already exists")
			} else {
				exec.Command("cp", "-LRp", "./kubespray/inventory/sample", "./kubespray/inventory/stackcluster").Run()

				// Copyt the hosts file after the infrastructre has been deployed by terraform

				exec.Command("cp", "-rf", "./kubespray/contrib/terraform/openstack/hosts.ini", "./kubespray/inventory/stackcluster/hosts.ini").Run()

				//Enable load balancer api access and copy the kubeconfig file locally, // Get Load Balancer IP and input to all.yml
				loadBalancerName, err := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./kubespray/inventory/stackcluster/hosts.ini | cut -d'=' -f2").CombinedOutput()
				if err != nil {
					fmt.Println("Problem getting the load balancer domain name", err)
				}

				DomainName := strings.TrimSpace(string(loadBalancerName))

				loadBalancerDomainName := "apiserver_loadbalancer_domain_name: " + DomainName

				//Read Configuration File
				viper.SetConfigName("network-config")

				viper.AddConfigPath(".")
				viper.AddConfigPath("./kubespray/contrib/terraform/openstack/")
				verr := viper.ReadInConfig() // Find and read the config file
				if verr != nil {             // Handle errors reading the config file
					panic(fmt.Errorf("Fatal error config file: %s \n", verr))
				}

				LBIP := viper.GetString("floating-master-lb-vip")

				LBaasSubnetID := viper.GetString("lbaas-private-subnet-id")

				LBaaSFloatNetworkID := viper.GetString("lbaas-floating-network-id")

				g, err := os.OpenFile("./kubespray/inventory/stackcluster/group_vars/all.yml", os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					panic(err)
				}

				defer g.Close()

				fmt.Fprintf(g, "#Set cloud provider to Openstack\n")
				fmt.Fprintf(g, "cloud_provider: 'openstack'\n")
				fmt.Fprintf(g, "openstack_lbaas_enabled: True\n")
				fmt.Fprintf(g, "openstack_lbaas_subnet_id: %s\n", strconv.Quote(LBaasSubnetID))
				fmt.Fprintf(g, "openstack_lbaas_floating_network_id: %s\n", strconv.Quote(LBaaSFloatNetworkID))
				fmt.Fprintf(g, "#Load Balancer Configuration\n")
				fmt.Fprintf(g, "loadbalancer_apiserver_localhost: false\n")
				fmt.Fprintf(g, "%s\n", loadBalancerDomainName)
				fmt.Fprintf(g, "loadbalancer_apiserver:\n")
				fmt.Fprintf(g, "  address: %s\n", LBIP)
				fmt.Fprintf(g, "  port: 6443\n")

				//Make a copy of kubeconfig on localhost, still buggy, so I disabled this pending resolution
				//f, err := os.OpenFile("./kubespray/inventory/stackcluster/group_vars/k8s-cluster.yml", os.O_APPEND|os.O_WRONLY, 0600)
				//if err != nil {
				//	panic(err)
				//}

				//defer f.Close()
				//fmt.Fprintf(f, "kubectl_localhost: true\n")

				//Set Network Plugin to Flannel

				exec.Command("sh", "-c", "sed -i 's/kube_network_plugin: calico/kube_network_plugin: flannel/g' ./kubespray/inventory/stackcluster/group_vars/k8s-cluster.yml").Run()
			}

			//Get credentials from clouds.yml and export as environment varaibles
			//Read Configuration File
			viper.SetConfigName("clouds")

			viper.AddConfigPath(".")
			viper.AddConfigPath("./kubespray/contrib/terraform/openstack/")
			venv := viper.ReadInConfig() // Find and read the config file
			if venv != nil {             // Handle errors reading the config file
				panic(fmt.Errorf("Fatal error config file: %s \n", venv))
			}

			OsAuthURL := viper.GetString("clouds.mycloud.auth.auth_url")

			OsProjectDomainName := viper.GetString("clouds.mycloud.auth.user_domain_name")

			OsUserDomainName := viper.GetString("clouds.mycloud.auth.user_domain_name")

			OsProjectName := viper.GetString("clouds.mycloud.auth.project_name")

			OsTenantName := viper.GetString("clouds.mycloud.auth.project_name")

			OsTenantID := viper.GetString("clouds.mycloud.auth.tenant_id")

			OsUserName := viper.GetString("clouds.mycloud.auth.username")

			OsPassword := viper.GetString("clouds.mycloud.auth.password")

			OsRegionName := viper.GetString("clouds.mycloud.region_name")

			OsInterface := viper.GetString("clouds.mycloud.interface")

			OsIdentityAPIVersion := viper.GetString("clouds.mycloud.identity_api_version")

			//Export Openstack Credentials to the environment

			os.Setenv("OS_AUTH_URL", OsAuthURL)

			os.Setenv("OS_PROJECT_DOMAIN_NAME", OsProjectDomainName)

			os.Setenv("OS_USER_DOMAIN_NAME", OsUserDomainName)

			os.Setenv("OS_PROJECT_NAME", OsProjectName)

			os.Setenv("OS_TENANT_NAME", OsTenantName)

			os.Setenv("OS_TENANT_ID", OsTenantID)

			os.Setenv("OS_USERNAME", OsUserName)

			os.Setenv("OS_PASSWORD", OsPassword)

			os.Setenv("OS_REGION_NAME", OsRegionName)

			os.Setenv("OS_INTERFACE", OsInterface)

			os.Setenv("OS_IDENTITY_API_VERSION", OsIdentityAPIVersion)

			kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/stackcluster/hosts.ini", "./cluster.yml", "-e ansible_user=centos", "-b", "--become-user=root", "--flush-cache")
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

		if create {
			// check if terraform is installed
			terr, err := exec.LookPath("terraform")
			if err != nil {
				log.Fatal("Terraform command not found, kindly check")
			}
			fmt.Printf("Found terraform at %s\n", terr)
			rr, err := exec.Command("terraform", "version").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))

			// Overwrite the terraform openstack folder, added the LB module so as to deploy the k8s master(s) behind LBaaS

			exec.Command("cp", "-rf", "./openstack/", "./kubespray/contrib/terraform/").Run()

			// Check the Openstack credentials file, if found copy to the stackcluster directory as clouds.yaml
			fmt.Println("Checking if openstack credential file exists")
			if _, err := os.Stat("./openstack/stack_credentials.yaml"); err != nil {
				fmt.Println("./openstack/stack_credentials.yaml credential file not found, please check")
				os.Exit(1)
			} else {
				fmt.Println("Credentials file exists, copying to the stackcluster directory")
				exec.Command("cp", "-rfp", "./openstack/stack_credentials.yaml", "./kubespray/contrib/terraform/openstack/clouds.yaml").Run()

			}

			// Check if Openstack cluster specific file exists (stack_cluster.tf), if found copy to the stackcluster directory as cluster.tf
			fmt.Println("Checking if cluster config file exists")
			if _, err := os.Stat("./openstack/cluster.tfvars"); err != nil {
				fmt.Println("./openstack/stack_cluster.tf configuration file not found, please check")
				os.Exit(1)
			} else {
				fmt.Println("Cluster config file exists, copying to the terraform openstack directory")
				exec.Command("cp", "-rfp", "./openstack/cluster.tfvars", "./kubespray/contrib/terraform/openstack/cluster.tfvars").Run()

			}

			// Terraform Initialization and create the infrastructure

			terrInit := exec.Command("terraform", "init")
			terrInit.Dir = "./kubespray/contrib/terraform/openstack/"
			out, _ := terrInit.StdoutPipe()
			terrInit.Start()
			scanInit := bufio.NewScanner(out)
			for scanInit.Scan() {
				m := scanInit.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			terrInit.Wait()

			fmt.Println("Starting terraform apply")
			terrSet := exec.Command("terraform", "apply", "-var-file=cluster.tfvars", "-auto-approve")
			terrSet.Dir = "./kubespray/contrib/terraform/openstack/"
			stdout, err := terrSet.StdoutPipe()
			terrSet.Stderr = terrSet.Stdout
			terrSet.Start()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			terrSet.Wait()

			os.Exit(0)

		}

		if destroy {

			// check if terraform is installed
			terr, err := exec.LookPath("terraform")
			if err != nil {
				log.Fatal("Terraform command not found, kindly check")
			}
			fmt.Printf("Found terraform at %s\n", terr)
			rr, err := exec.Command("terraform", "version").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(rr))

			// Remove ssh bastion file

			if _, err := os.Stat("./kubespray/ssh-bastion.conf"); err == nil {
				os.Remove("./kubespray/ssh-bastion.conf")
			}

			// Remove the cluster inventory folder
			err = os.RemoveAll("./kubespray/inventory/stackcluster")
			if err != nil {
				fmt.Println(err)
			}

			// Check if credentials file exist, if it exists skip asking to input the AWS values
			if _, err := os.Stat("./kubespray/contrib/terraform/openstack/clouds.yaml"); err == nil {
				fmt.Println("Credentials file already exists, creation skipped")
			} else {

				fmt.Println("Checking if openstack credential file exists")
				if _, err := os.Stat("./stack_credentials.yaml"); err != nil {
					fmt.Println("./tk8/stack_credentials.yaml credential file not found, please check")
					os.Exit(1)
				} else {
					fmt.Println("Credentials file exists, copying to the stackcluster directory")
					exec.Command("cp", "-rfp", "./stack_credentials.yaml", "./kubespray/contrib/terraform/openstack/clouds.yaml").Run()

				}
			}
			terrSet := exec.Command("terraform", "destroy", "-var-file=cluster.tfvars", "-force")
			terrSet.Dir = "./kubespray/contrib/terraform/openstack/"
			stdout, _ := terrSet.StdoutPipe()
			terrSet.Stderr = terrSet.Stdout
			error := terrSet.Start()
			if error != nil {
				fmt.Println(error)
			}
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			terrSet.Wait()

			os.Exit(0)
		}

	},
}

func init() {
	clusterCmd.AddCommand(openstackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openstackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openstackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	openstackCmd.Flags().BoolVarP(&create, "create", "c", false, "Deploy the Openstack infrastructure using terraform")

	openstackCmd.Flags().BoolVarP(&destroy, "destroy", "d", false, "Destroy the infrastructure using terraform")

	openstackCmd.Flags().BoolVarP(&install, "install", "i", false, "Install kuberntes on the Openstack infrastructure")
}
