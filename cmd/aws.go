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
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var install, create, destroy bool

var ec2IP string

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

			//check if ansible host file exists

			if _, err := os.Stat("./kubespray/inventory/hosts"); err != nil {
				fmt.Println("./kubespray/inventory/host inventory file not found")
				os.Exit(1)
			}

			// Copy the configuraton files as indicated in the kubespray docs

			if _, err := os.Stat("./kubespray/inventory/awscluster"); err == nil {
				fmt.Println("Configuration folder already exists")
			} else {
				//os.MkdirAll("./kubespray/inventory/awscluster/group_vars", 0755)
				exec.Command("cp", "-rfp", "./kubespray/inventory/sample/", "./kubespray/inventory/awscluster/").Run()

				exec.Command("cp", "./kubespray/inventory/hosts", "./kubespray/inventory/awscluster/hosts").Run()

				//Enable load balancer api access and copy the kubeconfig file locally
				loadBalancerName, err := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./kubespray/inventory/hosts | cut -d'=' -f2").CombinedOutput()
				if err != nil {
					fmt.Println("Problem getting the load balancer domain name", err)
				} else {
					//Make a copy of kubeconfig on Ansible host
					f, err := os.OpenFile("./kubespray/inventory/awscluster/group_vars/k8s-cluster.yml", os.O_APPEND|os.O_WRONLY, 0600)
					if err != nil {
						panic(err)
					}

					defer f.Close()
					fmt.Fprintf(f, "kubeconfig_localhost: true\n")

					g, err := os.OpenFile("./kubespray/inventory/awscluster/group_vars/all.yml", os.O_APPEND|os.O_WRONLY, 0600)
					if err != nil {
						panic(err)
					}

					defer g.Close()

					// Resolve Load Balancer Domain Name and pick the first IP

					s, _ := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./kubespray/inventory/hosts | cut -d'=' -f2 | sed 's/\"//g'").CombinedOutput()
					// Convert the Domain name to string and strip all spaces so that Lookup does not return errors
					r := string(s)
					t := strings.TrimSpace(r)

					fmt.Println(t)
					node, err := net.LookupHost(t)

					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					ec2IP := node[0]

					fmt.Println(node)

					DomainName := strings.TrimSpace(string(loadBalancerName))
					loadBalancerDomainName := "apiserver_loadbalancer_domain_name: " + DomainName

					fmt.Fprintf(g, "#Set cloud provider to AWS\n")
					fmt.Fprintf(g, "cloud_provider: 'aws'\n")
					fmt.Fprintf(g, "#Load Balancer Configuration\n")
					fmt.Fprintf(g, "loadbalancer_apiserver_localhost: false\n")
					fmt.Fprintf(g, "%s\n", loadBalancerDomainName)
					fmt.Fprintf(g, "loadbalancer_apiserver:\n")
					fmt.Fprintf(g, "  address: %s\n", ec2IP)
					fmt.Fprintf(g, "  port: 6443\n")
				}
			}
			kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/awscluster/hosts", "./cluster.yml", "--timeout=60", "-e ansible_user=core", "-e bootstrap_os=coreos", "-b", "--become-user=root", "--flush-cache")
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
			err = os.RemoveAll("./kubespray/inventory/awscluster")
			if err != nil {
				fmt.Println(err)
			}

			// Check if credentials file exist, if it exists skip asking to input the AWS values
			if _, err := os.Stat("./kubespray/contrib/terraform/aws/credentials.tfvars"); err == nil {
				fmt.Println("Credentials file already exists, creation skipped")
			} else {

				fmt.Println("Please enter your AWS access key ID")
				var awsAccessKeyID string
				fmt.Scanln(&awsAccessKeyID)

				fmt.Println("Please enter your AWS SECRET ACCESS KEY")
				var awsSecretKey string
				fmt.Scanln(&awsSecretKey)

				fmt.Println("Please enter your AWS SSH Key Name")
				var awsAccessSSHKey string
				fmt.Scanln(&awsAccessSSHKey)

				fmt.Println("Please enter your AWS Default Region")
				var awsDefaultRegion string
				fmt.Scanln(&awsDefaultRegion)

				file, err := os.Create("./kubespray/contrib/terraform/aws/credentials.tfvars")
				if err != nil {
					log.Fatal("Cannot create file", err)
				}
				defer file.Close()

				fmt.Fprintf(file, "AWS_ACCESS_KEY_ID = %s\n", awsAccessKeyID)
				fmt.Fprintf(file, "AWS_SECRET_ACCESS_KEY = %s\n", awsSecretKey)
				fmt.Fprintf(file, "AWS_SSH_KEY_NAME = %s\n", awsAccessSSHKey)
				fmt.Fprintf(file, "AWS_DEFAULT_REGION = %s\n", awsDefaultRegion)
			}
			terrSet := exec.Command("terraform", "destroy", "-var-file=credentials.tfvars", "-force")
			terrSet.Dir = "./kubespray/contrib/terraform/aws/"
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

			// Check if credentials file exist, if it exists skip asking to input the AWS values
			if _, err := os.Stat("./kubespray/contrib/terraform/aws/credentials.tfvars"); err == nil {
				fmt.Println("Credentials file already exists, creation skipped")
			} else {

				//Read Configuration File
				viper.SetConfigName("config")

				viper.AddConfigPath(".")
				viper.AddConfigPath("/tk8")
				verr := viper.ReadInConfig() // Find and read the config file
				if verr != nil {             // Handle errors reading the config file
					panic(fmt.Errorf("fatal error config file: %s", verr))
				}

				awsAccessKeyID := viper.GetString("aws.aws_access_key_id")

				awsSecretKey := viper.GetString("aws.aws_secret_access_key")

				awsAccessSSHKey := viper.GetString("aws.aws_ssh_keypair")

				awsDefaultRegion := viper.GetString("aws.aws_default_region")

				file, err := os.Create("./kubespray/contrib/terraform/aws/credentials.tfvars")
				if err != nil {
					log.Fatal("Cannot create file", err)
				}
				defer file.Close()

				fmt.Fprintf(file, "AWS_ACCESS_KEY_ID = %s\n", strconv.Quote(awsAccessKeyID))
				fmt.Fprintf(file, "AWS_SECRET_ACCESS_KEY = %s\n", strconv.Quote(awsSecretKey))
				fmt.Fprintf(file, "AWS_SSH_KEY_NAME = %s\n", strconv.Quote(awsAccessSSHKey))
				fmt.Fprintf(file, "AWS_DEFAULT_REGION = %s\n", strconv.Quote(awsDefaultRegion))

			}
			// Remove tftvars file

			err = os.Remove("./kubespray/contrib/terraform/aws/terraform.tfvars")
			if err != nil {
				fmt.Println(err)
			}

			//Read Configuration File
			viper.SetConfigName("config")

			viper.AddConfigPath(".")
			verr := viper.ReadInConfig() // Find and read the config file
			if verr != nil {             // Handle errors reading the config file
				panic(fmt.Errorf("fatal error config file: %s", verr))
			}

			awsClusterName := viper.GetString("aws.clustername")
			awsVpcCidrBlock := viper.GetString("aws.aws_vpc_cidr_block")
			awsCidrSubnetsPrivate := viper.GetString("aws.aws_cidr_subnets_private")
			awsCidrSubnetsPublic := viper.GetString("aws.aws_cidr_subnets_public")
			awsBastionSize := viper.GetString("aws.aws_bastion_size")
			awsKubeMasterNum := viper.GetString("aws.aws_kube_master_num")
			awsKubeMasterSize := viper.GetString("aws.aws_kube_master_size")
			awsEtcdNum := viper.GetString("aws.aws_etcd_num")
			awsEtcdSize := viper.GetString("aws.aws_etcd_size")
			awsKubeWorkerNum := viper.GetString("aws.aws_kube_worker_num")
			awsKubeWorkerSize := viper.GetString("aws.aws_kube_worker_size")
			awsElbAPIPort := viper.GetString("aws.aws_elb_api_port")
			k8sSecureAPIPort := viper.GetString("aws.k8s_secure_api_port")
			kubeInsecureApiserverAddress := viper.GetString("aws.")

			tfile, err := os.Create("./kubespray/contrib/terraform/aws/terraform.tfvars")
			if err != nil {
				log.Fatal("Cannot create file", err)
			}
			defer tfile.Close()

			fmt.Fprintf(tfile, "aws_cluster_name = %s\n", strconv.Quote(awsClusterName))
			fmt.Fprintf(tfile, "aws_vpc_cidr_block = %s\n", strconv.Quote(awsVpcCidrBlock))
			fmt.Fprintf(tfile, "aws_cidr_subnets_private = %s\n", awsCidrSubnetsPrivate)
			fmt.Fprintf(tfile, "aws_cidr_subnets_public = %s\n", awsCidrSubnetsPublic)

			fmt.Fprintf(tfile, "aws_bastion_size = %s\n", strconv.Quote(awsBastionSize))
			fmt.Fprintf(tfile, "aws_kube_master_num = %s\n", awsKubeMasterNum)
			fmt.Fprintf(tfile, "aws_kube_master_size = %s\n", strconv.Quote(awsKubeMasterSize))
			fmt.Fprintf(tfile, "aws_etcd_num = %s\n", awsEtcdNum)

			fmt.Fprintf(tfile, "aws_etcd_size = %s\n", strconv.Quote(awsEtcdSize))
			fmt.Fprintf(tfile, "aws_kube_worker_num = %s\n", awsKubeWorkerNum)
			fmt.Fprintf(tfile, "aws_kube_worker_size = %s\n", strconv.Quote(awsKubeWorkerSize))
			fmt.Fprintf(tfile, "aws_elb_api_port = %s\n", awsElbAPIPort)
			fmt.Fprintf(tfile, "k8s_secure_api_port = %s\n", k8sSecureAPIPort)
			fmt.Fprintf(tfile, "kube_insecure_apiserver_address = %s\n", strconv.Quote(kubeInsecureApiserverAddress))

			fmt.Fprintf(tfile, "default_tags = {\n")
			fmt.Fprintf(tfile, "#  Env = 'devtest'\n")
			fmt.Fprintf(tfile, "#  Product = 'kubernetes'\n")
			fmt.Fprintf(tfile, "}")

			//fmt.Println("Please enter your AWS access key ID")
			//var awsAccessKeyID string
			//fmt.Scanln(&awsAccessKeyID)

			//fmt.Println("Please enter your AWS SECRET ACCESS KEY")
			//var awsSecretKey string
			//fmt.Scanln(&awsSecretKey)

			//fmt.Println("Please enter your AWS SSH Key Name")
			//var awsAccessSSHKey string
			//fmt.Scanln(&awsAccessSSHKey)

			//fmt.Println("Please enter your AWS Default Region")
			//var awsDefaultRegion string
			//fmt.Scanln(&awsDefaultRegion)

			terrInit := exec.Command("terraform", "init")
			terrInit.Dir = "./kubespray/contrib/terraform/aws/"
			out, _ := terrInit.StdoutPipe()
			terrInit.Start()
			scanInit := bufio.NewScanner(out)
			for scanInit.Scan() {
				m := scanInit.Text()
				fmt.Println(m)
				//log.Printf(m)
			}

			terrInit.Wait()

			terrSet := exec.Command("terraform", "apply", "-var-file=credentials.tfvars", "-auto-approve")
			terrSet.Dir = "./kubespray/contrib/terraform/aws/"
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

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
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
