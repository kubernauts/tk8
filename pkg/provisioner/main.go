package provisioner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kubernauts/tk8/pkg/common"
)

var IOnly bool = false
var Interactive bool

type Provisioner interface {
	Init(args []string)
	Setup(args []string)
	Scale(args []string)
	Remove(args []string)
	Reset(args []string)
	Upgrade(args []string)
	Destroy(args []string)
}

func NotImplemented() {
	fmt.Println("Not implemented yet. Coming soon...")
	os.Exit(0)
}

func ExecuteTerraform(command string, path string) {

	common.DependencyCheck("terraform")
	var terrSet *exec.Cmd

	if strings.Compare(strings.TrimRight(command, "\n"), "init") == 0 {
		terrSet = exec.Command("terraform", command, "-var-file=credentials.tfvars")
	} else if strings.Compare(command, "apply") == 0 {
		terrSet = exec.Command("terraform", command, "-var-file=credentials.tfvars", "-auto-approve")
	} else {
		terrSet = exec.Command("terraform", command, "-var-file=credentials.tfvars", "-force")
	}

	terrSet.Dir = path
	stdout, _ := terrSet.StdoutPipe()
	terrSet.Stderr = terrSet.Stdout
	error := terrSet.Start()
	if error != nil {
		fmt.Println(error)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		if strings.Contains(m, "Error: Error applying plan") {
			fmt.Println("Terraform could not setup the infrastructure")
			os.Exit(1)
		}
	}

	terrSet.Wait()
}
