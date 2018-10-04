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
