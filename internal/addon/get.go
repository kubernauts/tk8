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
	"strings"

	"github.com/kubernauts/tk8/internal"
)

func GetAddon(addonNameOrGitPath string) string {
	addonName := extractAddonName(addonNameOrGitPath)
	fmt.Println("Search local for", addonName)

	if checkLocalPath(addonName) {
		fmt.Println("Found", addonName, "local.")
		return addonName
	}
	if checkLocalPath("tk8-addon-" + addonName) {
		fmt.Println("Found", addonName, "local.")
		return "tk8-addon-" + addonName
	}
	if !checkLocalPath(addonName) {
		fmt.Println("check if provided a url")
		if strings.Contains(addonNameOrGitPath, "http://") || strings.Contains(addonNameOrGitPath, "https://") {
			fmt.Println("Load Addon from external path", addonNameOrGitPath)
			common.CloneGit("./addons", addonNameOrGitPath, addonName)
			return addonName
		}

		fmt.Println("Search addon on kubernauts space.")
		common.CloneGit("./addons", "https://github.com/kubernauts/tk8-addon-"+addonName, addonName)
		return "tk8-addon-" + addonName

	}
	return "tk8-addon-" + addonName

}
