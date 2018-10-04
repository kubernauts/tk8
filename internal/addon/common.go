package addon

import (
	"fmt"
	"os"

	"github.com/kubernauts/tk8/internal"
)

// KubeConfig provide the path to the local kube config
var KubeConfig string

func cloneExample(addonName string) {
	if _, err := os.Stat("./addons/" + addonName); err == nil {
		fmt.Println("Addon", addonName, "already exist")
		os.Exit(0)
	}
	common.CloneGit("./addons", "https://github.com/kubernauts/tk-addon-develop", addonName)
}
