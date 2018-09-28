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
