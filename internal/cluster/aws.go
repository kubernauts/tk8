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
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/kubernauts/tk8/internal/templates"
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
	if awsInstanceOS == "custom" {
		go ParseTemplate(templates.CustomInfrastructure, "./inventory/"+Name+"/provisioner/create-infrastructure.tf", DistOSMap[awsInstanceOS])
	} else {
		go ParseTemplate(templates.Infrastructure, "./inventory/"+Name+"/provisioner/create-infrastructure.tf", DistOSMap[awsInstanceOS])
	}

	go ParseTemplate(templates.Credentials, "./inventory/"+Name+"/provisioner/credentials.tfvars", GetCredentials())
	go ParseTemplate(templates.Variables, "./inventory/"+Name+"/provisioner/variables.tf", DistOSMap[awsInstanceOS])
	go ParseTemplate(templates.Terraform, "./inventory/"+Name+"/provisioner/terraform.tfvars", GetClusterConfig())
	return DistOSMap[awsInstanceOS].User, awsInstanceOS
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
	// check if terraform is available
	DependencyCheck("terraform")

	if _, err := os.Stat("./inventory/" + Name + "/installer"); err == nil {
		fmt.Println("Configuration folder already exists")
	} else {
		os.MkdirAll("./inventory/"+Name+"/provisioner", 0755)
		exec.Command("cp", "-rfp", "./kubespray/contrib/terraform/aws/", "./inventory/"+Name+"/provisioner").Run()
		distSelect()
		terrInit := exec.Command("terraform", "init")
		terrInit.Dir = "./inventory/" + Name + "/provisioner/"
		out, _ := terrInit.StdoutPipe()
		terrInit.Start()
		scanInit := bufio.NewScanner(out)
		for scanInit.Scan() {
			m := scanInit.Text()
			fmt.Println(m)
		}
		terrInit.Wait()
	}

	terrSet := exec.Command("terraform", "apply", "-var-file=credentials.tfvars", "-auto-approve")
	terrSet.Dir = "./inventory/" + Name + "/provisioner/"
	stdout, _ := terrSet.StdoutPipe()
	terrSet.Stderr = terrSet.Stdout
	terrSet.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	terrSet.Wait()
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
		mvShhBastion := exec.Command("mv", "./kubespray/ssh-bastion.conf", "./inventory/"+Name+"/ssh-bastion.conf")
		mvShhBastion.Run()
		mvShhBastion.Wait()
		//os.MkdirAll("./inventory/"+Name+"/installer/group_vars", 0755)
		cpSample := exec.Command("cp", "-rfp", "./kubespray/inventory/sample/", "./inventory/"+Name+"/installer/")
		cpSample.Run()
		cpSample.Wait()

		cpKube := exec.Command("cp", "-rfp", "./kubespray/", "./inventory/"+Name+"/installer/")
		cpKube.Run()
		cpKube.Wait()

		mvInstallerHost := exec.Command("cp", "./inventory/"+Name+"/hosts", "./inventory/"+Name+"/installer/hosts")
		mvInstallerHost.Run()
		mvInstallerHost.Wait()

		// Check if Kubeadm is enabled
		EnableKubeadm()

		//Start Kubernetes Installation
		//check if ansible host file exists
		_, err = os.Stat("./inventory/" + Name + "/hosts")

		ErrorCheck("./inventory/"+Name+"/installer/hosts inventory file not found", err)

		//Enable load balancer api access and copy the kubeconfig file locally
		loadBalancerName, err := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./inventory/"+Name+"/hosts | cut -d'=' -f2").CombinedOutput()

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

			elbNameRaw, _ := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./inventory/"+Name+"/hosts | cut -d'=' -f2 | sed 's/\"//g'").CombinedOutput()

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
	sshUser, osLabel := distSelect()
	fmt.Printf("\nStarting playbook for user %s with os %s\n", sshUser, osLabel)
	kubeSet := exec.Command("ansible-playbook", "-i", "hosts", "cluster.yml", "--timeout=60", "-e ansible_user="+sshUser, "-e bootstrap_os="+osLabel, "-b", "--become-user=root", "--flush-cache")
	kubeSet.Dir = "./inventory/" + Name + "/installer/"
	stdout, _ := kubeSet.StdoutPipe()
	kubeSet.Stderr = kubeSet.Stdout
	kubeSet.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	kubeSet.Wait()

	return
}

