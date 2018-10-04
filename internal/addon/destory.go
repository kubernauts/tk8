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
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func DestroyAddon(addonNameOrGitPath string) {
	addonName := GetAddon(addonNameOrGitPath)
	fmt.Println("Destroying", strings.Replace(addonName, "tk-addon-", "", 1))
	executeMainSh(addonName)
	deleteMainYml(addonName, "yml")
	deleteMainYml(addonName, "yaml")
	fmt.Println(strings.Replace(addonName, "tk-addon-", "", 1), "destroy complete")

}

func deleteMainYml(addonName string, fileType string) {

	var cEx *exec.Cmd
	if _, err := os.Stat("./addons/" + addonName + "/main." + fileType); err == nil {
		fmt.Println("delete", strings.Replace(addonName, "tk-addon-", "", 1), "from cluster")
		if len(KubeConfig) > 1 {
			cEx = exec.Command("kubectl", "--kubeconfig="+KubeConfig, "delete", "-f", "main."+fileType)
		} else {
			cEx = exec.Command("kubectl", "delete", "-f", "main."+fileType)
		}
		cEx.Dir = "./addons/" + addonName
		printTerminalLog(cEx)
		cEx.Wait()
		return
	}
}
