package common

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	Name string
	// GITCOMMIT will hold the commit SHA to be used in the version command.
	GITCOMMIT = "0"
	// VERSION will hold the version number to be used in the version command.
	VERSION = "dev"
)

// ErrorCheck is responsbile to check if there is any error returned by a command.
func ErrorCheck(msg string, err error) {
	if err != nil {
		ExitErrorf(msg, err)
	}
}

// DependencyCheck check if the binary is installed
func DependencyCheck(bin string) {
	_, err := exec.LookPath(bin)
	ErrorCheck(bin+" not found.", err)

	_, err = exec.Command(bin, "--version").Output()
	ErrorCheck("Error executing "+bin, err)
}

// ExitErrorf exits the program with an error code of '1' and an error message.
func ExitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

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

// GetFilePath fetches and returns the current working directory.
func GetFilePath(fileName string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, fileName)
}
