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

package addon

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InstallAddon(addonNameOrGitPath string) {
	addonName := GetAddon(addonNameOrGitPath)
	fmt.Println("Install", strings.Replace(addonName, "tk8-addon-", "", 1), addonName)

	executeMainSh(addonName)
	applyMainYml(addonName, "yml")
	applyMainYml(addonName, "yaml")
	fmt.Println(addonName, "installation complete")

}

func applyMainYml(addonName string, fileType string) {

	var cEx *exec.Cmd
	_, err := os.Stat("./addons/" + addonName + "/main." + fileType)
	if err == nil {
		fmt.Println("apply main." + fileType)
		if len(KubeConfig) > 1 {
			cEx = exec.Command("kubectl", "--kubeconfig", KubeConfig, "apply", "-f", "main."+fileType)
		} else {
			cEx = exec.Command("kubectl", "apply", "-f", "main."+fileType)
		}
		cEx.Dir = "./addons/" + addonName
		printTerminalLog(cEx)
		cEx.Wait()
		return
	}
}

func executeMainSh(addonName string) {
	if _, err := os.Stat("./addons/" + addonName + "/main.sh"); err == nil {
		fmt.Println("execute main.sh")
		cEx := exec.Command("/bin/sh", "./main.sh")
		cEx.Dir = "./addons/" + addonName
		printTerminalLog(cEx)
		cEx.Wait()
		return
	}
}

func printTerminalLog(cEx *exec.Cmd) {
	cExOutput, _ := cEx.StdoutPipe()
	cEx.Stderr = cEx.Stdout
	cEx.Start()
	scanner := bufio.NewScanner(cExOutput)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func extractAddonName(addonNameOrGitPath string) string {
	if strings.Contains(addonNameOrGitPath, "/") {
		stringParts := strings.Split(addonNameOrGitPath, "/")
		return stringParts[len(stringParts)-1:][0]
	}
	return addonNameOrGitPath
}

func checkLocalPath(addonName string) bool {
	if _, err := os.Stat("./addons/" + addonName); err == nil {
		fmt.Println("Addon", addonName, "already exist")
		return true
	}
	return false
}
