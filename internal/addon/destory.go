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
