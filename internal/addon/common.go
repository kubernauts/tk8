package addon

import (
	"fmt"
	"os"
	"os/exec"
)

// KubeConfig provide the path to the local kube config
var KubeConfig string

func cloneExample(addonName string) {
	if _, err := os.Stat("./addons/" + addonName); err == nil {
		fmt.Println("Addon", addonName, "already exist")
		os.Exit(0)
	}
	cEx := exec.Command("git", "clone", "https://github.com/kubernauts/tk-addon-develop", addonName)
	cEx.Dir = "./addons"
	fmt.Println("Create example for your addon", addonName)
	cEx.Run()
	cEx.Wait()
}

func cloneGit(addonPath string) error {
	cEx := exec.Command("git", "clone", addonPath)
	cEx.Dir = "./addons"
	err := cEx.Run()
	if err != nil {
		fmt.Println("Addon could not found. Please check the name or create your own addon. See addon --help")
		os.Exit(0)
	}
	err = cEx.Wait()
	fmt.Println("Addon download finished.")
	return nil
}

func replaceGit(addonName string) {
	cEx := exec.Command("rm", "-rf", ".git")
	cEx.Dir = "./addons/" + addonName
	cEx.Run()
	cEx.Wait()
}
