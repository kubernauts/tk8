package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alecthomas/template"
)

var Name string

type AzureConfig struct {
	ClusterName         string
	AWSRegion           string
	NodeInstanceType    string
	DesiredCapacity     int
	AutoScallingMinSize int
	AutoScallingMaxSize int
	KeyPath             string
}

// GetFilePath fetches and returns the current working directory.
func GetFilePath(fileName string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, fileName)
}

func ParseTemplate(templateString string, outputFileName string, data interface{}) {
	// open template
	template := template.New("template")
	template, _ = template.Parse(templateString)
	// open output file
	outputFile, err := os.Create(GetFilePath(outputFileName))
	defer outputFile.Close()
	if err != nil {
		ExitErrorf("Error creating file %s: %v", outputFile, err)
	}
	err = template.Execute(outputFile, data)
	ErrorCheck("Error executing template: %v", err)

}

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

func SetClusterName() {

}