// AWSDestroy is used to destroy the infrastructure created.
func AWSDestroy() {
	// check if terraform is installed
	DependencyCheck("terraform")

	// Check if credentials file exist, if it exists skip asking to input the AWS values
	if _, err := os.Stat("./inventory/" + Name + "/provisioner/credentials.tfvars"); err == nil {
		fmt.Println("Credentials file already exists, creation skipped")
	} else {

		ParseTemplate(templates.Credentials, "./inventory/"+Name+"/provisioner/credentials.tfvars", GetCredentials())
	}
	cpHost := exec.Command("cp", "./inventory/"+Name+"/hosts", "./inventory/hosts")

	cpHost.Run()
	cpHost.Wait()
	terrSet := exec.Command("terraform", "destroy", "-var-file=credentials.tfvars", "-force")
	terrSet.Dir = "./inventory/" + Name + "/provisioner/"
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
	}

	terrSet.Wait()

	exec.Command("rm", "./inventory/hosts").Run()
	exec.Command("rm", "-rf", "./inventory/"+Name).Run()

	return
}

// AWSScale is used to scale the AWS infrastructure and Kubernetes cluster.
func AWSScale() {
	var confirmation string
	// Scale the AWS infrastructure
	fmt.Printf("\t\t===============Starting AWS Scaling====================\n\n")
	DependencyCheck("terraform")

	distSelect()

	terrSet := exec.Command("terraform", "apply", "-var-file=credentials.tfvars", "-auto-approve")
	terrSet.Dir = "./kubespray/contrib/terraform/aws/"
	terrSet.Stdout = os.Stdout
	terrSet.Stderr = os.Stderr

	err := terrSet.Start()
	ErrorCheck("Error in destroying the cluster: ", err)
	terrSet.Wait()

	// Scale the Kubernetes cluster
	fmt.Printf("\n\n\t\t===============Starting Kubernetes Scaling====================\n\n")
	DependencyCheck("ansible")
	_, err = os.Stat("./kubespray/inventory/awscluster")
	ErrorCheck("No host file found.", err)

	fmt.Printf("\n\nThis will overwrite the previous host file with a new one. Type \"yes\" to confirm:\n")
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		fmt.Printf("Confirmation denied. Exiting...")
		os.Exit(0)
	}
	// Update the inventory host file
	exec.Command("cp", "-rfp", "./kubespray/inventory/sample/", "./kubespray/inventory/awscluster/").Run()
	exec.Command("cp", "./kubespray/inventory/hosts", "./kubespray/inventory/awscluster/hosts").Run()

	sshUser, osLabel := distSelect()
	fmt.Printf("\n\nStarting Scaling playbook for user %s with os %s\n", sshUser, osLabel)
	scale := exec.Command("ansible-playbook", "-i", "./inventory/awscluster/hosts", "./scale.yml", "-e ansible_user="+sshUser, "-v", "-b", "--timeout=60")
	scale.Dir = "./kubespray/"
	scale.Stdout = os.Stdout
	scale.Stderr = os.Stderr

	scale.Start()
	scale.Wait()

	return
}

// AWSReset is used to reset the AWS infrastructure and removing Kubernetes from it.
func AWSReset() {
	DependencyCheck("ansible")
	sshUser, osLabel := distSelect()
	fmt.Printf("\nStarting playbook for user %s with os %s\n", sshUser, osLabel)
	reset := exec.Command("ansible-playbook", "-i", "./inventory/awscluster/hosts", "./reset.yml", "--timeout=60", "-e ansible_user="+sshUser, "-e delete_nodes_confirmation=yes", "-v", "-b", "--become-user=root", "--flush-cache")
	reset.Dir = "./kubespray/"
	reset.Stdout = os.Stdout
	reset.Stdin = os.Stdin
	reset.Stderr = os.Stderr

	reset.Start()
	reset.Wait()

	return
}

func AWSRemove() {

}
