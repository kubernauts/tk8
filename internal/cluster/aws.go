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
	"text/template"

	"github.com/kubernauts/tk8/internal/templates"
)

var ec2IP string

func parseTemplate(templateString string, outputFileName string, data interface{}) {
	// open template
	template := template.New("template")
	template, _ = template.Parse(templateString)
	// open output file
	outputFile, err := os.Create(GetFilePath(outputFileName))
	defer outputFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(outputFile, data)
}

func distSelect() (string, string) {
	var sshUser string

	//Read Configuration File
	ReadViperConfigFile("config")
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
		go parseTemplate(templates.CustomInfrastructure, "./kubespray/contrib/terraform/aws/create-infrastructure.tf", DistOSMap[awsInstanceOS])
	} else {
		go parseTemplate(templates.Infrastructure, "./kubespray/contrib/terraform/aws/create-infrastructure.tf", DistOSMap[awsInstanceOS])
	}

	go parseTemplate(templates.Credentials, "./kubespray/contrib/terraform/aws/credentials.tfvars", GetCredentials())
	go parseTemplate(templates.Variables, "./kubespray/contrib/terraform/aws/variables.tf", DistOSMap[awsInstanceOS])
	go parseTemplate(templates.Terraform, "./kubespray/contrib/terraform/aws/terraform.tfvars", GetClusterConfig())

	return DistOSMap[awsInstanceOS].User, DistOSMap[awsInstanceOS].OS
}

func AWSCreate() {
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

	distSelect()

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

func AWSInstall() {
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
	sshUser, osLabel := distSelect()
	kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/awscluster/hosts", "./cluster.yml", "--timeout=60", "-e ansible_user="+sshUser, "-e bootstrap_os="+osLabel, "-b", "--become-user=root", "--flush-cache")
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

func AWSDestroy() {
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
		parseTemplate(templates.Credentials, "./kubespray/contrib/terraform/aws/credentials.tfvars", GetCredentials())
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
