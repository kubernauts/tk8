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

package cluster

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kubernauts/tk8-provisioner-aws/internal/templates"
)

var ec2IP string

func distSelect() (string, string) {
	//Read Configuration File
	awsAmiID, awsInstanceOS, sshUser := GetDistConfig()

	if awsAmiID != "" && sshUser == "" {
		log.Fatal("SSH Username is required when using custom AMI")
		return "", ""
	}
	if awsAmiID == "" && awsInstanceOS == "" {
		log.Fatal("Provide either of AMI ID or OS in the config file.")
		return "", ""
	}

	if awsAmiID != "" && sshUser != "" {
		awsInstanceOS = "custom"
		DistOSMap["custom"] = DistOS{
			User:     sshUser,
			AmiOwner: awsAmiID,
			OS:       "custom",
		}
	}

	return DistOSMap[awsInstanceOS].User, awsInstanceOS
}

func prepareConfigFiles(awsInstanceOS string) {
	if awsInstanceOS == "custom" {
		ParseTemplate(templates.CustomInfrastructure, "./inventory/"+Name+"/provisioner/create-infrastructure.tf", DistOSMap[awsInstanceOS])
	} else {
		ParseTemplate(templates.Infrastructure, "./inventory/"+Name+"/provisioner/create-infrastructure.tf", DistOSMap[awsInstanceOS])
	}

	ParseTemplate(templates.Credentials, "./inventory/"+Name+"/provisioner/credentials.tfvars", GetCredentials())
	ParseTemplate(templates.Variables, "./inventory/"+Name+"/provisioner/variables.tf", DistOSMap[awsInstanceOS])
	ParseTemplate(templates.Terraform, "./inventory/"+Name+"/provisioner/terraform.tfvars", GetClusterConfig())
}

func prepareInventoryGroupAllFile(fileName string) *os.File {
	groupVars, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	ErrorCheck("Error while trying to open "+fileName+": %v.", err)
	return groupVars
}

func prepareInventoryClusterFile(fileName string) *os.File {
	k8sClusterFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	defer k8sClusterFile.Close()
	ErrorCheck("Error while trying to open "+fileName+": %v.", err)
	fmt.Fprintf(k8sClusterFile, "kubeconfig_localhost: true\n")
	return k8sClusterFile
}

// AWSCreate is used to create a infrastructure on AWS.
func AWSCreate() {

	if _, err := os.Stat("./inventory/" + Name + "/provisioner/.terraform"); err == nil {
		fmt.Println("Configuration folder already exists")
	} else {
		sshUser, osLabel := distSelect()
		fmt.Printf("Prepairing Setup for user %s on %s\n", sshUser, osLabel)
		os.MkdirAll("./inventory/"+Name+"/provisioner", 0755)
		err := exec.Command("cp", "-rfp", "./kubespray/contrib/terraform/aws/.", "./inventory/"+Name+"/provisioner").Run()
		ErrorCheck("provisioner could not provided: %v", err)
		prepareConfigFiles(osLabel)
		ExecuteTerraform("init", "./inventory/"+Name+"/provisioner/")
	}

	ExecuteTerraform("apply", "./inventory/"+Name+"/provisioner/")

	// waiting for Loadbalancer and other not completed stuff
	fmt.Println("Infrastructure is upcoming.")
	time.Sleep(15 * time.Second)
	return

}

