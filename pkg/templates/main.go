package templates

import (
	"os"

	"github.com/alecthomas/template"
	"github.com/kubernauts/tk8/pkg/common"
)

func ParseTemplate(templateString string, outputFileName string, data interface{}) {
	// open template
	template := template.New("template")
	template, _ = template.Parse(templateString)
	// open output file
	outputFile, err := os.Create(common.GetFilePath(outputFileName))
	defer outputFile.Close()
	if err != nil {
		common.ExitErrorf("Error creating file %s: %v", outputFile, err)
	}
	err = template.Execute(outputFile, data)
	common.ErrorCheck("Error executing template: %v", err)

}
