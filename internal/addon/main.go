package addon

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kubernauts/tk8/pkg/common"
)

type Addon struct {
}

func (a *Addon) Create(addonName string) (error, string) {
	cloneExample(addonName)
	common.ReplaceGit("./addons/" + addonName)
	return nil, addonName
}

func (a *Addon) Destroy(addonNameOrGitPath string) (error, string) {
	_, addonName := a.Get(addonNameOrGitPath)
	fmt.Println("Destroying", strings.Replace(addonName, "tk8-addon-", "", 1))
	executeMainSh(addonName)
	deleteMainYml(addonName, "yml")
	deleteMainYml(addonName, "yaml")
	fmt.Println(strings.Replace(addonName, "tk8-addon-", "", 1), "destroy complete")

	return nil, addonName
}

func (a *Addon) Get(addonNameOrGitPath string) (error, string) {
	addonName := extractAddonName(addonNameOrGitPath)
	fmt.Println("Search local for", addonName)

	if checkLocalPath(addonName) {
		fmt.Println("Found", addonName, "local.")
		return nil, addonName
	}
	if checkLocalPath("tk8-addon-" + addonName) {
		fmt.Println("Found", addonName, "local.")
		return nil, "tk8-addon-" + addonName
	}
	if !checkLocalPath(addonName) {
		fmt.Println("check if provided a url")
		if strings.Contains(addonNameOrGitPath, "http://") || strings.Contains(addonNameOrGitPath, "https://") {
			fmt.Println("Load Addon from external path", addonNameOrGitPath)
			common.CloneGit("./addons", addonNameOrGitPath, addonName)
			return nil, addonName
		}

		fmt.Println("Search addon on kubernauts space.")
		common.CloneGit("./addons", "https://github.com/kubernauts/tk8-addon-"+addonName, addonName)
		return nil, "tk8-addon-" + addonName

	}
	return nil, "tk8-addon-" + addonName

}

func (a *Addon) Install(addonNameOrGitPath string) {
	_, addonName := a.Get(addonNameOrGitPath)
	fmt.Println("Install", strings.Replace(addonName, "tk8-addon-", "", 1), addonName)

	executeMainSh(addonName)
	applyMainYml(addonName, "yml")
	applyMainYml(addonName, "yaml")
	fmt.Println(addonName, "installation complete")

}

// KubeConfig provide the path to the local kube config
var KubeConfig string

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

func cloneExample(addonName string) {
	if _, err := os.Stat("./addons/" + addonName); err == nil {
		fmt.Println("Addon", addonName, "already exist")
		os.Exit(0)
	}
	common.CloneGit("./addons", "https://github.com/kubernauts/tk8-addon-develop", addonName)
}

func deleteMainYml(addonName string, fileType string) {

	var cEx *exec.Cmd
	if _, err := os.Stat("./addons/" + addonName + "/main." + fileType); err == nil {
		fmt.Println("delete", strings.Replace(addonName, "tk8-addon-", "", 1), "from cluster")
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
