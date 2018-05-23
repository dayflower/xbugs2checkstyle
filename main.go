package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/phayes/checkstyle"
)

func addCheckstyleError(files map[string]*checkstyle.File, bug BugInstance, source SourceLine) {
	sourcePath := source.SourcePath

	file, exist := files[sourcePath]
	if !exist {
		file = checkstyle.NewFile(sourcePath)
		files[sourcePath] = file
	}

	file.AddError(checkstyle.NewError(source.Start, 1, checkstyle.SeverityWarning, bug.Type, bug.Type))
}

func main() {
	body, _ := ioutil.ReadAll(os.Stdin)

	document := BugCollection{}
	xml.Unmarshal(body, &document)

	checkstyleFiles := make(map[string]*checkstyle.File)

	for _, bug := range document.BugInstances {
		for _, source := range bug.ClassSourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
		for _, source := range bug.MethodSourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
		for _, source := range bug.FieldSourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
		for _, source := range bug.SourceLines {
			addCheckstyleError(checkstyleFiles, bug, source)
		}
	}

	checkstyleResult := checkstyle.New()
	for _, file := range checkstyleFiles {
		checkstyleResult.AddFile(file)
	}

	buf, _ := xml.MarshalIndent(checkstyleResult, "", "\t")
	fmt.Println(string(buf))
}
