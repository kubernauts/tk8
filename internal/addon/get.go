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
)

func GetAddon(addonNameOrGitPath string) string {
	addonName := extractAddonName(addonNameOrGitPath)
	fmt.Println("Search local for", addonName)

	if checkLocalPath(addonName) {
		fmt.Println("Found", addonName, "local.")
		return addonName
	}
	if checkLocalPath("tk-addon-" + addonName) {
		fmt.Println("Found", addonName, "local.")
		return "tk-addon-" + addonName
	}
	if !checkLocalPath(addonName) {
		fmt.Println("check if provided a url")
		if strings.Contains(addonNameOrGitPath, "http://") || strings.Contains(addonNameOrGitPath, "https://") {
			fmt.Println("Load Addon from external path", addonNameOrGitPath)
			cloneGit(addonNameOrGitPath)
			return addonName
		}

		fmt.Println("Search addon on kubernauts space.")
		cloneGit("https://github.com/kubernauts/tk-addon-" + addonName)
		return "tk-addon-" + addonName

	}
	return "tk-addon-" + addonName

}
