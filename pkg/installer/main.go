package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kubernauts/tk8/pkg/common"
)

func RunPlaybook(path string, file string, sshUser string, osLabel string) {
	common.DependencyCheck("ansible")
	fmt.Printf("\nStarting playbook for user %s with os %s\n", sshUser, osLabel)
	ansiblePlaybook := exec.Command("ansible-playbook", "-i", "hosts", file, "--timeout=60", "-e ansible_user="+sshUser, "-e ansible_user="+sshUser, "-e bootstrap_os="+osLabel, "-b", "--become-user=root", "--flush-cache")
	ansiblePlaybook.Dir = path
	ansiblePlaybook.Stdout = os.Stdout
	ansiblePlaybook.Stdin = os.Stdin
	ansiblePlaybook.Stderr = os.Stderr

	ansiblePlaybook.Start()
	ansiblePlaybook.Wait()
}
