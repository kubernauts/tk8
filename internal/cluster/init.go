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
	"os"
	"os/exec"
)

// KubesprayInit is responsible for cloning the kubespray diretory in CWD.
func KubesprayInit() {
	if _, err := os.Stat("./kubespray"); err == nil {
		fmt.Println("Kubespray clone on this system already exists")
		return
	}
	fmt.Println("Initialising kubespray git repo")
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal("either 'git' is not installed or not found in $PATH, kindly check and fix")
		os.Exit(1)
	} else {
		// issue-23 kubespray upstream
		gitClone := exec.Command("git", "clone", "-b", kubesprayVersion, "--single-branch", "https://github.com/kubernauts/kubespray")
		stdout, err := gitClone.StdoutPipe()
		gitClone.Stderr = gitClone.Stdout
		gitClone.Start()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
		gitClone.Wait()
		if err != nil {
			log.Fatalf("Seems there is a problem cloning the kubespray repo, %v", err)
			os.Exit(1)
		}
	}
	if _, err := exec.LookPath("pip"); err != nil {
		log.Fatal("either 'pip' is not installed or not found in $PATH, kindly check and fix")
		os.Exit(1)
	} else {
		// Ensure to have all the dependencies of Kubespray
		err := exec.Command("pip", "install", "-r", "kubespray/requirements.txt").Run()
		if err != nil {
			log.Fatalf("Seems there is a problem installing the kubespray dependencies. " +
				"Please run following command:" +
				"\n\n\t" +
				"pip install -r kubespray/requirements.txt" +
				"\nNOTE: Elevated permission may require.")
			os.Exit(1)
		}
	}
}
