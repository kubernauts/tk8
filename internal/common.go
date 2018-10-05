package common

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func CloneGit(executeDir string, gitUrl string, targetFolder string) error {
	os.Mkdir(executeDir, 0755)
	cEx := exec.Command("git", "clone", gitUrl, targetFolder)
	cEx.Dir = executeDir
	stdout, _ := cEx.StdoutPipe()
	cEx.Stderr = cEx.Stdout
	error := cEx.Start()
	if error != nil {
		fmt.Println(error)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	cEx.Wait()
	return nil
}

func ReplaceGit(executeDir string) {
	cEx := exec.Command("rm", "-rf", ".git")
	cEx.Dir = executeDir
	cEx.Run()
	cEx.Wait()
}