// AWSInstall is used for installing Kubernetes on the available infrastructure.
func AWSInstall() {
	// check if ansible is installed
	DependencyCheck("ansible")

	// Copy the configuraton files as indicated in the kubespray docs
	if _, err := os.Stat("./inventory/" + Name + "/installer"); err == nil {
		fmt.Println("Configuration folder already exists")
	} else {
		os.MkdirAll("./inventory/"+Name+"/installer", 0755)
		mvHost := exec.Command("mv", "./inventory/hosts", "./inventory/"+Name+"/hosts")
		mvHost.Run()
		mvHost.Wait()
		mvShhBastion := exec.Command("cp", "./kubespray/ssh-bastion.conf", "./inventory/"+Name+"/ssh-bastion.conf")
		mvShhBastion.Run()
		mvShhBastion.Wait()
		//os.MkdirAll("./inventory/"+Name+"/installer/group_vars", 0755)
		cpSample := exec.Command("cp", "-rfp", "./kubespray/inventory/sample/.", "./inventory/"+Name+"/installer/")
		cpSample.Run()
		cpSample.Wait()

		cpKube := exec.Command("cp", "-rfp", "./kubespray/.", "./inventory/"+Name+"/installer/")
		cpKube.Run()
		cpKube.Wait()

		mvInstallerHosts := exec.Command("cp", "./inventory/"+Name+"/hosts", "./inventory/"+Name+"/installer/hosts")
		mvInstallerHosts.Run()
		mvInstallerHosts.Wait()
		mvProvisionerHosts := exec.Command("cp", "./inventory/"+Name+"/hosts", "./inventory/"+Name+"/installer/hosts")
		mvProvisionerHosts.Run()
		mvProvisionerHosts.Wait()

		// Check if Kubeadm is enabled
		EnableKubeadm()

		//Start Kubernetes Installation
		//Enable load balancer api access and copy the kubeconfig file locally
		loadBalancerName, err := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./inventory/"+Name+"/installer/hosts | cut -d'=' -f2").CombinedOutput()

		if err != nil {
			fmt.Println("Problem getting the load balancer domain name", err)
		} else {
			var groupVars *os.File
			//Make a copy of kubeconfig on Ansible host
			if kubesprayVersion == "develop" {
				// Set Kube Network Proxy
				SetNetworkPlugin("./inventory/" + Name + "/installer/group_vars/k8s-cluster")
				prepareInventoryClusterFile("./inventory/" + Name + "/installer/group_vars/k8s-cluster/k8s-cluster.yml")
				groupVars = prepareInventoryGroupAllFile("./inventory/" + Name + "/installer/group_vars/all/all.yml")
			} else {
				// Set Kube Network Proxy
				SetNetworkPlugin("./inventory/" + Name + "/installer/group_vars")
				prepareInventoryClusterFile("./inventory/" + Name + "/installer/group_vars/k8s-cluster.yml")
				groupVars = prepareInventoryGroupAllFile("./inventory/" + Name + "/installer/group_vars/all.yml")
			}
			defer groupVars.Close()
			// Resolve Load Balancer Domain Name and pick the first IP

			elbNameRaw, _ := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./inventory/"+Name+"/installer/hosts | cut -d'=' -f2 | sed 's/\"//g'").CombinedOutput()

			// Convert the Domain name to string, strip all spaces so that Lookup does not return errors
			elbName := strings.TrimSpace(string(elbNameRaw))
			fmt.Println(elbName)
			node, err := net.LookupHost(elbName)
			ErrorCheck("Error resolving ELB name: %v", err)
			elbIP := node[0]
			fmt.Println(node)

			DomainName := strings.TrimSpace(string(loadBalancerName))
			loadBalancerDomainName := "apiserver_loadbalancer_domain_name: " + DomainName

			fmt.Fprintf(groupVars, "#Set cloud provider to AWS\n")
			fmt.Fprintf(groupVars, "cloud_provider: 'aws'\n")
			fmt.Fprintf(groupVars, "#Load Balancer Configuration\n")
			fmt.Fprintf(groupVars, "loadbalancer_apiserver_localhost: false\n")
			fmt.Fprintf(groupVars, "%s\n", loadBalancerDomainName)
			fmt.Fprintf(groupVars, "loadbalancer_apiserver:\n")
			fmt.Fprintf(groupVars, "  address: %s\n", elbIP)
			fmt.Fprintf(groupVars, "  port: 6443\n")
		}
	}

	RunPlaybook("./inventory/"+Name+"/installer/", "cluster.yml")

	return
}

// AWSDestroy is used to destroy the infrastructure created.
func AWSDestroy() {
	// Check if credentials file exist, if it exists skip asking to input the AWS values
	if _, err := os.Stat("./inventory/" + Name + "/provisioner/credentials.tfvars"); err == nil {
		fmt.Println("Credentials file already exists, creation skipped")
	} else {

		ParseTemplate(templates.Credentials, "./inventory/"+Name+"/provisioner/credentials.tfvars", GetCredentials())
	}
	cpHost := exec.Command("cp", "./inventory/"+Name+"/hosts", "./inventory/hosts")
	cpHost.Run()
	cpHost.Wait()

	ExecuteTerraform("destroy", "./inventory/"+Name+"/provisioner/")

	exec.Command("rm", "./inventory/hosts").Run()
	exec.Command("rm", "-rf", "./inventory/"+Name).Run()

	return
}

// AWSScale is used to scale the AWS infrastructure and Kubernetes
func AWSScale() {
	var confirmation string
	// Scale the AWS infrastructure
	fmt.Printf("\t\t===============Starting AWS Scaling====================\n\n")

	// Scale the Kubernetes cluster
	fmt.Printf("\n\n\t\t===============Starting Kubernetes Scaling====================\n\n")
	_, err := os.Stat("./inventory/" + Name + "/provisioner/hosts")
	ErrorCheck("No host file found.", err)
	fmt.Printf("\n\nThis will overwrite the previous host file with a new one. Type \"yes\" to confirm:\n")
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		fmt.Printf("Confirmation denied. Exiting...")
		os.Exit(0)
	}
	ExecuteTerraform("apply", "./inventory/"+Name+"/provisioner/")
	mvHost := exec.Command("mv", "./inventory/hosts", "./inventory/"+Name+"/hosts")
	mvHost.Run()
	mvHost.Wait()
	RunPlaybook("./inventory/"+Name+"/installer/", "scale.yml")

	return
}

// AWSReset is used to reset the AWS infrastructure and removing Kubernetes from it.
func AWSReset() {
	RunPlaybook("./inventory/"+Name+"/installer/", "reset.yml")
	return
}

func AWSRemove() {
	NotImplemented()
}
